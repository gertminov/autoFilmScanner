package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Arduino     `yaml:"arduino"`
	Images      `yaml:"images"`
	WithUI      bool   `yaml:"withUI"`
	Calibrate   bool   `yaml:"calibrate"`
	DefaultPath string `yaml:"defaultPath"`
}

type Arduino struct {
	BoutRate      int    `yaml:"bout_rate"`
	StepsPerPhoto uint32 `yaml:"steps_per_photo"`
	StepsPerDia   uint32 `yaml:"steps_per_dia"`
	ZeroPoint     uint32 `yaml:"zero_point"`
	GoBackSteps   uint32 `yaml:"go_back_steps"`
	GoBack        bool   `yaml:"go_back"`
	TimeOut       int    `yaml:"time_out"`
	StartPosition uint32 `yaml:"start_position"`
	FirstDiaPos   uint32 `yaml:"first_dia_position"`
	FirstImagePos uint32 `yaml:"first_image_position"`
}

type Images struct {
	ImagesPerStrip int  `yaml:"images_per_strip"`
	Dia            bool `yaml:"dia"`
}

var defaultConfig = Config{
	Arduino:   Arduino{9600, 945, 1530, 4, 5, false, 10, 1500, 3500, 3450},
	Images:    Images{6, false},
	WithUI:    false,
	Calibrate: true,
}

func InitConfig() Config {
	config := defaultConfig
	exists, err := checkForConfigFile()
	if err != nil {
		log.Fatal(err)
	}

	if exists {
		config = loadConfig()
	} else {
		WriteConfig(&config)
		fmt.Println("Konfigurationsdatei wure erstellt")
	}
	return config
}

func loadConfig() Config {
	file, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	config := Config{}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Println("config.yaml file seems to be broken")
		log.Fatal(err)
	}

	return config
}

func checkForConfigFile() (bool, error) {
	_, err := os.Stat("./config.yaml")
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func WriteConfig(config *Config) {
	marshal, err := yaml.Marshal(&config)
	fmt.Println(config.ImagesPerStrip)
	if err != nil {
		log.Fatal(err)
	}

	fileName := "config.yaml"
	err = os.WriteFile(fileName, marshal, 0777)
	if err != nil {
		log.Fatal(err)
	}
}
