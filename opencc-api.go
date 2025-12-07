package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	occ "github.com/longbridgeapp/opencc"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	return rw.ResponseWriter.Write(b)
}

type Ret struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ConvertRequest struct {
	Text string `json:"text"`
}

type ConvertResponse struct {
	Converted string `json:"converted"`
}

var SCHEMES = []string{
	"s2t", "t2s", "s2tw", "tw2s", "s2hk", "hk2s", "s2twp", "tw2sp", "t2tw", "t2hk",
}

var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	log.Println("OpenCC API in Go by Colin")
	log.Printf("Version: %s", Version)
	log.Printf("Build Time: %s", BuildTime)
	log.Printf("Git Commit: %s", GitCommit)
	log.Printf("Go Version: %s", runtime.Version())
	log.Printf("Platform: %s/%s", runtime.GOOS, runtime.GOARCH)

	http.HandleFunc("/", handler)
	http.HandleFunc("/api/", apiHandler)
	log.Println("Server start")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	rw := &responseWriter{ResponseWriter: w}

	var contentLength int
	var mode string

	defer func() {
		log.Printf("[INFO] Title: %q | Content-Length: %d | Mode: %s | Duration: %v | Status: %d", "-", contentLength, mode, time.Since(start), rw.statusCode)
	}()

	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mode = strings.TrimPrefix(r.URL.Path, "/api/")
	valid := false
	for _, value := range SCHEMES {
		if value == mode {
			valid = true
			break
		}
	}

	if !valid {
		http.Error(rw, "Invalid convert scheme", http.StatusBadRequest)
		return
	}

	var req ConvertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(rw, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	contentLength = len(req.Text)

	cc, err := occ.New(mode)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Failed to initialize OpenCC: %v", err), http.StatusInternalServerError)
		return
	}

	converted, err := cc.Convert(req.Text)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Conversion failed: %v", err), http.StatusInternalServerError)
		return
	}

	resp := ConvertResponse{
		Converted: converted,
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	rw := &responseWriter{ResponseWriter: w}

	var title string
	var contentLength int
	var mode string

	defer func() {
		log.Printf("[INFO] Title: %q | Content-Length: %d | Mode: %s | Duration: %v | Status: %d", title, contentLength, mode, time.Since(start), rw.statusCode)
	}()

	r.ParseForm()
	title = r.FormValue("title")
	content := r.FormValue("content")
	contentLength = len(content)

	modeStr := string(r.URL.Path)
	parts := strings.Split(modeStr, "/")
	if len(parts) > 1 {
		mode = parts[1]
	}

	valid := false
	for _, value := range SCHEMES {
		if value == mode {
			valid = true
			break
		}
	}
	if len(mode) == 0 {
		mode = "t2s"
	} else {
		if !valid {
			// Keeping existing behavior: status 200 with error message
			fmt.Fprint(rw, "Invalid convert scheme.")
			return
		}
	}

	cc, err := occ.New(mode)
	if err != nil {
		// Keeping existing behavior: status 200 (implied) and print error to response?
		// Previous code: fmt.Println(err) -> logged to server stdout, returned 200 with empty body?
		// "fmt.Println(err); return"
		// This means client got 200 OK and empty body.
		// I will log error and return.
		log.Printf("Error initializing OpenCC: %v", err)
		return
	}
	output, err := cc.Convert(title)
	if err != nil {
		log.Printf("Error converting title: %v", err)
	}
	titleConverted := strings.TrimSpace(output)
	output, err = cc.Convert(content)
	if err != nil {
		log.Printf("Error converting content: %v", err)
	}
	content = output

	ret := new(Ret)
	ret.Title = titleConverted
	ret.Content = content
	retJson, e := json.Marshal(ret)
	if e != nil {
		log.Printf("Error marshaling response: %v", e)
	}
	fmt.Fprint(rw, string(retJson))
}
