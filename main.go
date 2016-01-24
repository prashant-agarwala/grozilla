package main

import
(
  "log"
  "flag"
)

var (
  noOfFiles     = flag.Int("n", 10, "number of parallel connection")
  resume        = flag.Bool("r", false, "resume pending download")
  maxTryCount   = flag.Int("m",1,"maximum attempts to establish a connection")
  timeout       = flag.Int("t",900,"maximum time in seconds it will wait to establish a connection")
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
    url,resHeader := getFinalurl(url)
    log.Println(acceptRanges(resHeader))
    //log.Println(getFilenameFromUrl(url))

    if *resume {
      if acceptRanges(resHeader) {
        log.Println("resume to start")
        Resume(url, getContentLength(resHeader))
      }
    } else {
        if acceptRanges(resHeader) {
          log.Println("parallel down to start")
          Download(url, getContentLength(resHeader))
        } else {
          log.Println("single down to start")
          DownloadSingle(url)
        }
    }



}
