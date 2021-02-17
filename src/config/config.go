package config

import (
	"encoding/json"
	"osoba/auth"
)

type Config struct {
	Auth *auth.Config `json: "auth"`
}

func Configure(data []byte) (Config, error) {
	c := Config{}

	err := json.Unmarshal([]byte(data), &c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}
