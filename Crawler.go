package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Please enter a start web address...")
		os.Exit(1)
	}
	retrivedBody := retrieveURL(args[0])
	fmt.Println("Retrieved body is ", retrivedBody)
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
