package main

import (
	"flag"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

var (
	port = flag.String("port", ":8888",
		"the port to run the release notes aggregator on")
	conf        = flag.String("conf", "./example/repos.yml", "path to the config file")
	confFileDir string
	wg          sync.WaitGroup
)

func main() {
	flag.Parse()
	if len(*conf) == 0 {
		panic("Configuration file is not set.")
	}

	log.Printf("Using configuration file: %s\n", *conf)
	repos := parseConfigFile()

	confString := *conf
	confFileDir = confString[:strings.LastIndex(confString, "/")]

	output := make(chan []byte)

	if len(repos) > 0 {
		getReleaseNotes(repos, output)
		writeToFile(output)
	} else {
		panic("repos is nil")
	}
}

// parseConfigFile parses the configugration yaml file to which the global
// variable points
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

func getReleaseNotes(repos []string, output chan []byte) {
	for _, repo := range repos {
		wg.Add(1)
		go func(repo string) {
			var path string

			if repo[0:1] == "/" { // using absolute path
				path = repo
			} else {
				if repo[:1] == "./" { // relative path dot slash
					path = confFileDir + "/" + repo[2:]
				} else { // relative path naked
					path = confFileDir + "/" + repo
				}
			}

			file, err := ioutil.ReadFile(path + "/RELEASE.md")
			generalErrorHandler(err)

			output <- blackfriday.MarkdownBasic(file)
			wg.Done()
		}(repo)
	}

	go func() {
		wg.Wait()
		close(output)
	}()
}

func writeToFile(output <-chan []byte) {
	for notes := range output {
		log.Printf("%s\n", notes)
	}
}
