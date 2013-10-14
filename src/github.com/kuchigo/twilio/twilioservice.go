package main

import (
    "code.google.com/p/gorest"
    "github.com/garyburd/redigo/redis"
    "github.com/kuchigo"
    "net/http"
    "fmt"
)

func main(){
    gorest.RegisterService(new(TwilioService))
    http.Handle("/",gorest.Handle())
    http.ListenAndServe(":8787",nil)
}

//************************Define Service***************************

type TwilioService struct{

    //Service level config
    gorest.RestService  `root:"/fetch/" consumes:"application/json"  produced:"application/json"`

    userDetails gorest.EndPoint `method:"GET" path:"/items/{number:string}/{thresholdTime:float64}" output:"string"`
}

//Handler Methods: Method names must be the same as in config, but exported (starts with uppercase)

func(serv TwilioService) UserDetails(number string, thresholdTime float64)(resp string){
    c, err := redis.Dial("tcp", ":6379")
    defer c.Close()
    if err != nil {
      return fmt.Sprint("redis dial", err)
    }
    value,err := redis.Int(c.Do("GET", number))
    if err != nil {
        value = 0
    }
    newIndex, listing := craigslisting.GetListings(thresholdTime, value)
    c.Do("SET", number, newIndex)
    return listing
}

