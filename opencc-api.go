package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"

	occ "github.com/longbridgeapp/opencc"
)

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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mode := strings.TrimPrefix(r.URL.Path, "/api/")
	valid := false
	for _, value := range SCHEMES {
		if value == mode {
			valid = true
			break
		}
	}

	if !valid {
		http.Error(w, "Invalid convert scheme", http.StatusBadRequest)
		return
	}

	var req ConvertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	cc, err := occ.New(mode)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to initialize OpenCC: %v", err), http.StatusInternalServerError)
		return
	}

	converted, err := cc.Convert(req.Text)
	if err != nil {
		http.Error(w, fmt.Sprintf("Conversion failed: %v", err), http.StatusInternalServerError)
		return
	}

	resp := ConvertResponse{
		Converted: converted,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("URL= %q \n", r.URL.Path)

	r.ParseForm()
	title := r.FormValue("title")
	content := r.FormValue("content")
	// fmt.Println("title: ", title)

	modeStr := string(r.URL.Path)
	mode := strings.Split(modeStr, "/")[1]
	// fmt.Println(mode)

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
			fmt.Fprint(w, "Invalid convert scheme.")
			return
		}
	}

	cc, err := occ.New(mode)
	if err != nil {
		fmt.Println(err)
		return
	}
	output, err := cc.Convert(title)
	if err != nil {
		fmt.Println(err)
	}
	title = strings.TrimSpace(output)
	output, err = cc.Convert(content) // 如有err返回空字符串
	if err != nil {
		fmt.Println(err)
	}
	content = output
	fmt.Println("Converted")
	fmt.Println("Title: ", title)
	fmt.Println("Content Length: ", len(content))

	ret := new(Ret)
	ret.Title = title
	ret.Content = content
	retJson, e := json.Marshal(ret)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Fprint(w, string(retJson))
}
