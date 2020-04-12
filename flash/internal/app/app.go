package app

import (
	"html/template"
	"log"
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

	db    ResourcesRepository
	cache []byte
}

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

// NewApp is exported.
func NewApp(server, version string, repo ResourcesRepository) *App {
	return &App{
		server:  server,
		version: version,

		db: repo,
	}
}

func (a *App) apiIndexRoute(w http.ResponseWriter, r *http.Request) {
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

func (a *App) indexRoute(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("resources/index.html")
	if err != nil {
		log.Println("failed to parse template", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resources, err := a.db.LoadResources()
	if err != nil {
		log.Println("failed to load data from database", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	articles := []Article{}
	for _, v := range resources {
		article := Article{Title: v.URL, URL: v.URL, Created: v.Created}
		if v.Source != nil {
			article.Source = *v.Source
		}

		articles = append(articles, article)
	}

	data := IndexPage{
		Title:    "Recent Articles",
		Articles: articles,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("failed to write response", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Register is exported.
func (a *App) register() {
	// website routes
	http.HandleFunc("/", a.indexRoute)

	// api routes
	http.HandleFunc("/api", Gzip(http.HandlerFunc(a.apiIndexRoute)))
}

// Run is exported.
func (a *App) Run() {
	a.register()
	err := http.ListenAndServe(a.server, nil)
	if err != nil {
		log.Println("failed to start server", err)
	}
}
