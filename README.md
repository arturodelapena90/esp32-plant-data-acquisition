# ESP32 S3 Habanero Plant Data Acquisition

Simple, clean Go system to acquire sensor data from an Habanero plant and publish via MQTT.

## Sensors

- **BH1750** - Light intensity (I2C)
- **DHT22** - Temperature & Humidity (GPIO)
- **2x Capacitive Soil Moisture** - Redundant moisture sensors (ADC)

## Architecture

```
main.go:
├── connect WiFi (tinygo.org/x/espradio)
├── dial + connect MQTT (github.com/soypat/natiu-mqtt)
├── configure shared I2C bus
├── light.Start()    → sends light.Reading to a channel
├── climate.Start()  → sends climate.Reading to a channel
├── soil.Start() x2  → sends soil.Reading to a channel each
├── aggregator.Start() → merges all readings into mqtt.Data
└── mqttClient.Publish() → publishes mqtt.Data to the broker
```

Each sensor runs in its own goroutine (`internal/sensor/{light,climate,soil}`). The aggregator waits for a reading from each sensor, combines them into `mqtt.Data`, and sends it to the MQTT publisher. See `CLAUDE.md` for the full data-flow writeup.

## Configuration

Runtime config (WiFi credentials, Pi IP, MQTT topic/client ID) comes from environment variables — see `.env`. Hardware pin assignments and read interval are defaults in `internal/config/config.go`:

```go
DHT22Pin:     4, // GPIO4
SoilPin1:     7, // GPIO7 / ADC1_CH6
SoilPin2:     6, // GPIO6 / ADC1_CH5
I2CSDAPin:    8, // GPIO8
I2CSCLPin:    9, // GPIO9
ReadInterval: 30 * time.Second,
```

## Building for TinyGo

Requires TinyGo 0.41+ (for native ESP32 WiFi support via `tinygo.org/x/espradio`).

```bash
# Build
tinygo build -target=esp32s3-generic -o firmware.bin main.go

# Flash
tinygo flash -target=esp32s3-generic main.go

# Monitor
tinygo monitor -port=/dev/ttyUSB0 -baud=921600
```

## MQTT Message

Published to `MQTT_TOPIC` (see `.env`):

```json
{
  "timestamp": 1688123456,
  "light_lux": 45000.50,
  "temperature_c": 28.5,
  "humidity_percent": 65.2,
  "moisture1_percent": 72.8,
  "moisture2_percent": 71.9
}
```
