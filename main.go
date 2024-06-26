package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	link "github.com/sawilkhan/gophercises-html-parser"
)

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the url you want to build a sitemap for")
	flag.Parse()


	fmt.Println(*urlFlag)

	resp, err := http.Get(*urlFlag)
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
	
	links, _ := link.Parse(resp.Body)
	var hrefs[]string
	for _, l := range links {
		switch{
		case strings.HasPrefix(l.Href, "/"):
			hrefs = append(hrefs, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			hrefs = append(hrefs, l.Href)
		}
		
	}

	for _, href := range hrefs{
		fmt.Println(href)
	}

}
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