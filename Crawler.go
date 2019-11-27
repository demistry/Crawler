package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/jackdanger/collectlinks"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Please enter a start web address...")
		os.Exit(1)
	}

	//creating a channel to hold links as they come in
	queue := make(chan string)

	go func() {
		queue <- args[0]
	}()

	for uri := range queue {
		enqueue(uri, queue)
	}

}

func enqueue(uri string, queue chan string) {
	fmt.Println("Fetching uri...")

	tlsConfig := &tls.Config{InsecureSkipVerify: true} //skips ssl configured sites
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := http.Client{Transport: transport}

	resp, err := client.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	links := collectlinks.All(resp.Body)

	for _, link := range links {
		absolute := fixURL(link, uri)
		if uri != "" {
			go func() { queue <- absolute }()
		}

	}
}

func fixURL(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}

	baseURL, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseURL.ResolveReference(uri)
	return uri.String()
}

func retrieveURL(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
