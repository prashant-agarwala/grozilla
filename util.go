package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func acceptRanges(m http.Header) bool {
	for _, v := range m["Accept-Ranges"] {
		if v == "bytes" {
			return true
		}
	}
	return false
}

func getFilenameFromURL(url string) string {
	file := url[strings.LastIndex(url, "/")+1:]
	if strings.Index(file, "?") != -1 {
		return file[:strings.Index(file, "?")]
	}
	return file
}

func getContentLength(m http.Header) int {
	length, _ := strconv.Atoi(m["Content-Length"][0])
	return length
}

func getFinalurl(url string) (string, http.Header) {
	client := &http.Client{}
	res, err := client.Head(url)
	if err != nil {
		log.Fatal(err)
	}
	responseURL := res.Request.URL.String()
	if responseURL != url {
		return getFinalurl(responseURL)
	}
	return responseURL, res.Header
}

func validateFlags() {
	if *noOfFiles <= 0 || *maxTryCount <= 0 || *timeout <= 0 {
		log.Println("Give a value greater than 0")
		flag.Usage()
		os.Exit(1)
	}
	if !(*ovrdConnLimit) {
		if *noOfFiles > 20 {
			log.Println("Connection limit restricted to 20, either use lower value or override using -N")
			flag.Usage()
			os.Exit(1)
		}
	}
}
