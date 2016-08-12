package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/jenrik/golem"
	"io/ioutil"
	"os"
)

type Links []string

func (i *Links) String() string {
	return ""
}

func (i *Links) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	// Setup logging
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.ErrorLevel)

	// Commandline args
	queueAddr := flag.String("amqp", "amqp://localhost:5672/", "Address to connect queue")
	storageAddr := flag.String("storage", "postgresql://golem:golem@localhost:5432/golem", "Address to connect to postgresql")
	var links Links
	flag.Var(&links, "link", "Link to begin scrapping from")
	defsFile := flag.String("defs", "./definitions.json", "JSON file containing definitions")
	flag.Parse()

	err := golem.Connect(*queueAddr, *storageAddr)
	if err != nil {
		panic(err)
	}
	defer golem.Disconnect()

	data, err := ioutil.ReadFile(*defsFile)
	if err != nil {
		panic(err)
	}

	defs, err := golem.ParseDefs(&data)
	if err != nil {
		panic(err)
	}

	jobId, err := golem.SubmitJob(defs, (*[]string)(&links))
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Job succesfully submitted with id: %v\n", jobId)
	}
}
