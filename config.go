package osoba

import (
	"encoding/json"
	"os"
)

type Config struct {
	DocumentRoot         string `documentRoot`
	TmpDirContentsPrefix string `tmpDirContentsPrefix`
	Contents             []struct {
		Path   string `string`
		URL    string `url`
		Secret string `secret`
	} `contents`
}

func (a *App) LoadConfig() error {
	b, err := os.ReadFile("config.json")
	if err != nil {
		return err
	}

	c := Config{}
	if err := json.Unmarshal(b, &c); err != nil {
		return err
	}

	if c.DocumentRoot == "" {
		a.DocumentRoot = "www"
	} else {
		a.DocumentRoot = c.DocumentRoot
	}

	if c.TmpDirContentsPrefix == "" {
		a.TmpDirContentsPrefix = "contents"
	} else {
		a.TmpDirContentsPrefix = c.TmpDirContentsPrefix
	}

	a.Contents = Contents{
		V: map[string]Content{},
	}
	for _, v := range c.Contents {
		if err := a.Contents.Create(v.Path, v.URL, v.Secret); err != nil {
			return err
		}
	}

	return nil
}
