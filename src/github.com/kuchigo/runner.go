package main

import (
  "github.com/kuchigo/h5parser"
  "strings"
  "fmt"
  "time"
  )

func main(){
 returnedUrl := h5parser.ParseAndPrint()
 for _,urlData := range returnedUrl {
  //get the google encoded address
  address, err := h5parser.ParseForAddress(urlData)
  if err == nil && strings.TrimSpace(address) != "" {
   //split address to get google recognizable address
   locs := strings.Split(address, "loc%3A+")
   //google registration required to make lot of queries
   time.Sleep(5 * time.Second)
   if len(locs) == 2 {
    //TODO : Add the start time as parameter, currently it is encoded to some future date
    duration := h5parser.TransitTimeCaluclator(locs[1],"625+Harrison+Street+SanFrancisco+CA")
    //if it is less than 1 hour 10 minutes get it
    if duration < 4200 {
     fmt.Println(urlData, " takes time ", duration)
    }
   }
  }
 }
}
