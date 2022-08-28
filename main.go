package main

import (
	"flag"
	"fmt"
	"github.com/sadex11/gopher02urlshorthener/urlshort"
	"net/http"
)

const (
	yamlPathDefault = "path.yaml"
	jsonPathDefault = "path.json"
	boltPathDefault = "bolt.db"
)

func main() {
	yamlFileConfig := flag.String("yaml", yamlPathDefault, "A path to a file with yaml path config")
	jsonFileConfig := flag.String("json", jsonPathDefault, "A path to a file with json path config")
	boltConfig := flag.Bool("boldb", false, "Use bolt database")

	flag.Parse()

	mux := defaultMux()

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	handlerType := "map"
	handler := mapHandler

	if *yamlFileConfig != yamlPathDefault {
		handler = urlshort.CreateYamlHandler(*yamlFileConfig, mapHandler)
		handlerType = "yaml"
	} else if *jsonFileConfig != jsonPathDefault {
		handler = urlshort.CreateJsonHandler(*jsonFileConfig, mapHandler)
		handlerType = "json"
	} else if *boltConfig {
		handler = urlshort.CreateBoltConfig(jsonPathDefault, mapHandler)
	}

	fmt.Println("Starting the server on :8080 - used handler", handlerType)
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
