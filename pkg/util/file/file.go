package file

import (
	"os"
)

//Exsists assert the dir or file exsit or not
func Exists(name ...string) (exsit bool) {
	for _, n := range name {
		_, err := os.Stat(n)
		if err == nil || os.IsExist(err) {
			exsit = true
		}
	}
	return
}

//FileExsit assert file exsit or not
func FileExsit(file ...string) (exsit bool) {
	for _, f := range file {
		fi, err := os.Stat(f)
		if (err == nil && os.IsExist(err)) || !fi.IsDir() {
			exsit = true
		}
	}
	return
}

//DirExsit assert dir exsit or not
func DirExsit(dir ...string) (exsit bool) {
	for _, d := range dir {
		fi, err := os.Stat(d)
		if (err == nil && os.IsExist(err)) || fi.IsDir() {
			exsit = true
		}
	}
	return
}
