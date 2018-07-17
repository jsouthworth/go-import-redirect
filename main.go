package main

import (
	"flag"
	"fmt"
	"jsouthworth.net/go/go-import-redirector/godoc"
	"log"
	"net/http"
	"net/http/fcgi"
	"os"
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
	log.SetPrefix("go-import-redirect: ")
	flag.Parse()
	if flag.NArg() != 3 {
		flag.Usage()
	}
	importPath := flag.Arg(0)
	repoPath := flag.Arg(1)
	servePath := flag.Arg(2)
	mux := http.NewServeMux()
	mux.Handle(servePath, godoc.Redirect(vcs, importPath, repoPath))
	if fastcgi {
		log.Fatal(fcgi.Serve(nil, mux))
	} else {
		log.Fatal(http.ListenAndServe(addr, mux))
	}
}
