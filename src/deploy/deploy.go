package deploy

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Config struct {
	Path       string
	RootPath   string
	TmpPath    string
	ReleaseURL string
}

func (c Config) Deploy() error {
	c.TmpPath = "/tmp"

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

	for _, f := range zr.File {
		fpath := filepath.Join(c.TmpPath, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(c.TmpPath)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			if err = os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		defer outFile.Close()

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}
	}

	os.Remove(c.Path)

	distPath := filepath.Join(c.RootPath, c.Path)
	if err := os.MkdirAll(distPath, os.ModePerm); err != nil {
		return err
	}
	os.Remove(distPath)

	for _, f := range zr.File {
		fpath := filepath.Join(c.RootPath, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(c.RootPath)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			if err = os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		defer outFile.Close()

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}
	}

	/*
		if err := os.Rename(filepath.Join(c.TmpPath, "dist"), distPath); err != nil {
			return err
		}
	*/
	if err := exec.Command("mv", filepath.Join(c.TmpPath, "dist"), distPath).Run(); err != nil {
		return err
	}

	return nil
}
