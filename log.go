package main

import (
    "github.com/cheggaaa/pb"
    "strconv"
    "encoding/binary"
    "os"
    "io/ioutil"
    "log"
)

type ConnectionLog struct{
    stats []ConnectionStat
    pool *pb.Pool
    totalbar *pb.ProgressBar
}

type ConnectionStat struct{
    connectionIndex int
    pbar *pb.ProgressBar
    Err error
}

var connLog ConnectionLog

func SetupLog(length ,noOfConn int) error {
    connLog.stats = make([]ConnectionStat,noOfConn)
    barArray := make([]*pb.ProgressBar,noOfConn+1)
    len_sub := length / noOfConn
    for i := 0 ; i< noOfConn; i++ {
      file_begin := len_sub * i
      file_end := len_sub * (i + 1)
      if (i == noOfConn - 1) {
        file_end = length
      }
      bar := pb.New(file_end - file_begin).Prefix("Connection " + strconv.Itoa(i+1) + " ")
      customizeBar(bar)
      connLog.stats[i] = ConnectionStat{connectionIndex: i, pbar: bar}
      barArray[i] = bar
    }
    bar := pb.New(length).Prefix("Total ")
    customizeBar(bar)
    connLog.totalbar = bar
    barArray[noOfConn] = bar
    var err error
    connLog.pool, err = pb.StartPool(barArray...)
    if err != nil{
      return err
    }
    return nil
}

func customizeBar( bar *pb.ProgressBar){
    bar.ShowCounters = true
    bar.ShowTimeLeft = false
    bar.ShowSpeed = true
    bar.SetMaxWidth(80)
    bar.SetUnits(pb.U_BYTES)
}
func SetupResumeLog(filename string,length,noOfConn int) error {
    connLog.stats = make([]ConnectionStat,noOfConn)
    barArray := make([]*pb.ProgressBar,noOfConn+1)
    totalbar := pb.New(length).Prefix("Total ")
    len_sub := length / noOfConn
    for i := 0; i < noOfConn ; i++ {
      part_filename := "temp/" + filename + "_" + strconv.Itoa(i)
      if _, err := os.Stat(part_filename); err == nil {
          reader,err := ioutil.ReadFile(part_filename)
          if (err != nil){
            return err
          }
          header := reader[:16]
          file_begin := int(binary.LittleEndian.Uint64(header[0:8]))
          file_end := int(binary.LittleEndian.Uint64(header[8:16]))
          bar := pb.New(file_end - file_begin).Prefix("Connection " + strconv.Itoa(i+1) + " ")
          for j := 0 ; j < len(reader)-16 ; j++{
            bar.Increment()
            totalbar.Increment()
          }
          customizeBar(bar)
          connLog.stats[i] = ConnectionStat{connectionIndex: i , pbar: bar}
          barArray[i] = bar
      } else {
          file_begin := len_sub * i
          file_end := len_sub * (i + 1)
          if (i == noOfConn - 1) {
            file_end = length
          }
          bar := pb.New(file_end - file_begin).Prefix("Connection " + strconv.Itoa(i+1) + " ")
          customizeBar(bar)
          connLog.stats[i] = ConnectionStat{connectionIndex: i , pbar: bar}
          barArray[i] = bar
      }
    }
    customizeBar(totalbar)
    connLog.totalbar = totalbar
    barArray[noOfConn] = totalbar
    var err error
    connLog.pool, err = pb.StartPool(barArray...)
    if err != nil{
      return err
    }
    return nil
}


func UpdateStat(i int,file_begin int,file_end int){
   for j := file_begin ; j < file_end; j++ {
     connLog.stats[i].pbar.Increment()
     connLog.totalbar.Increment()
   }
}

func FinishLog(){
   connLog.pool.Stop()
}

func ReportErrorStat(i int,err error,noOfConn int){
   connLog.stats[i].Err = err
   connLog.pool.Stop()
   log.Println()
   log.Println("Error in connection " + strconv.Itoa(i+1) + " : "+ err.Error())
   log.Println()
   barArray := make([]*pb.ProgressBar,noOfConn)
   for i := 0 ; i < noOfConn ; i++ {
     barArray[i] = connLog.stats[i].pbar
   }
   connLog.pool, _ = pb.StartPool(barArray...)
}
