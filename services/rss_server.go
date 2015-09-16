package services

import (
	"github.com/gdg-belfast/gross/controllers"

	"net/http"
)

const (
	RSS_FEED = "/"
	FILE_URL = "/file/"
)

type Rss struct {
	Router *http.ServeMux
	Server controllers.Rss
}

func (r *Rss) SetupMux() {
	r.Router.HandleFunc(RSS_FEED, r.Server.Feed)
	r.Router.HandleFunc(FILE_URL, r.Server.ServeFile)
}
