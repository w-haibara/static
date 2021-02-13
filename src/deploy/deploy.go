package deploy

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

	tmpRootDir, err := ioutil.TempDir("", strings.Replace("osoba-"+c.RootPath+c.Path, string(os.PathSeparator), "", -1))
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpRootDir)
	log.Println("tmp file dir:", tmpRootDir)

	log.Println("[zip expand starting]")
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			if err = os.MkdirAll(filepath.Join(tmpRootDir, f.Name), os.ModePerm); err != nil {
				return err
			}
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		buf := new(bytes.Buffer)
		io.Copy(buf, rc)

		if err := ioutil.WriteFile(filepath.Join(tmpRootDir, f.Name), buf.Bytes(), 0666); err != nil {
			return err
		}
	}
	log.Println("[zip expand success]")

	distPath := filepath.Join(c.RootPath, c.Path)
	if err := os.MkdirAll(distPath, os.ModePerm); err != nil {
		return err
	}
	os.Remove(distPath)
	if err := os.MkdirAll(distPath, os.ModePerm); err != nil {
		return err
	}

	files, err := ioutil.ReadDir(tmpRootDir)
	if err != nil {
		return err
	}

	parentDir := ""
	if len(files) == 1 {
		parentDir = files[0].Name()
	}

	log.Println("[deploy starting]")
	if err := dirCopyAll(filepath.Join(tmpRootDir, parentDir), distPath); err != nil {
		return err
	}
	log.Println("[deploy success]")

	return nil
}

func dirCopyAll(src, dst string) error {
	log.Println(src, "-->", dst)
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			log.Println("dir ", string(f.Name()+"/"))
			if err := os.MkdirAll(filepath.Join(dst, f.Name()), os.ModePerm); err != nil {
				return err
			}
			dirCopyAll(filepath.Join(src, f.Name()), filepath.Join(dst, f.Name()))
			continue
		}

		log.Println("file", f.Name())
		b, err := ioutil.ReadFile(filepath.Join(src, f.Name()))
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(filepath.Join(dst, f.Name()), b, 0666); err != nil {
			return err
		}
	}

	return nil
}
