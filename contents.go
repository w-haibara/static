package osoba

func (c Contents) Create(path, url string) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if _, ok := c.V[Path(path)]; ok {
		panic("content exist: " + path)
	}

	c.V[Path(path)] = URL(url)
}

func (c Contents) Update(path, url string) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if _, ok := c.V[Path(path)]; !ok {
		panic("content not exist: " + path)
	}

	c.V[Path(path)] = URL(url)
}

func (c Contents) Delete(path string) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if _, ok := c.V[Path(path)]; !ok {
		panic("content not exist: " + path)
	}

	c.V[Path(path)] = URL("")
}

func (c Contents) DeleteDir(path string) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if _, ok := c.V[Path(path)]; !ok {
		panic("content not exist: " + path)
	}

	delete(c.V, Path(path))
}
