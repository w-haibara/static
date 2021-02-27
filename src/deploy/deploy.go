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

type Info struct {
	Path       string
	RootPath   string
	ReleaseURL string
}

func AwaitDeploy(info chan Info) {
	for {
		var msg string = ""
		select {
		case i := <-info:
			msg += "deploy to " + string(i.Path) + ", "
			if err := i.Deploy(); err != nil {
				msg += "error: " + err.Error()
			}
			msg += "success"
		}
		log.Println(msg)
	}
}

func (i Info) Deploy() error {
	r, err := http.Get(i.ReleaseURL)
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

	tmpRootDir, err := ioutil.TempDir("", strings.Replace("osoba-"+i.RootPath+i.Path, string(os.PathSeparator), "-", -1))
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
			return err
		}
		defer rc.Close()

		buf := new(bytes.Buffer)
		io.Copy(buf, rc)

		if err := ioutil.WriteFile(filepath.Join(tmpRootDir, f.Name), buf.Bytes(), 0666); err != nil {
			return err
		}
	}

	distPath := filepath.Join(i.RootPath, i.Path)
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

	if err := dirCopyAll(filepath.Join(tmpRootDir, parentDir), distPath); err != nil {
		return err
	}

	return nil
}

func dirCopyAll(src, dst string) error {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			if err := os.MkdirAll(filepath.Join(dst, f.Name()), os.ModePerm); err != nil {
				return err
			}
			if err := dirCopyAll(filepath.Join(src, f.Name()), filepath.Join(dst, f.Name())); err != nil {
				return err
			}
			continue
		}

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

func (i Info) Delete() error {
	log.Println("rm:", filepath.Join(i.RootPath, i.Path))
	return os.RemoveAll(filepath.Join(i.RootPath, i.Path))
}
