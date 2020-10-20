package uploadserver

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type UploadServer struct{}

func New() (*UploadServer, error) {
	us := &UploadServer{}
	return us, nil
}

func (us *UploadServer) handleUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Uploading file\n")

	// 50Mb max
	r.ParseMultipartForm(50 << 21)
	file, handler, err := r.FormFile("myFile")

	if err != nil {
		fmt.Fprintf(w, "Error getting file\n")
		fmt.Fprintf(w, err.Error())
		return
	}

	defer file.Close()

	fmt.Fprintf(w, "File size: %d\n", handler.Size)

	tempfile, err := ioutil.TempFile("temp", "upload-*.png")

	if err != nil {
		fmt.Fprintf(w, "error getting temp file\n")
		fmt.Fprintf(w, err.Error()+"\n")
		return
	}
	defer tempfile.Close()

	filebytes, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Fprintf(w, "error reading file to bytes\n")
		return
	}

	tempfile.Write(filebytes)
	fmt.Fprintf(w, "successfully uploaded file\n")
}

func (us *UploadServer) Start() {
	http.HandleFunc("/upload", us.handleUpload)
	http.ListenAndServe(":8081", nil)
}

func (us *UploadServer) Stop() { /* TODO */ }
