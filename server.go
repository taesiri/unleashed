package main

import (
    "net/http"
    "strings"
    "log"
    "io/ioutil"
    "encoding/base64"
    "net/url"
    "net/http/cookiejar"
)


func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {

    http.HandleFunc("/", hello)

    http.HandleFunc("/enc/", func(w http.ResponseWriter, r *http.Request) {
        theUrl := strings.SplitN(r.URL.Path, "/", 3)[2]

        log.Printf("ENC")
        log.Printf(theUrl)

        sDec, _ := base64.StdEncoding.DecodeString(theUrl)
        sDec, _ = base64.StdEncoding.DecodeString(string(sDec[:]))


        address := string(sDec[:])
        log.Printf(address)


        u, err := url.Parse(address)
        check(err)
        if u.Host == "" {

        }


        options := cookiejar.Options{
               PublicSuffixList: publicsuffix.List,
        }

        cookieJar, _ := cookiejar.New(&options)

        client := &http.Client {
          Jar: cookieJar,
        }


        resp, err := client.Get(address)
        body, err := ioutil.ReadAll(resp.Body)
        check(err)



        w.Header().Set("GO!", "Unleashed!")
        w.WriteHeader(200)
        w.Write(body)
    })

    http.ListenAndServe(":8000", nil)
}


func hello(w http.ResponseWriter, r *http.Request) {
        dat, err := ioutil.ReadFile("home.html")
        check(err)
        w.Write(dat)
}
