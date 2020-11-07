package uploadserver

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/sethpyle376/cs-statman/processor/matchprocessor"
)

type UploadServer struct {
	mp *matchprocessor.MatchProcessor
}

func New() (*UploadServer, error) {
	us := &UploadServer{}
	mp, err := matchprocessor.New()
	us.mp = mp
	return us, err
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

	tempfile, err := ioutil.TempFile("temp", "upload-*.dem")

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

	err = us.mp.ProcessMatch(tempfile)
	if err != nil {
		panic("error processing match")
	}
}

func (us *UploadServer) Start() {
	http.HandleFunc("/upload", us.handleUpload)
	uploadPort, ok := os.LookupEnv("UPLOAD_PORT")

	if !ok {
		uploadPort = "8081"
	}

	http.ListenAndServe(":"+uploadPort, nil)
}

func (us *UploadServer) Stop() { /* TODO */ }
