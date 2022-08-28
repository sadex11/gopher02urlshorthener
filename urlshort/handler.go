package urlshort

import (
	"encoding/json"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type mappedHandler struct {
	urlMapper map[string]string
	fallback  http.Handler
}

type boltHandler struct {
	fallback http.Handler
}

type RedirectPath struct {
	Path string `json:"path" yaml:"path"`
	Url  string `json:"url" yaml:"url"`
}

func (h mappedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if redirectPath, ok := h.urlMapper[r.URL.Path]; ok {
		http.Redirect(w, r, redirectPath, http.StatusPermanentRedirect)
		return
	}

	h.fallback.ServeHTTP(w, r)
}

func (h boltHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if redirectPath, err := getUrlByPath(r.URL.Path); err == nil && redirectPath != "" {
		http.Redirect(w, r, redirectPath, http.StatusPermanentRedirect)
		return
	}

	h.fallback.ServeHTTP(w, r)
}

func MapHandler(pathConfig map[string]string, fallback http.Handler) http.HandlerFunc {
	return mappedHandler{pathConfig, fallback}.ServeHTTP
}

func CreateYamlHandler(filePath string, fallback http.HandlerFunc) http.HandlerFunc {
	yamlData, err := readDataFile(filePath)

	if err != nil {
		panic("Invalid yaml file")
	}

	yamlHandler, err := YAMLHandler(yamlData, fallback)

	if err != nil {
		panic(err)
	}

	return yamlHandler
}

func CreateJsonHandler(filePath string, fallback http.HandlerFunc) http.HandlerFunc {
	jsonData, err := readDataFile(filePath)

	if err != nil {
		panic("Invalid json file")
	}

	jsonHandler, err := JsonHandler(jsonData, fallback)

	if err != nil {
		panic(err)
	}

	return jsonHandler
}

func CreateBoltConfig(jsonConfigPath string, fallback http.HandlerFunc) http.HandlerFunc {
	SetupDb(jsonConfigPath)

	return boltHandler{fallback}.ServeHTTP
}

func readDataFile(filePath string) ([]byte, error) {
	f, err := os.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	return f, nil
}

func createConfig(pathData []RedirectPath) map[string]string {
	pathConfig := make(map[string]string)

	for _, v := range pathData {
		pathConfig[v.Path] = v.Url
	}

	return pathConfig
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathData []RedirectPath

	if err := yaml.Unmarshal(yml, &pathData); err != nil {
		return nil, err
	}

	return mappedHandler{createConfig(pathData), fallback}.ServeHTTP, nil
}

func JsonHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var pathData []RedirectPath

	if err := json.Unmarshal(yml, &pathData); err != nil {
		return nil, err
	}

	return mappedHandler{createConfig(pathData), fallback}.ServeHTTP, nil
}
