package main

import (
	"bytes"
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
var basePath string
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
	loadTemplate()
}

func loadTemplate() {
	var defaultTemplate string
	var err error
	var result = os.Getenv("CYBERMONDAY_TEMPLATE")
	if utf8.RuneCountInString(defaultTemplate) == 0 {
		var defaultTemplatePath = os.Getenv("CYBERMONDAY_TEMPLATE_PATH")
		if utf8.RuneCountInString(defaultTemplatePath) == 0 {
			defaultTemplatePath = os.Getenv("CYBERMONDAY_DEFAULT_TEMPLATE_PATH")
			if utf8.RuneCountInString(defaultTemplatePath) == 0 {
				defaultTemplatePath = "./resources/index.tplt.html"
			}
		}
		var data []byte
		data, err = os.ReadFile(defaultTemplatePath)
		if err != nil {
			log.Panic(err)
		}
		result = string(data)
	}
	defaultTemplate = result
	mdTemplate, err = template.New("template").Parse(defaultTemplate)
	if err != nil {
		log.Panic(err)
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
	var s = &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: &handler{http.FileServer(http.Dir(basePath))},
	}
	err = s.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

type handler struct {
	staticHandler http.Handler
}

func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var filepath = normalizePath(request.URL.EscapedPath())
	if !isMarkdownRequest(filepath) {
		h.staticHandler.ServeHTTP(writer, request)
		return
		//writer.WriteHeader(400)
		//writer.Write([]byte("400 invalid request"))
		//return
	}
	var data []byte
	var err error
	if !fileExists(filepath) {
		writer.WriteHeader(404)
		writer.Write([]byte("404 not found"))
		return
	}
	data, err = os.ReadFile(filepath)
	if err != nil {
		writer.Header().Add("Content-type", "text/plain")
		writer.WriteHeader(500)
		writer.Write([]byte(err.Error()))
		return
	}
	var d = templateData{
		Content: md.RenderToString(data),
		Env:     envs,
	}
	var buffer bytes.Buffer
	err = mdTemplate.Execute(&buffer, d)
	if err != nil {
		writer.Header().Add("Content-type", "text/plain")
		writer.WriteHeader(500)
		writer.Write([]byte(err.Error()))
		return
	}
	//log.Println(buffer.String())
	writer.Write(buffer.Bytes())
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
