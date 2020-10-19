package uploadserver

import (
	"fmt"
	"net/http"
)

type UploadServer struct{}

func New() (*UploadServer, error) {
	us := &UploadServer{}
	return us, nil
}

func (us *UploadServer) handleUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Uploading file")
}

func (us *UploadServer) Start() {
	http.HandleFunc("/upload", us.handleUpload)
	http.ListenAndServe(":8080", nil)
}

func (us *UploadServer) Stop() { /* TODO */ }
