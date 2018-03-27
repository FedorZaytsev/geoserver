package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/hishamkaram/uploadShapeFile/geoserver"
)

var uploadedPath = "./uploaded/"

var gsCatalog geoserver.GeoServer

func handleUploaded(file *bytes.Buffer, filename string) string {
	_ = os.Mkdir(uploadedPath, 0700)
	filepath := uploadedPath + filename
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	io.Copy(f, file)
	return filepath
}
func index(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		file, handler, err := r.FormFile("fileupload")
		defer file.Close()
		if err != nil {
			panic(err)
		}
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			panic(err)
		}
		uploadedPath := handleUploaded(buf, handler.Filename)
		fileLocation, _ := filepath.Abs(uploadedPath)
		success, _ := gsCatalog.UploadShapeFile(fileLocation, "")
		fmt.Println(success)
	}
	tmplt := template.New("home.html")
	tmplt, _ = tmplt.ParseFiles("templates/home.html")

	tmplt.Execute(w, gsCatalog)
}
func main() {
	fileLocation, _ := filepath.Abs("./config.yml")
	gsCatalog.LoadConfig(fileLocation)
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static/").Handler(s)
	http.Handle("/", r)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
