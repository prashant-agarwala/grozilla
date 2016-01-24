package main

import (
    "github.com/cheggaaa/pb"
    "strconv"
)

type ConnectionLog struct{
  stats []ConnectionStat
  pool *pb.Pool
}

type ConnectionStat struct{
  connectionIndex int
  pbar *pb.ProgressBar
  Err error
}

var connLog ConnectionLog

func SetupLog(length ,noOfConn int) error {
  connLog.stats = make([]ConnectionStat,noOfConn)
  barArray := make([]*pb.ProgressBar,noOfConn)
  len_sub := length / noOfConn
  for i := 0 ; i< noOfConn; i++ {
    file_begin := len_sub * i
    file_end := len_sub * (i + 1)
    if (i == noOfConn - 1) {
      file_end = length
    }
    bar := pb.New(file_end - file_begin).Prefix("Connection " + strconv.Itoa(i+1) + " ")
    connLog.stats[i] = ConnectionStat{connectionIndex: i, pbar: bar}
   barArray[i] = bar
  }
  var err error
  connLog.pool, err = pb.StartPool(barArray...)
  if err != nil{
    return err
  }
  return nil
}

func UpdateStat(i int,file_begin int,file_end int){
 for j := file_begin ; j < file_end; j++{
  connLog.stats[i].pbar.Increment()
  }
}

func FinishLog(){
  connLog.pool.Stop()
}
