package h5parser 

import (
   "code.google.com/p/go.net/html"
	"net/http"
	"log"
	"io/ioutil"
	"strings"
	"encoding/json"
    "regexp"
    "fmt"
)

func ParseAndPrint() map[string]string{
   //TODO : Take this url as parameter
   res, err := http.Get("http://sfbay.craigslist.org/search/apa/pen?query=&zoomToPosting=&srchType=A&minAsk=&maxAsk=2500&bedrooms=2&housing_type=&nh=77&nh=79&nh=81&nh=83&nh=84&nh=87")
   if err != nil {
     log.Fatal(err)
   }
   body, err := ioutil.ReadAll(res.Body)
   res.Body.Close()
   doc, err := html.Parse(strings.NewReader(string(body)))
   if err != nil {
     log.Fatal(err)
   }
   returnUrl := make(map[string]string)
   var checkForListings func(*html.Node)
   checkForListings = func(n *html.Node) {
     if n.Type == html.ElementNode && n.Data == "a" {
       for _,a := range n.Attr {
         if a.Key  == "href" && strings.HasPrefix(a.Val, "/pen/apa") {
          if n.FirstChild != nil {
           returnUrl["http://sfbay.craigslist.org" + a.Val] = n.FirstChild.Data
          }
         }
       }
     }
     for c := n.FirstChild; c != nil; c = c.NextSibling {
      checkForListings(c)
     }
   }
   checkForListings(doc)
   return returnUrl
}

func ParseForAddress(craigsUrl string) (address string, phoneNumber string, err error){
 res, err := http.Get(craigsUrl)
 //phoneRegex,_ := regexp.Compile(`(?:(?:(?:01\d{9}(?:[\- \,])*)|(?:2\d{7}[\- \,]*))){1,}`)
 phoneRegex,_ := regexp.Compile(`((\(\d{3}\))|(\d{3}-)) ?\d{3}-\d{4}`)
 if err != nil {
  log.Fatal(err)
 }
 body, err := ioutil.ReadAll(res.Body)
 res.Body.Close()
 doc, err := html.Parse(strings.NewReader(string(body)))
 var googleMapParser func(*html.Node)
 googleMapParser = func(n *html.Node) {
  if phoneNumber == "" && n.Type == html.ElementNode && n.Data == "section" {
   c := n.FirstChild
   if c == nil {
    return
   }
   for _,section := range n.Attr {
    if section.Key == "id" && section.Val == "postingbody"{
     for c := n.FirstChild; c != nil; c = c.NextSibling {
      phoneNumber = phoneRegex.FindString(c.Data)
      fmt.Println(c.Data, " daka ", phoneNumber)
      if phoneNumber != "" {
       break
      }
     }
    }
   }
  } else if n.Type == html.ElementNode && n.Data == "a" {
   c := n.FirstChild
   if c == nil {
    return
   }
   if strings.Contains(strings.ToLower(c.Data), "google") {
    for _,a := range n.Attr {
     if a.Key == "href" {
      address = a.Val
      break
     }
    }
   }
  }
  for c := n.FirstChild; c != nil; c = c.NextSibling {
   googleMapParser(c)
  }
 }
 googleMapParser(doc)
 return address, phoneNumber, err
}

func TransitTimeCaluclator(origin string, destination string) (duration float64){
 res, _ := http.Get("http://maps.googleapis.com/maps/api/directions/json?origin="+ origin +"&destination="+destination+"&sensor=false&mode=transit&departure_time=1379344500")
 body, _ := ioutil.ReadAll(res.Body)
 res.Body.Close()
 var f interface{}
 err := json.Unmarshal(body, &f)
 if err != nil {
  log.Fatal(err)
 }
 completeRes := f.(map[string]interface{})
 for _, v := range completeRes {
  switch vv := v.(type) {
   case []interface{}:
    for _, u := range vv {
     firstRoute := u.(map[string]interface{})
     legs := firstRoute["legs"].([]interface{})
     firstLeg := legs[0].(map[string]interface{})
     durationVal := firstLeg["duration"].(map[string]interface{})
     duration = durationVal["value"].(float64)
     return duration
    }   
  }
 }
 return duration
}
