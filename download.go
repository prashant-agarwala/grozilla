package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

//Download downloads a file from a given url by creating parallel connection
func Download(url string, length int) {
	partLength := length / *noOfFiles
	filename := getFilenameFromURL(url)
	filename = getFilename(filename)
	if _, err := os.Stat("temp/" + filename + "_0"); err == nil {
		log.Fatal("Downloading has already started, resume downloading.")
	}
	if err := SetupLog(length, *noOfFiles); err != nil {
		log.Fatal(err)
	}
	for i := 0; i < *noOfFiles; i++ {
		byteStart := partLength * (i)
		byteEnd := byteStart + partLength
		if i == *noOfFiles-1 {
			byteEnd = length
		}
		os.MkdirAll("temp/", 0777)
		createTempFile("temp/"+filename+"_"+strconv.Itoa(i), byteStart, byteEnd)
		wg.Add(1)
		go downloadPart(url, filename, i, byteStart, byteEnd)
	}
	wg.Wait()
	FinishLog()
	if !errorGoRoutine {
		mergeFiles(filename, *noOfFiles)
		clearFiles(filename, *noOfFiles)
		log.Println("download successful")
	} else {
		log.Println("download unsuccessful")
	}
}

//Resume resumes a interrupted download by creating same number of connection
func Resume(url string, length int) {
	filename := getFilenameFromURL(url)
	filename = getFilename(filename)
	*noOfFiles = noOfExistingConnection(filename, length)
	partLength := length / *noOfFiles
	if err := SetupResumeLog(filename, length, *noOfFiles); err != nil {
		log.Fatal(err)
	}
	for i := 0; i < *noOfFiles; i++ {
		partFilename := "temp/" + filename + "_" + strconv.Itoa(i)
		if _, err := os.Stat(partFilename); err != nil {
			byteStart := partLength * (i)
			byteEnd := byteStart + partLength
			if i == *noOfFiles-1 {
				byteEnd = length
			}
			wg.Add(1)
			go downloadPart(url, filename, i, byteStart, byteEnd)
		} else {
			byteStart, byteEnd := readHeader(partFilename)
			if byteStart < byteEnd {
				wg.Add(1)
				go downloadPart(url, filename, i, byteStart, byteEnd)
			}
		}
	}
	wg.Wait()
	FinishLog()
	if !errorGoRoutine {
		mergeFiles(filename, *noOfFiles)
		clearFiles(filename, *noOfFiles)
		log.Println("download successful")
	} else {
		log.Println("download unsuccessful")
	}
}

//DownloadSingle downloads a file from a given url by creating single connection
func DownloadSingle(url string) {
	filename := getFilenameFromURL(url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	reader, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(filename, reader, 0666)
	if err != nil {
		log.Fatal(err)
	}
}
