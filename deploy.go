package osoba

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (a App) DeployAll() {
	for k, _ := range a.Contents.V {
		a.Deploy(k)
	}
}

func (a App) Deploy(path Path) {
	dirName := filepath.Join(a.DocumentRoot, string(path))

	// make directory
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
			panic(err.Error())
		}
	}

	// fetch zip file
	a.Contents.Mu.Lock()
	r, err := http.Get(string(a.Contents.V[path]))
	if err != nil {
		panic(err.Error())
	}
	a.Contents.Mu.Unlock()
	defer r.Body.Close()

	// create tmpolary directory
	tmpDir, err := os.MkdirTemp("", strings.Replace("osoba"+string(path), string(os.PathSeparator), "", -1))
	if err != nil {
		panic(err.Error())
	}
	defer os.RemoveAll(tmpDir)

	// create zip file
	zipPath := filepath.Join(tmpDir, "donwload.zip")
	f, err := os.Create(zipPath)
	if err != nil {
		panic(err.Error())
	}
	io.Copy(f, r.Body)

	// expand zip file
	if err := a.unzip(zipPath, tmpDir); err != nil {
		panic(err.Error())
	}

	// copy contents from tmpolary directory to document root
	dirCopyAll(filepath.Join(tmpDir, a.TmpDirContentsPrefix), a.DocumentRoot)
}

func (a App) unzip(source, dir string) error {
	r, err := zip.OpenReader(source)
	if err != nil {
		return err
	}

	prefix := strings.Split(r.File[0].Name, string(os.PathSeparator))[0]
	for _, file := range r.File {
		tmp := strings.Split(file.Name, string(os.PathSeparator))[0]
		if tmp != prefix {
			prefix = ""
			break
		}
	}

	for _, file := range r.File {
		if file.Mode().IsDir() {
			continue
		}
		os.MkdirAll(filepath.Join(dir, filepath.Dir(file.Name)), os.ModePerm)
		rc, err := file.Open()
		if err != nil {
			return err
		}
		f, err := os.Create(filepath.Join(dir, file.Name))
		if err != nil {
			return err
		}
		f.ReadFrom(rc)
	}

	if prefix != "" {
		os.Rename(filepath.Join(dir, prefix), filepath.Join(dir, a.TmpDirContentsPrefix))
	}

	return nil
}
