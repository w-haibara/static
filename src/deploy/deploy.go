package deploy

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type Config struct {
	Path       string
	RootPath   string
	TmpPath    string
	ReleaseURL string
}

func (c Config) Deploy() error {
	r, err := http.Get(c.ReleaseURL)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	zr, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return err
	}

	tmpRootDir, err := ioutil.TempDir("", "example")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpRootDir)
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			if err = os.MkdirAll(filepath.Join(tmpRootDir, f.Name), os.ModePerm); err != nil {
				return err
			}
			continue
		}

		rc, err := f.Open()
		if err != nil {
			log.Println("zip file open err", err)
			panic(err)
		}
		defer rc.Close()

		buf := new(bytes.Buffer)
		io.Copy(buf, rc)

		if err := ioutil.WriteFile(filepath.Join(tmpRootDir, f.Name), buf.Bytes(), 0666); err != nil {
			panic(err)
		}
	}

	distPath := filepath.Join(c.RootPath, c.Path)
	if err := os.MkdirAll(distPath, os.ModePerm); err != nil {
		log.Println("MkdirAll err", err)
		return err
	}
	os.Remove(distPath)

	/*
		if err := os.Rename(filepath.Join(tmpRootDir, "dist"), distPath); err != nil {
			return err
		}
	*/

	if err := exec.Command("mv", filepath.Join(tmpRootDir, "dist"), distPath).Run(); err != nil {
		return err
	}

	return nil
}
