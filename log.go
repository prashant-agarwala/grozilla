package main

import (
	"encoding/binary"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/cheggaaa/pb"
)

//ConnectionLog keeps log of all connection throgh progressbar
type ConnectionLog struct {
	stats    []ConnectionStat
	pool     *pb.Pool
	totalbar *pb.ProgressBar
}

//ConnectionStat keeps statistic of each connection
type ConnectionStat struct {
	connectionIndex int
	pbar            *pb.ProgressBar
	Err             error
}

var connLog ConnectionLog

//SetupLog sets up initial ConnectionLog
func SetupLog(length, noOfConn int) error {
	connLog.stats = make([]ConnectionStat, noOfConn)
	barArray := make([]*pb.ProgressBar, noOfConn+1)
	lenSub := length / noOfConn
	for i := 0; i < noOfConn; i++ {
		fileBegin := lenSub * i
		fileEnd := lenSub * (i + 1)
		if i == noOfConn-1 {
			fileEnd = length
		}
		bar := pb.New(fileEnd - fileBegin).Prefix("Connection " + strconv.Itoa(i+1) + " ")
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
	if err != nil {
		return err
	}
	return nil
}

func customizeBar(bar *pb.ProgressBar) {
	bar.ShowCounters = true
	bar.ShowTimeLeft = false
	bar.ShowSpeed = true
	bar.SetMaxWidth(80)
	bar.SetUnits(pb.U_BYTES)
}

//SetupResumeLog sets up ConnectionLog for a resumed download
func SetupResumeLog(filename string, length, noOfConn int) error {
	connLog.stats = make([]ConnectionStat, noOfConn)
	barArray := make([]*pb.ProgressBar, noOfConn+1)
	totalbar := pb.New(length).Prefix("Total ")
	lenSub := length / noOfConn
	for i := 0; i < noOfConn; i++ {
		partFilename := "temp/" + filename + "_" + strconv.Itoa(i)
		if _, err := os.Stat(partFilename); err == nil {
			reader, err := ioutil.ReadFile(partFilename)
			if err != nil {
				return err
			}
			header := reader[:16]
			fileBegin := int(binary.LittleEndian.Uint64(header[0:8]))
			fileEnd := int(binary.LittleEndian.Uint64(header[8:16]))
			bar := pb.New(fileEnd - fileBegin).Prefix("Connection " + strconv.Itoa(i+1) + " ")
			for j := 0; j < len(reader)-16; j++ {
				bar.Increment()
				totalbar.Increment()
			}
			customizeBar(bar)
			connLog.stats[i] = ConnectionStat{connectionIndex: i, pbar: bar}
			barArray[i] = bar
		} else {
			fileBegin := lenSub * i
			fileEnd := lenSub * (i + 1)
			if i == noOfConn-1 {
				fileEnd = length
			}
			bar := pb.New(fileEnd - fileBegin).Prefix("Connection " + strconv.Itoa(i+1) + " ")
			customizeBar(bar)
			connLog.stats[i] = ConnectionStat{connectionIndex: i, pbar: bar}
			barArray[i] = bar
		}
	}
	customizeBar(totalbar)
	connLog.totalbar = totalbar
	barArray[noOfConn] = totalbar
	var err error
	connLog.pool, err = pb.StartPool(barArray...)
	if err != nil {
		return err
	}
	return nil
}

//UpdateStat updates statistic of a connection
func UpdateStat(i int, fileBegin int, fileEnd int) {
	for j := fileBegin; j < fileEnd; j++ {
		connLog.stats[i].pbar.Increment()
		connLog.totalbar.Increment()
	}
}

//FinishLog stops ConnectionLog pool
func FinishLog() {
	connLog.pool.Stop()
}

//ReportErrorStat reports a log if an error occurs in a connection
func ReportErrorStat(i int, err error, noOfConn int) {
	connLog.stats[i].Err = err
	connLog.pool.Stop()
	log.Println()
	log.Println("Error in connection " + strconv.Itoa(i+1) + " : " + err.Error())
	log.Println()
	barArray := make([]*pb.ProgressBar, noOfConn+1)
	for i := 0; i < noOfConn; i++ {
		barArray[i] = connLog.stats[i].pbar
	}
	barArray[noOfConn] = connLog.totalbar
	connLog.pool, _ = pb.StartPool(barArray...)
}
