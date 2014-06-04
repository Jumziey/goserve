package main

import (
	"fmt"
	"flag"
	"net/http"
	"log"
	"io"
	"os"
)

var port = flag.String("p", "8080", "Sets server port")

func createFileHandler(path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w,r, path)
	}
}

func createFileIndex(files []string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for _,v := range(files) {
		io.WriteString(w, "<a href="+v+" download>"+v+"</a>\n")
		}
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	
	if err == nil {
		f, err := os.Open(path) //open
		if err != nil {
			fmt.Println(err)
			return false
		}
		defer f.Close()
		fi, err := f.Stat()		//get file info
		if err != nil {
			fmt.Println(err)
			return false
		}
		if !fi.IsDir() {	//is file?
			return true
		}
	}
	return false
}

func main() {
	flag.Parse()
	var files []string
	
	for _,v := range(flag.Args()) {
		if fileExists(v){
			files = append(files, v)
			http.HandleFunc(v, createFileHandler(v))
		} else {
		fmt.Println("Can't find the file: "+v,"\tSkipping...")
		}
	}
	fmt.Println(files)
	fmt.Println("port: ", *port)
	
	fileIndex := createFileIndex(files)
	http.HandleFunc("/", fileIndex)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
