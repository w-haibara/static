package osoba

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (a App) Deploy(path string) error {
	dirName := filepath.Join(a.DocumentRoot, path)

	// make directory
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
			return err
		}
	}

	// fetch zip file
	var r *http.Response
	if err := func() error {
		a.Contents.Mu.RLock()
		defer a.Contents.Mu.RUnlock()
		var err error
		r, err = http.Get(a.Contents.V[path].URL)
		if err != nil {
			return err
		}
		return nil
	}(); err != nil {
		return err
	}
	defer r.Body.Close()

	// create tmpolary directory
	tmpDir, err := os.MkdirTemp("", strings.Replace("osoba"+path, string(os.PathSeparator), "", -1))
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	// create zip file
	zipPath := filepath.Join(tmpDir, "donwload.zip")
	f, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	io.Copy(f, r.Body)

	// expand zip file
	if err := a.unzip(zipPath, tmpDir); err != nil {
		return err
	}

	// copy contents from tmpolary directory to document root
	dirCopyAll(filepath.Join(tmpDir, a.TmpDirContentsPrefix), a.DocumentRoot)

	return nil
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
