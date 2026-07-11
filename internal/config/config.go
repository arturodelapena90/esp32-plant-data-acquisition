package config

import (
	"machine"
	"time"

	"github.com/caarlos0/env/v10"
)

type Config struct {

	// Raspberry Pi
	RaspberryPiIP string `env:"RASPBERRY_PI_IP,required"`

	// WiFi
	WifiSSID     string `env:"WIFI_SSID,required"`
	WifiPassword string `env:"WIFI_PASSWORD,required"`

	// MQTT
	MQTTTopic  string `env:"MQTT_TOPIC,required"`
	MQTTBroker string

	// Hardware Settings
	DHT22Pin     machine.Pin
	SoilPin1     machine.Pin
	SoilPin2     machine.Pin
	ReadInterval time.Duration
}

func (c *Config) GetMQTTBroker() string {
	return "mqtt://" + c.RaspberryPiIP + ":1883"
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		DHT22Pin:     4,
		SoilPin1:     7,
		SoilPin2:     8,
		ReadInterval: 30 * time.Second,
	}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	// add mqtt
	cfg.MQTTBroker = "mqtt://" + cfg.RaspberryPiIP + ":1883"
	return cfg, nil
}
