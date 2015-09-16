package controllers

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	"github.com/gdg-belfast/gross/usecases"

	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
)

func TestNewRssController(t *testing.T) {

	Convey(`When I call NewRssController`, t, func() {
		serverName := "http://test.server"
		serverPort := 8080

		controller := NewRssController(serverName, serverPort, &RssConfig{AuthorName: "Tyndyll"})

		Convey(`Then I will get a rssController`, func() {
			So(controller, ShouldHaveSameTypeAs, &RssController{})

			Convey(`And the ServerName will be populated`, func() {
				So(controller.ServerName, ShouldEqual, serverName)
			})

			Convey(`And the ServerPort will be populated`, func() {
				So(controller.ServerPort, ShouldEqual, serverPort)
			})

			Convey(`And the Feed title will be populated with the default`, func() {
				So(controller.FullFeed.Title, ShouldEqual, DEFAULT_TITLE)
			})

			Convey(`And the Feed Link HREF will be populated with the server URL`, func() {
				So(controller.FullFeed.Link.Href, ShouldEqual, serverName)
			})

			Convey(`And the Feed Description will be populated with the default`, func() {
				So(controller.FullFeed.Description, ShouldEqual, DEFAULT_DESCRIPTION)
			})

			Convey(`And the Feed Author Name will be populated with the provided value`, func() {
				So(controller.FullFeed.Author.Name, ShouldEqual, "Tyndyll")
			})
		})
	})
}

func TestRssController(t *testing.T) {
	Convey(`Given I have a rssController`, t, func() {
		var err error
		serverName := "http://example.com"
		serverPort := 8080

		mockFile := &mockFileProvider{}
		controller := NewRssController(serverName, serverPort, &RssConfig{})
		controller.File = mockFile

		Convey(`When I call Feed`, func() {
			req, err := http.NewRequest("GET", serverName, nil)
			if err != nil {
				panic(err)
			}
			resp := httptest.NewRecorder()

			controller.Feed(resp, req)

			Convey(`Then I will receive feeds.XML output`, func() {
				bodyText, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					panic(err)
				}
				rssFeed, err := controller.FullFeed.ToRss()
				if err != nil {
					panic(err)
				}
				So(rssFeed, ShouldEqual, string(bodyText))
			})

			Convey(`Then the status code will be StatusOK`, func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
		})

		Convey(`When I call ServeFile`, func() {
			var req *http.Request
			var resp *httptest.ResponseRecorder

			var requestUrl, fileHash, fileName string

			callMethod := func() {
				if req, err = http.NewRequest("GET", requestUrl, nil); err != nil {
					panic(err)
				}
				resp = httptest.NewRecorder()
				controller.ServeFile(resp, req)
			}

			Convey(`And the URL is valid`, func() {
				requestUrl = fmt.Sprintf("%s/file/%s/%s", serverName, fileHash, fileName)

				Convey(`And the file does not exist`, func() {
					fileHash = "unknown-test-text"
					fileName = "pointless-name.mp3"

					mockFile.GetFunc = func(hash, name string) (*os.File, error) {
						return nil, new(usecases.ErrNotFound)
					}
					callMethod()

					Convey(`Then the status code will be StatusNotFound`, func() {
						So(resp.Code, ShouldEqual, http.StatusNotFound)
					})
				})

				Convey(`And the file does exist`, func() {

					tempfile, err := ioutil.TempFile("", "")
					if err != nil {
						panic(err)
					}
					defer func() {
						if err := os.Remove(tempfile.Name()); err != nil {
							panic(err)
						}
					}()
					bodyContents := "Test data"
					fmt.Fprintf(tempfile, string(bodyContents))

					mockFile.GetFunc = func(hash, name string) (*os.File, error) {
						return tempfile, nil
					}
					callMethod()

					Convey(`Then the contents of the file will be returned`, func() {
						bodyText, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							panic(err)
						}
						So(string(bodyText), ShouldEqual, bodyContents)
					})

					Convey(`Then the Content-Disposition header will be set`, func() {
						So(resp.HeaderMap["Content-Disposition"][0], ShouldEqual, fmt.Sprintf("attachment; filename=%s", path.Base(tempfile.Name())))
					})
				})
			})

			Convey(`And the URL is invalid`, func() {
				requestUrl = fmt.Sprintf("%s/file/%s", serverName, fileHash)
				callMethod()

				Convey(`Then the status code will be BadRequest`, func() {
					So(resp.Code, ShouldEqual, http.StatusNotFound)
				})
			})
		})
	})
}

type mockFileProvider struct {
	GetFunc func(string, string) (*os.File, error)
}

func (mock mockFileProvider) Get(hash, name string) (*os.File, error) {
	return mock.GetFunc(hash, name)
}
