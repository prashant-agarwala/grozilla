package main

import
(
  "log"
  "net/http"
  "strings"
  "strconv"
  "os"
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

func getFinalurl(url string) (string, http.Header){
  client  := &http.Client{}
  res,err := client.Head(url)
  if err != nil{
    log.Fatal(err)
  }
  responseUrl := res.Request.URL.String()
  if responseUrl != url {
    return getFinalurl(responseUrl)
  }
  return responseUrl, res.Header
}

func validateFlags(){
  if (*noOfFiles <= 0 || *maxTryCount <= 0 || *timeout <= 0){
    log.Println("Give a value greater than 0")
    flag.Usage()
    os.Exit(1)
  }
  if !(*ovrrdConnLimit) {
    if (*noOfFiles > 20){
      log.Println("Connection limit restricted to 20, either use lower value or override using -N")
      flag.Usage()
      os.Exit(1)
    }
  }
}
