package process

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// type Posit struct {
// 	// text     string
// 	data     []byte
// 	filename string
// }

// func loadPosit(path string) (posit *Posit, err error) {
// 	dat, err := ioutil.ReadFile(path)
// 	if err != nil {
// 		return
// 	}

// 	posit.data = dat
// 	// posit.filename = filepath.Base(path)
// 	posit.filename = path

// 	return
// }

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, errors.New(fmt.Sprintf("%s: %v", path, err))
	}
	return fileInfo.IsDir(), err
}

func findPosits(path string) ([]string, error) {

	// determine if path is dir
	isDir, err := isDirectory(path)
	if err != nil {
		return nil, err
	}

	if isDir {
		positPaths, err := filepath.Glob(fmt.Sprintf("%s/*_emlbody.txt", path))
		if err != nil {
			return nil, err
		}

		return positPaths, nil
	} else {
		return nil, fmt.Errorf("%s is not a directory", path)
	}
}
