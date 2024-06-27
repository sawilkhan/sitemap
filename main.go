package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	link "github.com/sawilkhan/gophercises-html-parser"
)

/*
	USE CASES
	/some-path
	https://etcetc
	#fragment
	mailto:sawil.khan@gmail.com
*/

/*
	1. GET the webpage
	2. parse all the links on the page
	3. build proper urls with out links
	4. filter out any links w/ a diff domain
	5. Find all pages (BFS)
	6. print out XML
*/

var visited map[string]bool = make(map[string]bool)

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the url you want to build a sitemap for")
	flag.Parse()
	fmt.Println(*urlFlag)	
	
	pages := get(*urlFlag)
	i := 0
	for{
		if len(pages) == i{
			break
		}
		size := len(pages)
		for j := i; j < size; j++{
			i++
			pages = append(pages, get(pages[j])...)
		}
	}
	
	for _, page := range pages{
		fmt.Println(page)
	}

}

func get(urlString string) []string{
	resp, err := http.Get(urlString)
	if err != nil{
		panic(err)
	}
	defer resp.Body.Close()

	reqUrl := resp.Request.URL

	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host: reqUrl.Host,
	}
	base := baseUrl.String()
	return filter(hrefs(resp.Body, base), withPrefix(base))
}


func hrefs(r io.Reader, base string) []string{
	links, _ := link.Parse(r)
	var ret[]string
	for _, l := range links {
		switch{
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}
	return ret
}

func filter(links []string, keepFn func(string) bool) []string{
	var ret []string
	for _, link := range links{
		if keepFn(link){
			if exists := visited[link]; !exists{
				visited[link] = true
				ret = append(ret, link)
			}
		}
	}
	return ret
}

func withPrefix(pfx string) func(string) bool{
	return func(link string) bool{
		return strings.HasPrefix(link, pfx)
	}
}
