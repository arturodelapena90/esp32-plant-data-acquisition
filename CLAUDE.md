# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

Go firmware for an ESP32-S3 (generic S3-N16R8 board), compiled with TinyGo, that reads sensors on a habanero plant (BH1750 light over I2C, DHT22 temperature/humidity over GPIO, 2x capacitive soil moisture over ADC) and publishes JSON readings to an MQTT broker over WiFi. The broker and WiFi AP are both hosted on a Raspberry Pi.

## Build / flash / monitor

This is TinyGo firmware, not a regular Go binary — `go build` will fail because it imports the `machine` package, which only resolves under the TinyGo toolchain. Requires TinyGo 0.41+ (both `go` 1.24.4 and `tinygo` 0.41.1 are installed) — that's the release that added native ESP32/ESP32-S3 WiFi support via `tinygo.org/x/espradio`; there is no `machine.WIFI` API.

```bash
# Build
tinygo build -target=esp32s3-generic -o firmware.bin main.go

# Flash
tinygo flash -target=esp32s3-generic main.go

# Monitor serial output
tinygo monitor -port=/dev/ttyUSB0 -baud=921600
```

Note the target is `esp32s3-generic`, not `esp32-s3` (that name doesn't exist) or bare `esp32s3` (inheritable-only, TinyGo refuses to build with it directly).

Use `go vet ./...` / `gofmt` for quick static checks on non-hardware code, but treat `tinygo build` as the real compile check since hardware-specific files are gated behind `//go:build tinygo` (see below) — plain `go build`/`go vet` cannot type-check them at all.

## Configuration

Runtime config comes from environment variables (`.env`, loaded via `github.com/caarlos0/env/v10` in `internal/config/config.go`), not from flashed-in constants:

- `RASPBERRY_PI_IP` — Pi's IP; MQTT broker address is derived as `<ip>:1883` (no `mqtt://` scheme — it's dialed directly with `net.Dial("tcp", ...)`)
- `WIFI_SSID`, `WIFI_PASSWORD` — Pi's AP credentials
- `MQTT_TOPIC` — publish topic
- `MQTT_CLIENT_ID` — optional, defaults to `esp32-habanero-01`

Hardware pin assignments (`DHT22Pin`, `SoilPin1`, `SoilPin2`, `I2CSDAPin`, `I2CSCLPin`) and `ReadInterval` are hardcoded defaults in `config.LoadConfig()`, not env-driven. Pins are chosen to avoid collisions — the I2C bus uses GPIO8/9, so the soil ADC pins deliberately avoid those. `.env` is gitignored — don't commit real credentials into it or into source.

## Architecture

Each sensor type lives in its own package under `internal/sensor/` and follows the same shape:

- A platform-independent file (e.g. `light.go`, `climate.go`, `soil.go`) defining `Sensor`, `Reading`, `New(...)`, `Start(interval, chan<- Reading)`, and `Read()`.
- A `//go:build tinygo`-gated file (e.g. `bh1750.go`, `dht22.go`, `sto160.go`) with the actual register/protocol-level driver code. This split exists so the non-hardware logic can at least be parsed/vetted with plain `go`, while the real hardware access only builds under TinyGo.

`internal/sensor/climate/dhtdriver/` is a vendored, patched copy of `tinygo.org/x/drivers/dht` — that upstream package (as of v0.35.0 and the current `dev` branch) doesn't compile for esp32s3 (`machine.CPUFrequency()` doesn't exist there) and its bit-timing `counter` type would overflow at esp32s3's ~240MHz clock even if it did. See `dhtdriver/doc.go` for the specifics. Don't route DHT22 changes through the upstream `tinygo.org/x/drivers/dht` import — this package is a deliberate fork, not a mistake.

Data flow, wired up in `main.go`:

1. `config.LoadConfig()` reads env vars.
2. WiFi: `link.Esplink{}` (from `tinygo.org/x/espradio/netlink`) is registered via `netdev.UseNetdev()`, then `NetConnect()` joins the Pi's AP. This is the standard TinyGo pattern (`netdev`/`netlink` abstraction) — after this, ordinary `net.Dial` works.
3. A TCP connection is dialed to `cfg.MQTTBroker` (`ip:1883`), and `mqtt.SetupMQTT()` (wrapping `github.com/soypat/natiu-mqtt`) connects over it using `cfg.MQTTClientID`.
4. I2C bus (`machine.I2C1`, a `*machine.I2C`, on `cfg.I2CSDAPin`/`cfg.I2CSCLPin`) is configured once and shared; each sensor's `New()` is called to initialize it. `light.New` takes the bus by pointer (`*machine.I2C`) since all of `machine.I2C`'s methods have pointer receivers.
5. Each sensor runs `Start()` in its own goroutine, ticking on `cfg.ReadInterval` and pushing a `Reading` onto its own channel (light immediately reads once before the first tick; climate/soil wait for the first tick).
6. `aggregator.Start()` runs in its own goroutine: it blocks on receiving one reading from each of the four channels (light, climate, soil1, soil2), merges them into an `mqtt.Data` payload with a fresh Unix timestamp, and pushes it to `mqttChan`. Because it waits on all four channels every loop, publish rate is effectively gated by the slowest sensor — and a sensor that skips a send on error (light does; climate/soil always send, error or not) stalls that cycle until it succeeds again.
7. `mqttClient.Publish()` drains `mqttChan`, JSON-marshals each `mqtt.Data`, and publishes it (QoS0) to `cfg.MQTTTopic`.

All sensor `Reading` fields are `*float32` pointers (nil on read error), and `mqtt.Data` mirrors that — a nil field means that sensor failed on that cycle but the payload is still published with the other readings.

Logging is plain `fmt.Printf`/`fmt.Println` called directly at each site — there is no logger type or injected dependency. This was `go.uber.org/zap` originally, but measurement (`tinygo build -size short`) showed it cost ~150KB flash / ~43KB RAM on esp32s3 for structured/leveled logging with no consumer beyond a serial monitor, so it was dropped in favor of stdlib `fmt`. Fatal startup errors in `main.go` use `panic(fmt.Errorf("...: %w", err))` rather than a `Fatalf`-style helper — TinyGo has no meaningful `os.Exit`, and panic's own runtime output serves as the log line.
