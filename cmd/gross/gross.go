package main

import (
	"github.com/gdg-belfast/gross/domain"
	"github.com/gdg-belfast/gross/implementation"
	"github.com/gdg-belfast/gross/infrastructure"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatalln("No directories provided")
	}

	additions := make(chan *domain.MediaFile)
	defer close(additions)

	for _, directory := range os.Args[1:] {
		if _, err := ioutil.ReadDir(directory); err != nil {
			log.Fatalln(err)
		}
		go implementation.MonitorDirectory(directory, additions)
	}
	server := infrastructure.NewRssServer("http://localhost", 64055)
	server.SetFileInput(additions)
	log.Fatalln(server.Run())
}
