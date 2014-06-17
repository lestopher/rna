package main

import (
	"fmt"
	"flag"
	"io/ioutil"
	// "github.com/russross/blackfriday"
	"gopkg.in/yaml.v1"
)

var (
	port = flag.String("port", ":8888",
		"the port to run the release notes aggregator on")
	conf = flag.String("conf", "./example/repos.yml", "path to the config file")
)

func main() {
	flag.Parse()
	if len(*conf) == 0 {
		panic("Configuration file is not set.")
	}

	r :=  parseConfigFile()

	fmt.Println(r)
}

func parseConfigFile() []string {
	b, err := ioutil.ReadFile(*conf)
	generalErrorHandler(err)

	m := make(map[string][]string)

	err = yaml.Unmarshal(b, &m)
	generalErrorHandler(err)

	if val, ok := m["repos"]; ok {
		return val
	}

	return nil
}

func generalErrorHandler(err error) {
	if err != nil {
		panic(err)
	}
}
