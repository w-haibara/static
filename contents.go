package osoba

import (
	"fmt"
)

func (c Contents) Create(path, url, secret string) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if _, ok := c.V[path]; ok {
		return fmt.Errorf("content exist: " + path)
	}

	c.V[path] = Content{
		URL:    url,
		Secret: secret,
	}

	return nil
}

func (c Contents) Update(path, url, secret string) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if _, ok := c.V[path]; !ok {
		return fmt.Errorf("content not exist: " + path)
	}

	c.V[path] = Content{
		URL:    url,
		Secret: secret,
	}

	return nil
}

func (c Contents) Delete(path string) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if _, ok := c.V[path]; !ok {
		return fmt.Errorf("content not exist: " + path)
	}

	c.V[path] = Content{}

	return nil
}

func (c Contents) DeleteDir(path string) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if _, ok := c.V[path]; !ok {
		return fmt.Errorf("content not exist: " + path)
	}

	delete(c.V, path)

	return nil
}
