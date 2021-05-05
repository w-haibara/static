package osoba

import (
	"encoding/json"
	"os"
)

type Config struct {
	DocumentRoot         string `documentRoot`
	TmpDirContentsPrefix string `tmpDirContentsPrefix`
	Contents             []struct {
		Path `string`
		URL  `url`
	} `contents`
}

func (a *App) LoadConfig() {
	b, err := os.ReadFile("config.json")
	if err != nil {
		panic(err.Error())
	}

	c := Config{}
	if err := json.Unmarshal(b, &c); err != nil {
		panic(err.Error())
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
		V: map[Path]URL{},
	}
	for _, v := range c.Contents {
		a.Contents.Create(string(v.Path), string(v.URL))
	}
}
