package flash

import (
	"net/http"
	"os"
	"strings"

	"github.com/segmentio/encoding/json"
)

// Response is exported.
type Response struct {
	Version string
	Items   []Item
}

// NewResponse creates and returns a new instance of Response.
func NewResponse(version string) *Response {
	return &Response{
		Version: version,
		Items:   []Item{},
	}
}

// AddItem inserts a new Item into the Response list.
func (r *Response) AddItem(key, value string) {
	r.Items = append(r.Items, Item{
		Name:  key,
		Value: value,
	})
}

// Item is exported.
type Item struct {
	Name, Value string
}

// App is exported.
type App struct {
	server  string
	version string

	cache []byte
}

// NewApp is exported.
func NewApp(server, version string) *App {
	return &App{
		server:  server,
		version: version,
	}
}

func (a *App) indexRoute(w http.ResponseWriter, r *http.Request) {
	if a.cache == nil {
		response := NewResponse(a.version)
		for _, e := range os.Environ() {
			pair := strings.SplitN(e, "=", 2)
			response.AddItem(pair[0], pair[1])
		}

		js, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		a.cache = js
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(a.cache)
}

// Register is exported.
func (a *App) register() {
	http.HandleFunc("/", Gzip(http.HandlerFunc(a.indexRoute)))
}

// Run is exported.
func (a *App) Run() {
	a.register()
	http.ListenAndServe(a.server, nil)
}
