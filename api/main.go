package main

// Golang + MySQL: https://tutorialedge.net/golang/golang-mysql-tutorial/
import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strings"
)

type httpMethod string
type urlPattern string

type routeRules struct {
	methods map[httpMethod]http.Handler
}

type Router struct {
	routes map[urlPattern]routeRules
}

// ServeHTTP Response controller
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	foundRoute, exists := r.routes[urlPattern(req.URL.Path)]
	if !exists {
		http.NotFound(w, req)
		return
	}
	handler, exists := foundRoute.methods[httpMethod(req.Method)]
	if !exists {
		notAllowed(w, foundRoute)
		return
	}
	handler.ServeHTTP(w, req)
}

// HandleFunc Establishes association between any route and handler
func (r *Router) HandleFunc(method httpMethod, pattern urlPattern, f func(w http.ResponseWriter, req *http.Request)) {
	rules, exists := r.routes[pattern]
	if !exists {
		rules = routeRules{methods: make(map[httpMethod]http.Handler)}
		r.routes[pattern] = rules
	}
	rules.methods[method] = http.HandlerFunc(f)
}

// notAllowed Handles all unknown routes
func notAllowed(w http.ResponseWriter, r routeRules) {
	methods := make([]string, 1)
	for k := range r.methods {
		methods = append(methods, string(k))
	}
	w.Header().Set("Allow", strings.Join(methods, " "))
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// New Dynamically maps the routes
func New() *Router {
	return &Router{routes: make(map[urlPattern]routeRules)}
}

// handleHello Handles the /hello route
func handleHello(w http.ResponseWriter, req *http.Request) {
	err := hello().Render(req.Context(), w)
	if err != nil {
		return
	}
}

func main() {
	const PORT = 9000
	fmt.Printf("Running server at... %d\n", PORT)

	// DB
	db, err := sql.Open("mysql", "myrdsuser:myrdspassword@tcp(myrdsinstance.ct9ignv7dzxg.us-west-2.rds.amazonaws.com:3306)/myrdsinstance")
	if err != nil {
		panic(err.Error())
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	insert, err := db.Query("INSERT INTO Users VALUES ( 'Tod', 'tod@gmail.com', 'Todster1987!', now(), now() )")
	if err != nil {
		panic(err.Error())
	}
	defer func(insert *sql.Rows) {
		err := insert.Close()
		if err != nil {

		}
	}(insert)

	// Server
	r := New()
	r.HandleFunc(http.MethodGet, "/hello", handleHello)
	err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), r)
	if err != nil {
		fmt.Println(err)
	}
}
