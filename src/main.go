package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	sw "github.com/mbezuidenhout/tasmota.openapi/go"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Port     string `yaml:"port"`
	Keyfile  string `yaml:"keyfile"`
	Certfile string `yaml:"certfile"`
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func main() {
	config, err := NewConfig("config.yml")
	if err != nil {
		return
	}

	cert, _ := tls.LoadX509KeyPair(config.Certfile, config.Keyfile)

	s := &http.Server{
		Addr:    ":" + config.Port,
		Handler: sw.NewRouter(),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	log.Fatal(s.ListenAndServeTLS(config.Certfile, config.Keyfile))
}
