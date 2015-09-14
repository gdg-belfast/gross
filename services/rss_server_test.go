package services

import (
	. "github.com/smartystreets/goconvey/convey"

	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRss(t *testing.T) {

	Convey(`Given I have an RSS server instance`, t, func() {

		var validFile, serveFileCalled bool
		rssFeed := `<?xml version="1.0"?><feed></feed>`
		validFileContents := "validFileContents"

		mockRss := &mockRssController{
			MockFeed: func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, rssFeed)
			},
			MockServe: func(w http.ResponseWriter, r *http.Request) {
				serveFileCalled = true
				if validFile {
					fmt.Fprintf(w, validFileContents)
				} else {
					http.Error(w, "", http.StatusNotFound)
				}
			},
		}

		mux := http.NewServeMux()
		server := &Rss{
			Router: mux,
			Server: mockRss,
		}

		Convey(`And I have set up the routes`, func() {
			server.SetupMux()

			Convey(`When I go to the bare URL`, func() {
				req, err := http.NewRequest("GET", "http://localhost", nil)
				if err != nil {
					panic(err)
				}

				w := httptest.NewRecorder()
				server.Router.ServeHTTP(w, req)

				Convey(`Then I will get a 301 status code`, func() {
					So(w.Code, ShouldEqual, http.StatusMovedPermanently)
				})

				Convey(`Then the Location: header will be /`, func() {
					So(w.HeaderMap["Location"][0], ShouldEqual, buildUrl(`/`))
				})
			})

			Convey(`When I go to the RSS_FEED url`, func() {
				req, err := http.NewRequest("GET", buildUrl(`/`), nil)
				if err != nil {
					panic(err)
				}

				w := httptest.NewRecorder()
				server.Router.ServeHTTP(w, req)

				Convey(`Then I will get the RSS feed`, func() {
					bodyText, err := ioutil.ReadAll(w.Body)
					if err != nil {
						panic(err)
					}
					So(string(bodyText), ShouldEqual, rssFeed)
				})
			})

			Convey(`When I access a FILE_URL`, func() {

				Convey(`And the file URL is valid`, func() {
					validFile = true
					serveFileCalled = false

					req, err := http.NewRequest("GET", buildUrl(FILE_URL, "valid"), nil)
					if err != nil {
						panic(err)
					}

					w := httptest.NewRecorder()
					server.Router.ServeHTTP(w, req)

					Convey(`Then the ServeFile function will have been called`, func() {
						So(serveFileCalled, ShouldBeTrue)
					})

					Convey(`Then the status code will be 200`, func() {
						So(w.Code, ShouldEqual, http.StatusOK)
					})

					Convey(`Then the file contents will be returned`, func() {
						bodyText, err := ioutil.ReadAll(w.Body)
						if err != nil {
							panic(err)
						}
						So(string(bodyText), ShouldEqual, validFileContents)
					})
				})

				Convey(`And the file URL is invalid`, func() {
					validFile = false
					serveFileCalled = false

					req, err := http.NewRequest("GET", buildUrl(FILE_URL, "valid"), nil)
					if err != nil {
						panic(err)
					}

					w := httptest.NewRecorder()
					server.Router.ServeHTTP(w, req)

					Convey(`Then the ServeFile function will have been called`, func() {
						So(serveFileCalled, ShouldBeTrue)
					})

					Convey(`Then the status code will be 404`, func() {
						So(w.Code, ShouldEqual, http.StatusNotFound)
					})
				})
			})
		})
	})
}

type mockRssController struct {
	MockFeed  func(http.ResponseWriter, *http.Request)
	MockServe func(http.ResponseWriter, *http.Request)
}

func (stub *mockRssController) Feed(w http.ResponseWriter, r *http.Request) {
	stub.MockFeed(w, r)
}

func (stub *mockRssController) ServeFile(w http.ResponseWriter, r *http.Request) {
	stub.MockServe(w, r)
}

func buildUrl(path ...string) string {
	f := fmt.Sprintf("http://localhost%s", strings.Join(path, ``))
	//Println(f)
	return f
}
