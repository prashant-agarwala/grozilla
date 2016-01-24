package main

import
(
  "log"
  "net/http"
  "strings"
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

func getFinalurl(url string) (string, http.Header){
  client  := &http.Client{}
  res,err := client.Head(url)
  if err != nil{
    log.Fatal(err)
  }
  responseUrl := res.Request.URL.String()
  log.Println(responseUrl)
  if responseUrl != url {
    return getFinalurl(responseUrl)
  }
  return responseUrl, res.Header
}
