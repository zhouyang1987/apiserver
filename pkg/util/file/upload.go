package file

import (
	"io"
	"net/http"
	"os"
)

func Upload(w http.ResponseWriter, r *http.Request, filePath string) {
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
}
