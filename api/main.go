package main

// Golang + MySQL: https://tutorialedge.net/golang/golang-mysql-tutorial/
import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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

type Message struct {
	id         int
	message    string
	created_at time.Time
	user_id    int
}

type User struct {
	id         int
	username   string //remove
	password   string //remove
	created_at time.Time
	last_seen  time.Time
}

var dailyMessage = ""

func fck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func logHTTPReq(r *http.Request) {
	fmt.Println(fmt.Sprintf("Method:%s\nPath:%s\nUserAgent:%s\nRemote Address:%\nsSize:%s\n\n", r.Method, r.URL.String(), r.UserAgent(), r.RemoteAddr, strconv.FormatInt(r.ContentLength, 10)))
}

// ServeHTTP Response controller
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	logHTTPReq(req)
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

// handleHello Handles the /message route
func handleMessage(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("API-KEY", "YEA THIS IS THE FLAG")
	err := hello(dailyMessage).Render(req.Context(), w)
	fck(err)
}

// handleLoginForm Handles the /auth/loginForm route
func handleLoginForm(w http.ResponseWriter, req *http.Request) {
	err := loginForm().Render(req.Context(), w)
	fck(err)
}

// TEMPORARY CTF FOR TRENT
func handleDeleteAllUsers(w http.ResponseWriter, req *http.Request) {
	err := deleteAllUsers().Render(req.Context(), w)
	fck(err)
}
func handleDeleteAllUsersOptions(w http.ResponseWriter, req *http.Request) {
	err := deleteAllUsersOptions().Render(req.Context(), w)
	fck(err)
}

// handleLogin Handles the /auth/login route
func handleLogin(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	fck(err)
	username := ""
	password := ""
	for key, value := range req.Form {
		fmt.Println(key, value)
		if key == "username" {
			username = value[0]
		}
		if key == "password" {
			password = value[0]
		}
	}
	// DB
	db, err := sql.Open("mysql", "myrdsuser:myrdspassword@tcp(myrdsinstance.ct9ignv7dzxg.us-west-2.rds.amazonaws.com:3306)/dev")
	if err != nil {
		panic(err.Error())
	}
	defer func(db *sql.DB) {
		err := db.Close()
		fck(err)
	}(db)

	users, err := db.Query(fmt.Sprintf("SELECT id,username,password,created_at,last_seen FROM Users WHERE username=\"%s\" AND password=\"%s\"", username, password))
	if err != nil {
		panic(err.Error())
	}
	defer func(users *sql.Rows) {
		err := users.Close()
		fck(err)
	}(users)
	columnNames, err := users.Columns()
	fck(err)
	rc := NewMapStringScan(columnNames)
	user := ""
	pass := ""
	for users.Next() {
		err := rc.Update(users)
		fck(err)
		cv := rc.Get()
		user = cv["username"]
		pass = cv["password"]
	}

	if user != "" || pass != "" {
		err = successfulLogin().Render(req.Context(), w)
		fck(err)
	} else {
		err = unsuccessfulLogin().Render(req.Context(), w)
		fck(err)
	}
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
		err := rc.Update(messages)
		fck(err)
		cv := rc.Get()
		dailyMessage = cv["raw"]
		log.Printf("%#v\n\n", cv["created_at"])
	}

	// Server
	r := New()
	r.HandleFunc(http.MethodGet, "/message", handleMessage)
	r.HandleFunc(http.MethodGet, "/auth/loginForm", handleLoginForm)
	r.HandleFunc(http.MethodPost, "/auth/login", handleLogin)
	r.HandleFunc(http.MethodOptions, "/deleteAllUsers", handleDeleteAllUsersOptions)
	r.HandleFunc(http.MethodPost, "/deleteAllUsers", handleDeleteAllUsers)
	err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), r)
	fck(err)
}
