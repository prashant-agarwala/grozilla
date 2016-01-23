package main

import
(
  "log"
  "net/http"
)

func main(){
    log.Println("Hello world")
    url := "http://download.wavetlan.com/SVV/Media/HTTP/H264/Talkinghead_Media/H264_test1_Talkinghead_mp4_480x360.mp4"
    res, err := http.Head(url);
    if err != nil{
      log.Fatal(err)
    }
    log.Println(res.Header)
}
