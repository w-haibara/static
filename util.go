package osoba

import (
	"os"
	"path/filepath"
)

func dirCopyAll(src, dst string) error {
	files, err := os.ReadDir(src)
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

		b, err := os.ReadFile(filepath.Join(src, f.Name()))
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dst, f.Name()), b, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}
