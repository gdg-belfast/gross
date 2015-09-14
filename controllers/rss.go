package controllers

import (
	"net/http"
)

type Rss interface {
	Feed(http.ResponseWriter, *http.Request)
	ServeFile(http.ResponseWriter, *http.Request)
}

type rssHandler struct {
}

func (handler *rssHandler) Feed(w http.ResponseWriter, r *http.Request) {

}

func (handler *rssHandler) ServeFile(w http.ResponseWriter, r *http.Request) {

}
