package main

import (
	"flag"
	"log"
)

var (
	noOfFiles     = flag.Int("n", 10, "number of parallel connection")
	resume        = flag.Bool("r", false, "resume pending download")
	maxTryCount   = flag.Int("m", 1, "maximum attempts to establish a connection")
	timeout       = flag.Int("t", 900, "maximum time in seconds it will wait to establish a connection")
	ovrdConnLimit = flag.Bool("N", false, "maximum connection is restricted to 20, to force more connection")
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("Specify a file url to download")
	}
	validateFlags()
	url := args[0]
	url, resHeader := getFinalurl(url)
	if *resume {
		if acceptRanges(resHeader) {
			Resume(url, getContentLength(resHeader))
		}
	} else {
		if acceptRanges(resHeader) {
			Download(url, getContentLength(resHeader))
		} else {
			DownloadSingle(url)
		}
	}
}
