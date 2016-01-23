package main

import (
  "strconv"
  "log"
  "io/ioutil"
  "os"
  "encoding/binary"
  "bytes"
)

func createTempFile(part_filename string,fileBegin ,fileEnd int){
    buf := new(bytes.Buffer)
    if err := binary.Write(buf, binary.LittleEndian, int64(fileBegin)); err != nil {
      log.Fatal(err)
    }
    if err := binary.Write(buf, binary.LittleEndian, int64(fileEnd)); err != nil {
      log.Fatal(err)
    }
    if err := ioutil.WriteFile(part_filename, []byte(buf.Bytes()), 0666); err != nil {
      log.Fatal(err)
    }
}


func writeBytes(part_filename string, reader []byte, byteStart , byteEnd int) error{
    err := os.MkdirAll("temp/", 0777)
    if err != nil {
      return err
    }
    if _, err := os.Stat("temp/" + part_filename); err != nil {
      log.Println("new file to be created  " , part_filename)
      createTempFile("temp/" + part_filename,byteStart,byteEnd)
    }
    file, err := os.OpenFile("temp/" + part_filename, os.O_WRONLY|os.O_APPEND,0666)
    if err != nil {
      return err
    }
    defer file.Close()
    if _, err = file.WriteString(string(reader)); err != nil {
      return err
    }
    return nil
}

func readHeader(part_filename string) (int,int){
    reader,err := ioutil.ReadFile(part_filename)
    if (err != nil) {
      log.Fatal(err)
    }
    header := reader[:16]
    byteStart := int(binary.LittleEndian.Uint64(header[0:8])) + len(reader)-16
    byteEnd   := int(binary.LittleEndian.Uint64(header[8:16]))
    return byteStart,byteEnd
}

func mergeFiles(filename string){
    for i := 0; i < noOfFiles ; i++ {
        part_filename := "temp/" + filename + "_" + strconv.Itoa(i)
        file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
        if err != nil {
          log.Fatal(err)
        }
        defer file.Close()
        reader,err := ioutil.ReadFile(part_filename)
        reader = reader[16:]
        if err != nil {
          log.Fatal(err)
        }
        if _, err = file.WriteString(string(reader)); err != nil {
          log.Fatal(err)
        }
    }
}

func clearFiles(filename string){
  os.RemoveAll("temp")
  // for i := 0; i < noOfFiles ; i++ {
  //   part_filename := "temp/" + filename + "_" + strconv.Itoa(i)
  //   os.Remove(part_filename)
  // }
}
