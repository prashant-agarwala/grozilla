package main

import (
  "strconv"
  "net/http"
  "log"
  "io/ioutil"
  "sync"
  "os"
)

const PACKETLENGTH = 32000
var wg sync.WaitGroup

func downloadPacket(client *http.Client, req *http.Request,part_filename string,byteStart, byteEnd int){
    resp, err := client.Do(req)
    if err != nil {
      log.Fatal(err)
    }
    reader, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      log.Fatal(err)
    }
    log.Println(part_filename, len(reader))
    err = writeBytes(part_filename,reader,byteStart,byteEnd)
    if err != nil {
      log.Fatal(err)
    }
}

func downloadPart(url,filename string, index, byteStart, byteEnd int){
    client := &http.Client{}
    part_filename := filename + "_" + strconv.Itoa(index)
    noofpacket := (byteEnd - byteStart + 1)/PACKETLENGTH + 1
    for i := 0 ; i < noofpacket; i ++ {
      packetStart := byteStart + i*PACKETLENGTH
      packetEnd   := packetStart + PACKETLENGTH
      if (i == noofpacket - 1){
        packetEnd = byteEnd
      }
      range_header := "bytes=" + strconv.Itoa(packetStart) +"-" + strconv.Itoa(packetEnd-1)
      //log.Println(range_header)
      req, _ := http.NewRequest("GET",url, nil)
      req.Header.Add("Range", range_header)
      downloadPacket(client,req,part_filename,byteStart,byteEnd)
    }
    wg.Done()
}

func Download(url string,length int){
    partLength := length / *noOfFiles
    filename := getFilenameFromUrl(url)
    for i := 0 ; i < *noOfFiles ; i++ {
      byteStart := partLength * (i)
      byteEnd   := byteStart + partLength
      if (i == *noOfFiles - 1 ){
        byteEnd = length
      }
      os.MkdirAll("temp/", 0777)
      createTempFile("temp/" + filename + "_" + strconv.Itoa(i),byteStart,byteEnd)
      wg.Add(1)
      go downloadPart(url,filename,i,byteStart,byteEnd)
    }
    wg.Wait()
    mergeFiles(filename,*noOfFiles)
    clearFiles(filename,*noOfFiles)
    reader,_ := ioutil.ReadFile(filename)
    log.Println(len(reader))
}

func Resume(url string,length int){
    filename := getFilenameFromUrl(url)
    *noOfFiles = noOfExistingConnection(filename,length)
    partLength := length / *noOfFiles
    for i := 0 ; i < *noOfFiles ; i++ {
      part_filename := "temp/" +filename + "_" + strconv.Itoa(i)
      if _, err := os.Stat(part_filename); err != nil {
        byteStart := partLength * (i)
        byteEnd   := byteStart + partLength
        if (i == *noOfFiles - 1 ){
          byteEnd = length
        }
        wg.Add(1)
        go downloadPart(url,filename,i,byteStart,byteEnd)
      } else {
        byteStart, byteEnd := readHeader(part_filename)
        if (byteStart < byteEnd) {
          wg.Add(1)
          go downloadPart(url,filename,i,byteStart,byteEnd)
        }
      }
    }
    wg.Wait()
    mergeFiles(filename,*noOfFiles)
    clearFiles(filename,*noOfFiles)
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
