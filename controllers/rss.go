package controllers

import (
	"github.com/gdg-belfast/gross/usecases"
	"github.com/gorilla/feeds"

	"fmt"
	"net/http"
	"path"
	"strings"
)

const (
	DEFAULT_TITLE       = "Gross RSS Feed"
	DEFAULT_DESCRIPTION = "A feed provided by GDG Belfast"
	DEFAULT_AUTHOR_NAME = "GDG Belfast"
)

type RssConfig struct {
	Title       string
	Link        string
	Description string
	AuthorName  string
}

type Rss interface {
	Feed(http.ResponseWriter, *http.Request)
	ServeFile(http.ResponseWriter, *http.Request)
}

type RssController struct {
	// URL of the Server
	ServerName string
	// Port the server listens to
	ServerPort int
	// Feed Object
	FullFeed *feeds.Feed
	// File provider
	File usecases.FileProvider
}

func (controller *RssController) Feed(w http.ResponseWriter, r *http.Request) {
	controller.FullFeed.WriteRss(w)
}

func (controller *RssController) ServeFile(w http.ResponseWriter, r *http.Request) {
	keyParts := strings.Split(r.URL.Path, "/")
	// ["", files, "file hash", "file name"]
	if len(keyParts) != 4 {
		http.Error(w, "Invalid URL", http.StatusNotFound)
		return
	}

	file, err := controller.File.Get("filehash", "filename")
	if err != nil {
		switch err {
		case err.(*usecases.ErrNotFound):
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", path.Base(file.Name())))
	http.ServeFile(w, r, file.Name())
}

// NewRssController takes the server name and port from which files will be presented
// as well as the config for the RSS feed itself.
func NewRssController(serverName string, serverPort int, config *RssConfig) *RssController {
	controller := &RssController{
		ServerName: serverName,
		ServerPort: serverPort,
		FullFeed:   &feeds.Feed{},
	}

	getValue := func(value, defaultValue string) string {
		if value != "" {
			return value
		}
		return defaultValue
	}

	controller.FullFeed.Title = getValue(config.Title, DEFAULT_TITLE)
	controller.FullFeed.Link = &feeds.Link{
		Href: getValue(config.Link, serverName),
	}
	controller.FullFeed.Description = getValue(config.Description, DEFAULT_DESCRIPTION)
	controller.FullFeed.Author = &feeds.Author{
		Name: getValue(config.AuthorName, DEFAULT_AUTHOR_NAME),
	}
	return controller
}
