package server

import (
	"fmt"
	"github.com/gobelfast/gross/mediafile"
	"github.com/gorilla/feeds"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// RssServer is a HTTP server which maintains and presents a https://github.com/gorilla/feeds
// feed
//
// Exposes:
//    / - Prints the RSS XML
//    /file/<id>/filename - Downloads the file
type RssServer struct {
	// URL of the Server
	Server string
	// Port the server listens on
	Port int
	// Feed Object
	Feed *feeds.RssFeed
	// Internal map of files that we are making available
	Filemap map[string]*mediafile.File
}

// NewServer takes the URL and port that the server will listen and
// offer its wares from.
// It returns a pointer to a RssServer
func NewServer(serverName string, port int) *RssServer {
	s := &RssServer{
		Server:  serverName,
		Port:    port,
		Feed:    &feeds.RssFeed{},
		Filemap: make(map[string]*mediafile.File),
	}
	return s
}

// SetFileInput sets the input channel that will feed the RSS feed
// Starts a new go routine
func (s *RssServer) SetFileInput(additions chan *mediafile.File) {
	go func() {
		for {
			newFile := <-additions
			s.Filemap[newFile.Hash] = newFile
			item := s.CreateRssItem(newFile)
			s.Feed.Items = append(s.Feed.Items, item)
			s.Feed.PubDate = time.Now().String()
		}
	}()
}

// Run starts the HTTP listener.
// Returns whatever error the HTTP server may raise
func (s *RssServer) Run() error {
	http.HandleFunc("/file/", s.GetFile)
	http.HandleFunc("/", s.GetList)

	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil)
}

// GetList is the endpoint that will return the XML file
func (s *RssServer) GetList(w http.ResponseWriter, r *http.Request) {
	feeds.WriteXML(s.Feed, w)
}

// GetFile is the endpoint that will return the file requested, or raise
// a 404 error if the file cannot be found
func (s *RssServer) GetFile(w http.ResponseWriter, r *http.Request) {
	keyParts := strings.Split(r.URL.Path, "/")
	// ["", files, "file hash", "file name"]
	if len(keyParts) != 4 {
		InvalidUrl(w)
		return
	}
	if file, ok := s.Filemap[keyParts[2]]; ok {
		log.Println("Getting", file.Filepath)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Name()))
		http.ServeFile(w, r, file.Filepath)
	} else {
		InvalidUrl(w)
		return
	}
}

// InvalidUrl writes a 404 response and error message to the provided HTTP ResponseWriter
func InvalidUrl(w http.ResponseWriter) {
	http.Error(w, "Invalid URL", http.StatusBadRequest)
}

// CreateRssItem takes a provided file.MediaFile and converts it into a RssItem
func (s *RssServer) CreateRssItem(file *mediafile.File) *feeds.RssItem {
	item := &feeds.RssItem{}
	item.Title = file.Name()
	item.Link = s.MakeLinkUrl(file.Hash, file.Name())
	item.Description = ""

	item.Enclosure = &feeds.RssEnclosure{
		Url:    s.MakeLinkUrl(),
		Length: fmt.Sprintf("%d", file.Size()),
		Type:   mime.TypeByExtension(filepath.Ext(file.Filepath)),
	}
	return item
}

// MakeLinkUrl creates a URL from the provided string parts
func (s *RssServer) MakeLinkUrl(parts ...string) string {
	link := fmt.Sprintf("%s/file/", s.ServerAddress())
	link += strings.Join(parts, "/")
	return link
}

// ServerAddress returns the base URL
func (s *RssServer) ServerAddress() string {
	return fmt.Sprintf("%s:%d", s.Server, s.Port)
}
