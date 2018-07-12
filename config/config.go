// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period      time.Duration `config:"period"`
	APIClientID string        `config:"api_client_ID"`
	APIKey      string        `config:"api_Key"`
}

var DefaultConfig = Config{
	Period:      1 * time.Second,
	APIClientID: "f6cab9156394e0bc768b",
	APIKey:      "1fcd292c-2059-4235-8f76-62d1a6cf9db3",
}
