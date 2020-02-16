package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	time "time"

	"github.com/Masterminds/sprig"
)

// Log requests
// https://www.socketloop.com/tutorials/golang-how-to-log-each-http-request-to-your-web-server
func RequestLogger(targetMux http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		targetMux.ServeHTTP(res, req)
		log.Printf(
			"%s\t%s\t%s\t%v",
			// "%s\t\t%s\t\t%s\t\t%v",
			req.Method,
			// res.StatusCode,
			req.RequestURI,
			req.RemoteAddr, // IP
			time.Since(start),
		)
	})
}

type Swatch struct {
	Adj          string `json:"adj"`
	Noun         string `json:"noun"`
	R            int    `json:"r"`
	G            int    `json:"g"`
	B            int    `json:"b"`
	CreatedAtStr string `json:"created_at"`
	CreatedAt    time.Time
}

type Template struct {
	GoogleAnalytics string
	Swatches        []Swatch
}

func main() {

	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/", indexRoute)
	mux.HandleFunc("/style.css", cssRoute)

	// Serve views directory
	// fs := http.FileServer(http.Dir("views/"))
	// http.Handle("/views/", http.StripPrefix("/views/", fs))

	// Format the host from env variables or defaults
	hostname, exists := os.LookupEnv("HOSTNAME")
	if !exists {
		hostname = "0.0.0.0"
	}
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "80"
	}
	host := fmt.Sprintf("%s:%s", hostname, port)

	log.Printf("Server running at http://%s", host)
	log.Println("Close with CTRL + C")

	// Start Server
	http.ListenAndServe(host, RequestLogger(mux))
}

func indexRoute(res http.ResponseWriter, req *http.Request) {
	// The "/" matches anything not handled elsewhere. If it's not the root
	// then report not found.
	// https://stackoverflow.com/a/50282841/247218
	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}

	// Load the template
	tmpl := template.Must(template.New("").Funcs(sprig.HtmlFuncMap()).ParseFiles("./views/swatches.gohtml"))

	// Define the template data
	data := Template{
		GoogleAnalytics: "GA-00000000",
	}

	// Load and parse the JSON into data
	file, _ := ioutil.ReadFile("./data/swatches.json")
	json.Unmarshal([]byte(file), &data.Swatches)

	// Parse dates to time.Time
	for i, swatch := range data.Swatches {
		date, _ := time.Parse(time.RFC3339Nano, swatch.CreatedAtStr)
		data.Swatches[i].CreatedAt = date
	}

	// Execute the template as the response
	tmpl.ExecuteTemplate(res, "swatches.gohtml", data)
}

func cssRoute(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./views/style.css")
}
