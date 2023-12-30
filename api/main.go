package main

import (
	"fmt"
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

func (r *Router) HandleFunc(method httpMethod, pattern urlPattern, f func(w http.ResponseWriter, req *http.Request)) {
	rules, exists := r.routes[pattern]
	if !exists {
		rules = routeRules{methods: make(map[httpMethod]http.Handler)}
		r.routes[pattern] = rules
	}
	rules.methods[method] = http.HandlerFunc(f)
}

func notAllowed(w http.ResponseWriter, r routeRules) {
	methods := make([]string, 1)
	for k := range r.methods {
		methods = append(methods, string(k))
	}
	w.Header().Set("Allow", strings.Join(methods, " "))
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func New() *Router {
	return &Router{routes: make(map[urlPattern]routeRules)}
}

func handler(w http.ResponseWriter, req *http.Request) {
	err := hello().Render(req.Context(), w)
	if err != nil {
		return
	}
}

func main() {
	const PORT = 9000
	fmt.Printf("Running server at... %d\n", PORT)
	r := New()
	r.HandleFunc(http.MethodGet, "/hello", handler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), r)
	if err != nil {
		fmt.Println(err)
	}
}
