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

const serverIP = "http://127.0.0.1:8000/"
const serverPort = ":8000"


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

        siteUrl := "";

        u, err := url.Parse(address)
        check(err)

        if u.Host == "" {
          w.Header().Set("GO!", "Not found!!")
          w.WriteHeader(404)
          w.Write([]byte("Not Found!"))
        } else {

          cookieJar, _ := cookiejar.New(nil)

          client := &http.Client {
            Jar: cookieJar,
          }

          resp, err := client.Get(address)
          body, err := ioutil.ReadAll(resp.Body)
          check(err)

          siteUrl = u.Scheme + "://" + u.Host

          s := string(body[:])

          s = strings.Replace(s, "<img src=\"/", "<img src=\"" + serverIP + siteUrl + "/" , -1)
          s = strings.Replace(s, "url('/", "url('" + serverIP + siteUrl + "/" , -1)


          s = strings.Replace(s, "href=\"/", "href=\"" + serverIP + siteUrl + "/" , -1)

          s = strings.Replace(s, "src=\"/", "src=\"" + serverIP + siteUrl + "/" , -1)


          body = []byte(s)


          w.Header().Set("GO!", "Unleashed!")
          w.WriteHeader(200)
          w.Write(body)
        }
    })

    http.ListenAndServe(serverPort, nil)
}


func hello(w http.ResponseWriter, r *http.Request) {
        dat, err := ioutil.ReadFile("home.html")
        check(err)
        w.Write(dat)
}
