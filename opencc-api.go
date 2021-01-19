package main

import (
	"encoding/json"
	"fmt"
	occ "github.com/gwd0715/opencc"
	"log"
	"net/http"
	"strings"
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
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hit")
	fmt.Printf("URL= %q \n", r.URL.Path)

	r.ParseForm()
	title := r.FormValue("title")
	content := r.FormValue("content")
	fmt.Println("title: ", title)
	//fmt.Println(content)

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
	//if valid {
	cc, err := occ.NewOpenCC(mode)
	if err != nil {
		fmt.Println(err)
		return
	}
	output, err := cc.ConvertText(title)
	title = output
	output, err = cc.ConvertText(content)
	content = output
	fmt.Println("Converted")
	//}

	ret := new(Ret)
	ret.Title = title
	ret.Content = content
	retJson, e := json.Marshal(ret)
	if e != nil {
		fmt.Println(e)
	}
	//fmt.Println(string(retJson))
	fmt.Fprint(w, string(retJson))
}
