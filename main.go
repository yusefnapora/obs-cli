package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/andreykaipov/goobs"
	"github.com/muesli/coral"
)

var (
	host     string
	password string
	port     uint32
	version  string

	rootCmd = &coral.Command{
		Use:   "obs-cli",
		Short: "obs-cli is a command-line remote control for OBS",
	}

	client *goobs.Client
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if client != nil {
		_ = client.Disconnect()
	}
}

func init() {
	coral.OnInitialize(connectOBS)
	rootCmd.PersistentFlags().StringVar(&host, "host", "localhost", "host to connect to")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "password for connection")
	rootCmd.PersistentFlags().Uint32VarP(&port, "port", "p", 4455, "port to connect to")

	if host == "localhost" && password == "" && port == 4455 {
		type (
			connection struct {
				Host     string
				Port     int
				Password string
			}

			config struct {
				Connection map[string]connection
			}
		)

		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		f := filepath.Join(homeDir, ".obs_cli", "config.toml")

		if _, err := os.Stat(f); err == nil {
			var c config
			_, err := toml.DecodeFile(f, &c.Connection)
			if err != nil {
				log.Fatal(err)
			}
			conn := c.Connection["connection"]
			host = conn.Host
			port = uint32(conn.Port)
			password = conn.Password
		}
	}
}

func getUserAgent() string {
	userAgent := "obs-cli"
	if version != "" {
		userAgent += "/" + version
	}
	return userAgent
}

func connectOBS() {
	var err error
	client, err = goobs.New(
		host+fmt.Sprintf(":%d", port),
		goobs.WithPassword(password),
		goobs.WithRequestHeader(http.Header{"User-Agent": []string{getUserAgent()}}),
	)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
