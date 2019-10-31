package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

var indexPath string

func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	var urlPath = r.URL.Path[1:]
	fmt.Println(urlPath)
	if urlPath == "" {
		http.ServeFile(w, r, fmt.Sprintf("./%s", indexPath))
	} else {
		http.ServeFile(w, r, urlPath)
	}
}

func main() {
	port := *flag.Int("port", 3000, "Port for Shanty to serve on")
	flag.Parse()

	if flag.Arg(0) != "" {
		_indexPath, err := filepath.Rel(filepath.Dir(os.Args[0]), flag.Arg(0))
		indexPath = _indexPath
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("index= %s\n", indexPath)
		}
		http.HandleFunc("/", rootHandler)
	} else {
		http.Handle("/", http.FileServer(http.Dir("./")))
	}

	fmt.Printf("Listening on 127.0.0.1:%d...\n", port)
	go http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil)

	fmt.Println("Press Crtl+C to Exit")

	if flag.Arg(0) != "" {
		time.Sleep(time.Second)
		openURL(fmt.Sprintf("http://localhost:%d/", port))
	}
	select {}
}
