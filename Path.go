package utils

import (
	"runtime"
	"os/exec"
	"path/filepath"
	"os"
	"strings"
)

func GetCurrPath() string {
	if runtime.GOOS == "windows" {
		file, _ := exec.LookPath(os.Args[0])
		path, _ := filepath.Abs(file)
		path = path[0:strings.LastIndex(path,"\\")] + "\\"
		return path
	} else{
		return "./"
	}

}