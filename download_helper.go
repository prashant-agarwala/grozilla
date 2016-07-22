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

const PACKETLENGTH = 32000

var wg sync.WaitGroup
var errorGoRoutine bool

func downloadPacket(client *http.Client, req *http.Request, part_filename string, byteStart, byteEnd int) error {
	c := make(chan httpResponse, 1)
	go func() {
		resp, err := client.Do(req)
		http_response := httpResponse{resp, err}
		c <- http_response
	}()
	select {
	case http_response := <-c:
		if err := handleResponse(http_response, part_filename, byteStart, byteEnd); err != nil {
			return err
		}
	case <-time.After(time.Second * time.Duration(*timeout)):
		err := errors.New("Manual time out as response not recieved")
		return err
	}
	return nil
}

func handleResponse(http_response httpResponse, part_filename string, byteStart, byteEnd int) error {
	if http_response.err != nil {
		return http_response.err
	}
	defer http_response.resp.Body.Close()
	reader, err := ioutil.ReadAll(http_response.resp.Body)
	if err != nil {
		return err
	}
	err = writeBytes(part_filename, reader, byteStart, byteEnd)
	if err != nil {
		return err
	}
	return nil
}

func downloadPacketWithRetry(client *http.Client, req *http.Request, part_filename string, byteStart, byteEnd int) error {
	var err error
	for i := 0; i < *maxTryCount; i++ {
		err = downloadPacket(client, req, part_filename, byteStart, byteEnd)
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
	part_filename := filename + "_" + strconv.Itoa(index)
	noofpacket := (byteEnd-byteStart+1)/PACKETLENGTH + 1
	for i := 0; i < noofpacket; i++ {
		packetStart := byteStart + i*PACKETLENGTH
		packetEnd := packetStart + PACKETLENGTH
		if i == noofpacket-1 {
			packetEnd = byteEnd
		}
		range_header := "bytes=" + strconv.Itoa(packetStart) + "-" + strconv.Itoa(packetEnd-1)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Range", range_header)
		err := downloadPacketWithRetry(client, req, part_filename, byteStart, byteEnd)
		if err != nil {
			handleErrorInGoRoutine(i, err)
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
