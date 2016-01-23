package main

import (
  "strconv"
  "net/http"
  "log"
  "io/ioutil"
  "sync"
)

const noOfFiles = 20
var wg sync.WaitGroup

func downloadPart(url,filename string, index, byteStart, byteEnd int){
    client := &http.Client{}
    range_header := "bytes=" + strconv.Itoa(byteStart) +"-" + strconv.Itoa(byteEnd-1)
    log.Println(range_header)
    req, _ := http.NewRequest("GET",url, nil)
    req.Header.Add("Range", range_header)
    resp, err := client.Do(req)
    if err != nil {
      log.Fatal(err)
    }
    reader, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      log.Fatal(err)
    }
    log.Println(index, len(reader))
    part_filename := filename + "_" + strconv.Itoa(index)
    err = writeBytes(part_filename , reader)
    if err != nil {
      log.Fatal(err)
    }
    wg.Done()
}

func Download(url string,length int){
    partLength := length/noOfFiles
    filename := getFilenameFromUrl(url)
    for i:= 0 ; i < noOfFiles ; i++ {
      byteStart := partLength * (i)
      byteEnd   := byteStart + partLength
      if (i == noOfFiles - 1 ){
        byteEnd = length
      }
      wg.Add(1)
      go downloadPart(url,filename,i,byteStart,byteEnd)
    }
    wg.Wait()
    mergeFiles(filename)
    clearFiles(filename)
    reader,_ := ioutil.ReadFile(filename)
    log.Println(len(reader))
}

func DownloadSingle(url string){
    filename := getFilenameFromUrl(url)
    client := &http.Client{}
    req, _ := http.NewRequest("GET",url, nil)
    resp, err := client.Do(req)
    if err != nil {
      log.Fatal(err)
    }
    reader, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      log.Fatal(err)
    }
    log.Println(len(reader))
    err = ioutil.WriteFile(filename, reader,0666)
    if err != nil {
      log.Fatal(err)
    }
}
