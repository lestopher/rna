package main

import (
	"flag"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v1"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

// ReleaseFile is used to communicate on a channel the name and body
type ReleaseFile struct {
	Name string
	Body template.HTML
}

func (rf ReleaseFile) FormattedName() string {
	if rf.Name[:2] == "./" {
		return rf.Name[2:]
	}
	return rf.Name
}

var (
	port  = flag.String("port", ":8888",
		"the port to run the release notes aggregator on")
	conf        = flag.String("conf", "./example/repos.yml", "path to the config file")
	confFileDir string
	wg          sync.WaitGroup
)

func main() {
	tmpl, err := template.New("default.tmpl.html").ParseFiles("./templates/default.tmpl.html")
	generalErrorHandler(err)

	flag.Parse()
	if len(*conf) == 0 {
		panic("Configuration file is not set.")
	}

	log.Printf("Using configuration file: %s\n", *conf)
	repos := parseConfigFile()

	confString := *conf
	confFileDir = confString[:strings.LastIndex(confString, "/")]

	output := make(chan *ReleaseFile)

	if len(repos) > 0 {
		getReleaseNotes(repos, output)
		writeToFile(output, tmpl)
	} else {
		panic("repos is nil")
	}

	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/static/", StaticHandler)
	http.ListenAndServe(*port, nil)
}

// RootHandler takes care of stuff that comes out of the compiled_html folder
func RootHandler(rw http.ResponseWriter, r *http.Request) {
	http.ServeFile(rw, r, "./compiled_html/"+r.URL.Path[1:])
}

// StaticHandler takes care of the static assets
func StaticHandler(rw http.ResponseWriter, r *http.Request) {
	http.ServeFile(rw, r, r.URL.Path[1:])
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

func getReleaseNotes(repos []string, output chan *ReleaseFile) {
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

			output <- &ReleaseFile{
				Name: repo,
				Body: template.HTML(blackfriday.MarkdownBasic(file)),
			}
			wg.Done()
		}(repo)
	}

	// Start a goroutine to wait for the waitgroup to hit zero, because otherwise
	// we'll get deadlocked goroutines
	go func() {
		wg.Wait()
		close(output)
	}()
}

func writeToFile(output <-chan *ReleaseFile, tmpl *template.Template) {
	for note := range output {
		f, err := os.Create("./compiled_html/"+note.Name+".html")
		generalErrorHandler(err)

		defer f.Close()

		err = tmpl.Execute(f, note)

		generalErrorHandler(err)
	}
}
