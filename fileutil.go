package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func createTempFile(partFilename string, fileBegin, fileEnd int) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, int64(fileBegin)); err != nil {
		log.Fatal(err)
	}
	if err := binary.Write(buf, binary.LittleEndian, int64(fileEnd)); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(partFilename, buf.Bytes(), 0666); err != nil {
		log.Fatal(err)
	}
}

func writeBytes(partFilename string, reader []byte, byteStart, byteEnd int) error {
	if err := os.MkdirAll("temp/", 0777); err != nil {
		return err
	}
	if _, err := os.Stat("temp/" + partFilename); err != nil {
		createTempFile("temp/"+partFilename, byteStart, byteEnd)
	}
	file, err := os.OpenFile("temp/"+partFilename, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = file.WriteString(string(reader)); err != nil {
		return err
	}
	return nil
}

func readHeader(partFilename string) (int, int) {
	reader, err := ioutil.ReadFile(partFilename)
	if err != nil {
		log.Fatal(err)
	}
	header := reader[:16]
	byteStart := int(binary.LittleEndian.Uint64(header[0:8])) + len(reader) - 16
	byteEnd := int(binary.LittleEndian.Uint64(header[8:16]))
	return byteStart, byteEnd
}

func mergeFiles(filename string, count int) {
	tempFilename := strconv.Itoa(time.Now().Nanosecond()) + "_" + filename
	for i := 0; i < count; i++ {
		partFilename := "temp/" + filename + "_" + strconv.Itoa(i)
		file, err := os.OpenFile(tempFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		reader, err := ioutil.ReadFile(partFilename)
		reader = reader[16:]
		if err != nil {
			log.Fatal(err)
		}
		if _, err = file.WriteString(string(reader)); err != nil {
			log.Fatal(err)
		}
	}
	os.Rename(tempFilename, filename)
}
func isDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()
	_, err = f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func clearFiles(filename string, count int) {
	for i := 0; i < count; i++ {
		partFilename := "temp/" + filename + "_" + strconv.Itoa(i)
		os.Remove(partFilename)
	}
	empty, err := isDirEmpty("temp/")
	if err != nil {
		log.Fatal(err)
	}
	if empty {
		os.Remove("temp/")
	}
}

func noOfExistingConnection(filename string, length int) int {
	existingFilename := "temp/" + filename + "_0"
	if _, err := os.Stat(existingFilename); err != nil {
		log.Fatal("No file to resume downloading")
	}
	if _, err := os.Stat(existingFilename); err == nil {
		reader, err := ioutil.ReadFile(existingFilename)
		if err != nil {
			log.Fatal(err)
		}
		if len(reader) < 16 {
			log.Fatal("No file to resume downloading")
		}
		header := reader[:16]
		interval := int(binary.LittleEndian.Uint64(header[8:16])) - int(binary.LittleEndian.Uint64(header[0:8]))
		if interval == 0 {
			log.Fatal("No file to resume downloading")
		}
		return (length / interval)
	}
	return 0
}

func getFilename(filename string) string {
	j := 0
	for j = 0; ; j++ {
		if j == 1 {
			filename += "(1)"
		}
		if (j != 0) && (j != 1) {
			filename = strings.Replace(filename, "("+strconv.Itoa(j-1)+")", "("+strconv.Itoa(j)+")", 1)
		}
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			break
		}
	}
	if j != 0 && j != 1 {
		filename = strings.Replace(filename, "("+strconv.Itoa(j-1)+")", "("+strconv.Itoa(j)+")", 1)
	}
	return filename
}
