package config

import (
	"machine"
	"time"

	"github.com/caarlos0/env/v10"
)

const mqttPort = "1883"

type Config struct {

	// Raspberry Pi
	RaspberryPiIP string `env:"RASPBERRY_PI_IP,required"`

	// WiFi
	WifiSSID     string `env:"WIFI_SSID,required"`
	WifiPassword string `env:"WIFI_PASSWORD,required"`

	// MQTT
	MQTTTopic    string `env:"MQTT_TOPIC,required"`
	MQTTClientID string `env:"MQTT_CLIENT_ID,required"`
	MQTTBroker   string

	// Hardware Settings
	DHT22Pin     machine.Pin
	SoilPin1     machine.Pin
	SoilPin2     machine.Pin
	I2CSDAPin    machine.Pin
	I2CSCLPin    machine.Pin
	ReadInterval time.Duration
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		DHT22Pin:     4,
		SoilPin1:     7,
		SoilPin2:     6,
		I2CSDAPin:    8,
		I2CSCLPin:    9,
		ReadInterval: 30 * time.Second,
	}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	cfg.MQTTBroker = cfg.RaspberryPiIP + ":" + mqttPort
	return cfg, nil
}
