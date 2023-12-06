package handler

import (
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/gkawamoto/cybermonday/config"
	"gitlab.com/golang-commonmark/markdown"
)

type templateData struct {
	Content string
	Env     map[string]string
}

func New(staticHandler http.Handler, conf *config.Config) (http.Handler, error) {
	tplt, err := loadTemplate(conf)
	if err != nil {
		return nil, err
	}

	md := markdown.New(markdown.XHTMLOutput(true))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filepath := normalizePath(conf, r.URL.EscapedPath())
		if filepath == "" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 not found"))
			return
		}

		if !isMarkdownRequest(filepath) {
			staticHandler.ServeHTTP(w, r)
			return
		}

		if !fileExists(filepath) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 not found"))
			return
		}

		data, err := os.ReadFile(filepath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		d := templateData{
			Content: md.RenderToString(data),
			Env:     conf.Envs,
		}
		if err := tplt.Execute(w, d); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}), nil
}

func loadTemplate(conf *config.Config) (*template.Template, error) {
	result := conf.Template
	if result == "" {
		defaultTemplatePath := conf.TemplatePath
		if defaultTemplatePath == "" {
			defaultTemplatePath = conf.DefaultTemplatePath
		}

		log.Printf("Loading template from %s", defaultTemplatePath)

		data, err := os.ReadFile(defaultTemplatePath)
		if err != nil {
			return nil, err
		}

		result = string(data)
	}
	return template.New("").Parse(result)
}

func isMarkdownRequest(name string) bool {
	return strings.HasSuffix(name, ".md")
}

func normalizePath(conf *config.Config, originalPath string) string {
	filePath := path.Join(conf.BasePath, originalPath)
	if !isDirectory(filePath) {
		return filePath
	}

	for _, index := range []string{"index.md", "README.md"} {
		if !fileExists(path.Join(filePath, index)) {
			continue
		}

		return path.Join(filePath, index)
	}

	return ""
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
	s, err := os.Stat(name)
	if err != nil {
		return false
	}
	return s.IsDir()
}
