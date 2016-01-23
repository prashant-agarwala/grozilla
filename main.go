package main

import
(
  "log"
  "net/http"
  "strings"
  "sync"
  "strconv"
)

func acceptRanges(m http.Header) bool {
  for _,v := range m["Accept-Ranges"]{
    if v == "bytes"{
      return true
    }
  }
  return false
}

func getFilenameFromUrl(url string) string {
  file := url[strings.LastIndex(url,"/")+1:]
  if (strings.Index(file,"?") != -1) {
    return file[:strings.Index(file,"?")]
  }
  return file
}

func getContentLength(m http.Header) int {
  length, _ := strconv.Atoi(m["Content-Length"][0])
  return length
}

var wg sync.WaitGroup

func main(){
    log.Println("Hello world")
    url := "http://download.wavetlan.com/SVV/Media/HTTP/H264/Talkinghead_Media/H264_test1_Talkinghead_mp4_480x360.mp4"
    res, err := http.Head(url);
    if err != nil{
      log.Fatal(err)
    }
    log.Println(acceptRanges(res.Header))
    //log.Println(getFilenameFromUrl(url))

    Download(url, getContentLength(res.Header))

}
