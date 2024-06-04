package util

import (
	"fmt"
	"net/mail"
	"os"
	"path/filepath"
	"runtime"
)

func AppDir() (string, error) {
	_, fileStr, _, _ := runtime.Caller(0)
	dir := filepath.Dir(fileStr)

	for {
		files, err := os.ReadDir(dir)
		if err != nil {
			return "", fmt.Errorf("unable to read files")
		}

		for _, f := range files {
			if f.Name() == "go.mod" {
				return dir, nil
			}
		}

		dir = filepath.Dir(dir)
	}
}

func DaysInMilliseconds(days int64) int64 {
	return days * 24 * 60 * 60 * 1000
}

func ArrayContainsString(array []string, element string) bool {
	for _, el := range array {
		if el == element {
			return true
		}
	}
	return false
}

func EmailIsValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err != nil
}
