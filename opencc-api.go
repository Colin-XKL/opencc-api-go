package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	occ "github.com/gwd0715/opencc"
)

type Ret struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

var SCHEMES = []string{
	"s2t", "t2s", "s2tw", "tw2s", "s2hk", "hk2s", "s2twp", "tw2sp", "t2tw", "t2hk",
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("OpenCC API in Go by Colin")
	log.Println("Server start")
	log.Fatal(http.ListenAndServe("localhost:3000", nil))
}
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hit")
	fmt.Printf("URL= %q \n", r.URL.Path)

	r.ParseForm()
	title := r.FormValue("title")
	content := r.FormValue("content")
	fmt.Println("title: ", title)

	modeStr := string(r.URL.Path)
	mode := strings.Split(modeStr, "/")[1]
	fmt.Println(mode)

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

	cc, err := occ.NewOpenCC(mode)
	if err != nil {
		fmt.Println(err)
		return
	}
	output, err := cc.ConvertText(title)
	title = strings.TrimSpace(output)
	output, err = cc.ConvertText(content) // 如有err返回空字符串
	content = output
	fmt.Println("Converted")

	ret := new(Ret)
	ret.Title = title
	ret.Content = content
	retJson, e := json.Marshal(ret)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Fprint(w, string(retJson))
}
