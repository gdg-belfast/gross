package main

import (
	"github.com/codegangsta/cli"
	"github.com/gobelfast/gross/mediafile"
	"github.com/gobelfast/gross/server"
	"io/ioutil"
	"log"
	"os"
)

func RunServer(c *cli.Context) {
	if len(c.Args()) == 0 {
		log.Fatalln("No directories provided")
	}

	additions := make(chan *mediafile.File)
	defer close(additions)

	for _, directory := range c.Args() {
		if _, err := ioutil.ReadDir(directory); err != nil {
			log.Fatalln(err)
		}
		go mediafile.MonitorDirectory(directory, additions)
	}
	server := server.NewServer("http://localhost", 64055)
	server.SetFileInput(additions)
	log.Fatalln(server.Run())
}

func main() {
	log.Println("Running GRoSS")
	app := cli.NewApp()
	app.Name = "gross"
	app.Usage = "run an RSS server"
	app.Action = RunServer
	app.Run(os.Args)
}
