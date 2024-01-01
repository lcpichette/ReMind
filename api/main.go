package main

// Golang + MySQL: https://tutorialedge.net/golang/golang-mysql-tutorial/
import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strings"
)

type httpMethod string
type urlPattern string

type MapStringScan struct {
	// cp are the column pointers
	cp []interface{}
	// row contains the final result
	row      map[string]string
	colCount int
	colNames []string
}

type routeRules struct {
	methods map[httpMethod]http.Handler
}

type Router struct {
	routes map[urlPattern]routeRules
}

var dailyMessage = "Hello, World"

func fck(err error) {
	if err != nil {
		log.Fatal(err)
	}
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
	err := hello(dailyMessage).Render(req.Context(), w)
	fck(err)
}

func NewMapStringScan(columnNames []string) *MapStringScan {
	lenCN := len(columnNames)
	s := &MapStringScan{
		cp:       make([]interface{}, lenCN),
		row:      make(map[string]string, lenCN),
		colCount: lenCN,
		colNames: columnNames,
	}
	for i := 0; i < lenCN; i++ {
		s.cp[i] = new(sql.RawBytes)
	}
	return s
}
func (s *MapStringScan) Update(rows *sql.Rows) error {
	if err := rows.Scan(s.cp...); err != nil {
		return err
	}

	for i := 0; i < s.colCount; i++ {
		if rb, ok := s.cp[i].(*sql.RawBytes); ok {
			s.row[s.colNames[i]] = string(*rb)
			*rb = nil // reset pointer to discard current value to avoid a bug
		} else {
			return fmt.Errorf("Cannot convert index %d column %s to type *sql.RawBytes", i, s.colNames[i])
		}
	}
	return nil
}

func (s *MapStringScan) Get() map[string]string {
	return s.row
}

func main() {
	const PORT = 9000
	fmt.Printf("Running server at... %d\n", PORT)

	// DB
	db, err := sql.Open("mysql", "myrdsuser:myrdspassword@tcp(myrdsinstance.ct9ignv7dzxg.us-west-2.rds.amazonaws.com:3306)/dev")
	if err != nil {
		panic(err.Error())
	}
	defer func(db *sql.DB) {
		err := db.Close()
		fck(err)
	}(db)

	insert, err := db.Query("INSERT INTO Analytics (event,content_grouping,created_at) VALUES ( 'page_load', 'landing_page', now() )")
	if err != nil {
		panic(err.Error())
	}
	defer func(insert *sql.Rows) {
		err := insert.Close()
		fck(err)
	}(insert)

	messages, err := db.Query("SELECT raw,created_at FROM Messages WHERE user_id = 1")
	if err != nil {
		panic(err.Error())
	}
	defer func(messages *sql.Rows) {
		err := messages.Close()
		fck(err)
	}(messages)
	columnNames, err := messages.Columns()
	fck(err)
	rc := NewMapStringScan(columnNames)
	for messages.Next() {
		//        cv, err := rowMapString(columnNames, rows)
		//        fck(err)
		err := rc.Update(messages)
		fck(err)
		cv := rc.Get()
		dailyMessage = cv["raw"]
		log.Printf("%#v\n\n", cv["created_at"])
	}

	// Server
	r := New()
	r.HandleFunc(http.MethodGet, "/hello", handleHello)
	err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), r)
	fck(err)
}
