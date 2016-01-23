package main

import (
  "strconv"
  "net/http"
  "log"
  "io/ioutil"
  "os"
)

const noOfFiles = 40
func writeBytes(part_filename string, reader []byte) error{
    file, err := os.OpenFile(part_filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE,0666)
    if err != nil {
      return err
    }
    defer file.Close()
    if _, err = file.WriteString(string(reader)); err != nil {
      return err
    }
    return nil
}

func mergeFiles(filename string){
    for i := 0; i < noOfFiles ; i++ {
        part_filename := filename + "_" + strconv.Itoa(i)
        file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
        if err != nil {
          log.Fatal(err)
        }
        defer file.Close()
        reader,err := ioutil.ReadFile(part_filename)
        if err != nil {
          log.Fatal(err)
        }
        if _, err = file.WriteString(string(reader)); err != nil {
          log.Fatal(err)
        }
    }
}

func clearFiles(filename string){
  for i := 0; i < noOfFiles ; i++ {
    part_filename := filename + "_" + strconv.Itoa(i)
    os.Remove(part_filename)
  }
}

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
    log.Println(len(reader))
    part_filename := filename + "_" + strconv.Itoa(index)
    err = writeBytes(part_filename , reader)
    if err != nil {
      log.Fatal(err)
    }
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
      downloadPart(url,filename,i,byteStart,byteEnd)
    }
    mergeFiles(filename)
    clearFiles(filename)
    reader,_ := ioutil.ReadFile(filename)
    log.Println(len(reader))

}
