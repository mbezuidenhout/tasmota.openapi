package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	sw "github.com/mbezuidenhout/tasmota.openapi/go"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v2"
)

const VERSION = "0.0.1"

type Config struct {
	Port     string `yaml:"port"`
	Keyfile  string `yaml:"keyfile"`
	Certfile string `yaml:"certfile"`
	Webpath  string `yaml:"webpath,omitempty"`
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

func printHelp() {
	fmt.Printf("Usage: %s [-f path-to-config.yml]\n", os.Args[0])
}

func main() {
	if slices.Contains(os.Args[1:], "--help") {
		printHelp()
		return
	}

	configFile := "config.yml"
	if slices.Contains(os.Args, "-f") {
		pos := slices.Index(os.Args, "-f")
		configFile = os.Args[pos+1]
	}

	config, err := NewConfig(configFile)
	if err != nil {
		return
	}

	cert, _ := tls.LoadX509KeyPair(config.Certfile, config.Keyfile)
	webPath := "./dist/swaggerui"

	if len(config.Webpath) > 0 {
		webPath = config.Webpath
	}

	fmt.Printf("Starting %s version %s\n", os.Args[0], VERSION)

	s := &http.Server{
		Addr:    ":" + config.Port,
		Handler: sw.NewRouter(webPath),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}
	fmt.Printf("Listening on %s\n", config.Port)

	ticker := time.NewTicker(5 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				i := sw.CleanupConnections()
				if i > 0 {
					fmt.Printf("Closing %d connection(s)", i)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	// Automatically close connections that are open for more than 15 minutes

	log.Fatal(s.ListenAndServeTLS(config.Certfile, config.Keyfile))
}
