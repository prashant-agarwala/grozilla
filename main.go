package main

import
(
  "log"
  "net/http"
  "strings"
  "strconv"
  "flag"
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
var (
  noOfFiles     = flag.Int("n", 2, "number of parallel connection")
  resume        = flag.Bool("r", false, "resume pending download")
)
func main(){
    flag.Parse()
    log.Println("Hello world")
    args := flag.Args()
    if len(args)<1 {
      log.Fatal("Specify a file url to download")
    }
    url := args[0]
    //url := "http://download.wavetlan.com/SVV/Media/HTTP/H264/Talkinghead_Media/H264_test1_Talkinghead_mp4_480x360.mp4"
    //url := "https://nodejs.org/dist/v4.2.4/node-v4.2.4.tar.gz"
    //url := "https://storage.googleapis.com/golang/go1.5.linux-amd64.tar.gz"
    //url := "http://localhost/go1.5.linux-amd64.tar.gz"

    res, err := http.Head(url);
    if err != nil{
      log.Fatal(err)
    }
    log.Println(res.Header)
    log.Println(acceptRanges(res.Header))
    //log.Println(getFilenameFromUrl(url))

    if *resume {
      log.Println("resume to start")
      Resume(url, getContentLength(res.Header))
    } else {
      log.Println("down to start")
      Download(url, getContentLength(res.Header))
    }
    //DownloadSingle(url)


}
