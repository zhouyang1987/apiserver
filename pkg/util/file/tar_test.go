package file

import (
	"testing"
)

func TestTar(t *testing.T) {
	err := Tar("test.tar.gz", true, "file.go", "upload.go", "tar_test.go")
	if err != nil {
		t.Error(err)
	}
}

func TestUnTar(t *testing.T) {
	err := UnTar("test.tar.gz", "../")
	if err != nil {
		t.Error(err)
	}
}
