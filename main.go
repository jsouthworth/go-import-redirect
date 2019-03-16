package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/fcgi"
	"os"

	"jsouthworth.net/go/go-import-redirector/godoc"
)

var (
	addr    string
	vcs     string
	fastcgi bool
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: go-import-redirect <import> <repo>\n")
	fmt.Fprintf(os.Stderr, "options:\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	flag.StringVar(&addr, "addr", ":http", "serve http on `address`")
	flag.StringVar(&vcs, "vcs", "git", "set version control `system`")
	flag.BoolVar(&fastcgi, "fcgi", false, "use fcgi to serve requests")
	flag.Usage = usage
}

func main() {
	var importPath string
	var repoPath string
	var servePath string
	log.SetPrefix("go-import-redirect: ")
	flag.Parse()
	if flag.NArg() == 3 {
		importPath = flag.Arg(0)
		repoPath = flag.Arg(1)
		servePath = flag.Arg(2)
	}
	if importPath == "" {
		importPath = os.Getenv("IMPORT_PATH")
	}
	if repoPath == "" {
		repoPath = os.Getenv("REPO_PATH")
	}
	if servePath == "" {
		servePath = os.Getenv("SERVE_PATH")
	}
	log.Println("Import Path:", importPath)
	log.Println("Repo Path:", repoPath)
	log.Println("Serve Path:", servePath)
	if importPath == "" || repoPath == "" || servePath == "" {
		flag.Usage()
	}
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	mux := http.NewServeMux()
	mux.Handle(servePath, godoc.Redirect(vcs, importPath, repoPath))
	if fastcgi {
		log.Fatal(fcgi.Serve(nil, mux))
	} else {
		log.Fatal(http.ListenAndServe(addr, mux))
	}
}
