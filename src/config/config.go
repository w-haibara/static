package config

import (
	"encoding/json"
	"osoba/auth"
	"osoba/resource"
)

type Config struct {
	Auth *auth.Config     `json: "auth"`
	DB   *resource.Config `json: "db"`
}

func Configure(data []byte) (Config, error) {
	c := Config{}

	err := json.Unmarshal([]byte(data), &c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}
