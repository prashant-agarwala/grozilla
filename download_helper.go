package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type httpResponse struct {
	resp *http.Response
	err  error
}

//PACKETLENGTH is size of each packet in bytes
const PACKETLENGTH = 32000

var wg sync.WaitGroup
var errorGoRoutine bool

func downloadPacket(client *http.Client, req *http.Request, partFilename string, byteStart, byteEnd int) error {
	c := make(chan httpResponse, 1)
	go func() {
		resp, err := client.Do(req)
		httpResponse := httpResponse{resp, err}
		c <- httpResponse
	}()
	select {
	case httpResponse := <-c:
		if err := handleResponse(httpResponse, partFilename, byteStart, byteEnd); err != nil {
			return err
		}
	case <-time.After(time.Second * time.Duration(*timeout)):
		err := errors.New("Manual time out as response not recieved")
		return err
	}
	return nil
}

func handleResponse(httpResponse httpResponse, partFilename string, byteStart, byteEnd int) error {
	if httpResponse.err != nil {
		return httpResponse.err
	}
	defer httpResponse.resp.Body.Close()
	reader, err := ioutil.ReadAll(httpResponse.resp.Body)
	if err != nil {
		return err
	}
	err = writeBytes(partFilename, reader, byteStart, byteEnd)
	if err != nil {
		return err
	}
	return nil
}

func downloadPacketWithRetry(client *http.Client, req *http.Request, partFilename string, byteStart, byteEnd int) error {
	var err error
	for i := 0; i < *maxTryCount; i++ {
		err = downloadPacket(client, req, partFilename, byteStart, byteEnd)
		if err == nil {
			return nil
		} else if err.Error() == "Manual time out as response not recieved" {
			continue
		} else {
			return err
		}
	}
	return err
}

func downloadPart(url, filename string, index, byteStart, byteEnd int) {
	client := &http.Client{}
	partFilename := filename + "_" + strconv.Itoa(index)
	noofpacket := (byteEnd-byteStart+1)/PACKETLENGTH + 1
	for i := 0; i < noofpacket; i++ {
		packetStart := byteStart + i*PACKETLENGTH
		packetEnd := packetStart + PACKETLENGTH
		if i == noofpacket-1 {
			packetEnd = byteEnd
		}
		rangeHeader := "bytes=" + strconv.Itoa(packetStart) + "-" + strconv.Itoa(packetEnd-1)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Range", rangeHeader)
		err := downloadPacketWithRetry(client, req, partFilename, byteStart, byteEnd)
		if err != nil {
			handleErrorInGoRoutine(index, err)
			return
		}
		UpdateStat(index, packetStart, packetEnd)
	}
	wg.Done()
}

func handleErrorInGoRoutine(index int, err error) {
	ReportErrorStat(index, err, *noOfFiles)
	errorGoRoutine = true
	wg.Done()
}
