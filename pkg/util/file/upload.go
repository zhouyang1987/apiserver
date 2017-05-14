package file

import (
	"io"
	"net/http"
	"os"

	"apiserver/pkg/util/log"
)

/*func Upload(w http.ResponseWriter, r *http.Request, filePath string) {
	src, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.Copy(dst, src)
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, `{"err"：0,"msg":"upload success","status":"success"}`)
}*/

func Upload(r *http.Request, fileDir string) error {
	src, head, err := r.FormFile("file")
	log.Debugf("err == %v", err)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(fileDir, os.ModePerm); err != nil {
		log.Errorf("create dir %v err: %v", fileDir, err)
		return err
	}
	defer src.Close()
	filePath := fileDir + string(os.PathSeparator) + head.Filename
	log.Debugf("dst == %v", filePath)
	dst, err := os.Create(filePath)
	log.Debugf("err == %v", err)
	if err != nil {
		log.Errorf("upload file %v err: %v", head.Filename, err)
		return err
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		log.Debugf("err == %v", err)
		log.Errorf("upload file %v err: %v", head.Filename, err)
		return err
	}
	return nil
}
