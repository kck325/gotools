package craigslisting

import (
  "github.com/kuchigo/h5parser"
  "strings"
  "time"
  "fmt"
  )

func GetListings(thresholdTime float64, value int)(int, string){
 returnedUrl := h5parser.ParseAndPrint()
 currentIndex := -1
 for urlData,_ := range returnedUrl {
  currentIndex++
  if currentIndex < value {
   continue
  }
  //get the google encoded address
  address, phoneNumber, err := h5parser.ParseForAddress(urlData)
  if err == nil && strings.TrimSpace(address) != "" {
   fmt.Println(phoneNumber)
   //split address to get google recognizable address
   locs := strings.Split(address, "loc%3A+")
   //google registration required to make lot of queries
   time.Sleep(5 * time.Second)
   if len(locs) == 2 {
    //TODO : Add the start time as parameter, currently it is encoded to some future date
    duration := h5parser.TransitTimeCaluclator(locs[1],"625+Harrison+Street+SanFrancisco+CA")
    if duration < thresholdTime {
     return currentIndex,fmt.Sprint(urlData, " takes time ", duration, " ", phoneNumber)
    }
   }
  } else {
    fmt.Println(err)
  }
 }
 return -1,""
}
