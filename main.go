package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var basePath string

func Index(w http.ResponseWriter, r *http.Request) {
	html, err := Asset("assets/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(html)
}

func GetEntries(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	config := NewConfigStore(basePath)

	config.Load()

	entries := config.GetEntries()

	js, err := json.Marshal(entries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func PostEntry(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var entry ConfigEntry

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &entry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	config := NewConfigStore(basePath)

	config.Load()

	config.SetEntry(&entry)

	config.Save()

	configWriter := NewConfigWriter(basePath)

	configWriter.WriteSwitchVcl(config.GetEntries())

	w.WriteHeader(http.StatusCreated)
}

func DeleteEntry(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	config := NewConfigStore(basePath)

	config.Load()

	host := params.ByName("host")

	config.DeleteEntry(host)

	config.Save()

	configWriter := NewConfigWriter(basePath)

	configWriter.WriteSwitchVcl(config.GetEntries())

	w.WriteHeader(http.StatusAccepted)
}

func Apply(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	configWriter := NewConfigWriter(basePath)

	err := configWriter.ApplyConfiguration()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func BasicAuth(h http.Handler, user, pass []byte) SecureHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		const basicAuthPrefix string = "Basic "

		// Get the Basic Authentication credentials
		auth := r.Header.Get("Authorization")

		if strings.HasPrefix(auth, basicAuthPrefix) {
			// Check credentials
			payload, err := base64.StdEncoding.DecodeString(auth[len(basicAuthPrefix):])

			if err == nil {
				pair := bytes.SplitN(payload, []byte(":"), 2)

				if len(pair) == 2 &&
					bytes.Equal(pair[0], user) &&
					bytes.Equal(pair[1], pass) {

					// Delegate request to the given handle
					h.ServeHTTP(w, r)
					return
				}
			}
		}

		// Request Basic Authentication otherwise
		w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")

		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}

type SecureHandler func(http.ResponseWriter, *http.Request)

func (handler SecureHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler(w, r)
}

func main() {
	var err error

	if len(os.Args) < 3 {
		log.Fatal("Must specify listen address and password as command line arguments")
	}

	basePath, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	listen := os.Args[1]
	password := os.Args[2]

	api := httprouter.New()

	api.GET("/api/entries", GetEntries)
	api.POST("/api/entries", PostEntry)
	api.DELETE("/api/entries/:host", DeleteEntry)
	api.POST("/api/apply", Apply)

	mux := http.NewServeMux()

	mux.Handle("/api/", api)
	mux.HandleFunc("/", Index)

	log.Fatal(http.ListenAndServe(listen, BasicAuth(mux, []byte("admin"), []byte(password))))
}
