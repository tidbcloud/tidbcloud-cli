package util

import (
	"github.com/spf13/afero"
)

func RemoveFile(fileName string) error {
	fs := afero.NewOsFs()
	stat, err := fs.Stat(fileName)
	if err != nil {
		return err
	}
	if stat != nil {
		err := fs.Remove(fileName)
		if err != nil {
			return err
		}
	}

	return nil
}
