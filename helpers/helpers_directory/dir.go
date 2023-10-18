package helpers_directory

import (
	"os"
	"path/filepath"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func DirExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
func Dir(path string, upLevel int) string {
	for i := 0; i < upLevel; i++ {
		path = filepath.Dir(path)
	}
	return path
}
func Getwd(upLevel int) string {
	wd, _ := os.Getwd()
	return Dir(wd, upLevel)
}
