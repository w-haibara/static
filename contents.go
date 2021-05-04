package osoba

type Path string
type URL string
type Contents map[Path]URL

func (c Contents) Create(path, url string) {
	if _, ok := c[Path(path)]; ok {
		panic("content exist: " + path)
	}
	c[Path(path)] = URL(url)
}

func (c Contents) Update(path, url string) {
	if _, ok := c[Path(path)]; !ok {
		panic("content not exist: " + path)
	}
	c[Path(path)] = URL(url)
}

func (c Contents) Delete(path string) {
	if _, ok := c[Path(path)]; !ok {
		panic("content not exist: " + path)
	}
	c[Path(path)] = URL("")
}

func (c Contents) DeleteDir(path string) {
	if _, ok := c[Path(path)]; !ok {
		panic("content not exist: " + path)
	}
	delete(c, Path(path))
}
