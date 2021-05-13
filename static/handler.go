package static

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Handler struct{}

func suffixed(s string, l ...string) bool {
	for _, p := range l {
		if strings.HasSuffix(s, p) {
			return true
		}
	}
	return false
}

func headers(w http.ResponseWriter, path string) {
	if suffixed(path, ".html", ".html.gz") {
		w.Header().Set("Content-Type", "text/html")
	}
	if suffixed(path, ".css", ".css.gz") {
		w.Header().Set("Content-Type", "text/css")
	}
	if suffixed(path, ".js", ".js.gz") {
		w.Header().Set("Content-Type", "application/json")
	}

	if suffixed(path, ".gz") {
		w.Header().Set("Content-Encoding", "gzip")
	}
}

func fspath(req string) string {
	if req == "/" {
		return "/index.html"
	}
	if suffixed(req, "/") {
		return req + "index.html"
	}
	if suffixed(req, ".gz", ".html", ".css", ".js", ".png", ".gif", ".jpg") {
		return req
	}
	return req + "/index.html"
}

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		path := fspath(req.URL.Path)

		alt := os.Getenv("OPTIGIT_UI_ROOT")
		if len(Assets) == 0 || alt != "" {
			if alt == "" {
				alt = "assets"
			}
			alt = strings.TrimSuffix(alt, "/")
			b, err := ioutil.ReadFile(fmt.Sprintf("%s%s", alt, path))
			if err == nil {
				headers(w, path)
				w.WriteHeader(200)
				w.Write(b)
				return
			}
			w.WriteHeader(404)
			fmt.Fprintf(w, "404 not found\n")
			return
		}

		if b, ok := Assets[path]; ok {
			headers(w, path)
			w.WriteHeader(200)
			w.Write(b)
			return
		}
		w.WriteHeader(404)
		fmt.Fprintf(w, "404 not found\n")
	}
}
