package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"
	"unicode/utf8"

	"gitlab.com/golang-commonmark/markdown"
)

var md = markdown.New(markdown.XHTMLOutput(true))
var defaultTemplate string
var basePath, cybermondayBootstrapRef, cybermondayTitle string
var mdTemplate *template.Template
var envs = map[string]string{}

func init() {
	var env string
	for _, env = range os.Environ() {
		var parts = strings.Split(env, "=")
		envs[parts[0]] = strings.Join(parts[1:], "=")
	}
	basePath = os.Getenv("CYBERMONDAY_BASEPATH")
	if utf8.RuneCountInString(basePath) == 0 {
		basePath = "."
	}
	cybermondayTitle = os.Getenv("CYBERMONDAY_TITLE")
	if utf8.RuneCountInString(cybermondayTitle) == 0 {
		cybermondayTitle = "Home"
	}
	cybermondayBootstrapRef = os.Getenv("CYBERMONDAY_BOOTSTRAP_REF")
	if utf8.RuneCountInString(cybermondayBootstrapRef) == 0 {
		cybermondayBootstrapRef = "//stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css"
	}
	defaultTemplate = os.Getenv("CYBERMONDAY_TEMPLATE")
	if utf8.RuneCountInString(defaultTemplate) == 0 {
		defaultTemplate = `<!DOCTYPE html>
<html>
<head>
<link rel="stylesheet" href="{{.Bootstrap}}">
<style type="text/css">
main > div.container {
	margin-top: 30px;
}
</style>
</head>
<body>
<header>
	<div class="navbar navbar-dark bg-dark shadow-sm">
    <div class="container d-flex justify-content-between">
      <a href="/" class="navbar-brand d-flex align-items-center">{{ .Title }}</a>
    </div>
	</div>
</header>
<main>
	<div class="container">{{ .Content }}</div>
</main>
</body>
</html>`
	}
}

type templateData struct {
	Bootstrap string
	Content   string
	Title     string
	Env       map[string]string
}

func main() {
	var err error
	mdTemplate, err = template.New("template").Parse(defaultTemplate)
	if err != nil {
		log.Panic(err)
	}
	var s = &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: &handler{},
	}
	err = s.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var filepath = normalizePath(req.URL.EscapedPath())
	if !isMarkdownRequest(filepath) {
		w.WriteHeader(400)
		w.Write([]byte("400 invalid request"))
		return
	}
	var data []byte
	var err error
	if !fileExists(filepath) {
		w.WriteHeader(404)
		w.Write([]byte("404 not found"))
		return
	}
	data, err = ioutil.ReadFile(filepath)
	if err != nil {
		w.Header().Add("Content-type", "text/plain")
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	var d = templateData{
		Bootstrap: cybermondayBootstrapRef,
		Title:     cybermondayTitle,
		Content:   md.RenderToString(data),
		Env:       envs,
	}
	err = mdTemplate.Execute(w, d)
	if err != nil {
		w.Header().Add("Content-type", "text/plain")
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	//w.Write([]byte(defaultTemplate + content + cybermondayFooter))
}

func isMarkdownRequest(name string) bool {
	return strings.HasSuffix(name, ".md")
}

func normalizePath(originalPath string) string {
	var filePath = path.Join(basePath, originalPath)
	if isDirectory(filePath) {
		if fileExists(path.Join(filePath, "index.md")) {
			filePath = path.Join(filePath, "index.md")
		} else if fileExists(path.Join(filePath, "README.md")) {
			filePath = path.Join(filePath, "README.md")
		}
	}
	return filePath
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func isDirectory(name string) bool {
	var s, err = os.Stat(name)
	if err != nil {
		return false
	}
	return s.IsDir()
}
