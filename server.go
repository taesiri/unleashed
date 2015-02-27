package main

import (
        "io"
        "log"
        "io/ioutil"
        "net/http"
        "net/url"
        "strings"
)

const serverIP = "http://127.0.0.1:8000/"
const serverPort = ":8000"

func check(e error) {
    if e != nil {
        panic(e)
    }
}


func hello(w http.ResponseWriter, r *http.Request) {
        dat, err := ioutil.ReadFile("home.html")
        check(err)
        io.WriteString(w, string(dat))

}

var mux map[string]func(http.ResponseWriter, *http.Request)

func main() {
        server := http.Server{
                Addr:    serverPort,
                Handler: &myHandler{},
        }

        mux = make(map[string]func(http.ResponseWriter, *http.Request))
        mux["/"] = hello

        server.ListenAndServe()
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        if h, ok := mux[r.URL.String()]; ok {
                h(w, r)
                return
        }

        address := r.URL.String()
        address = address[1:len(address)]

        u, err := url.Parse(address)
            if err != nil {
                panic(err)
            }



        siteUrl := "";

        if u.Host != "" {
          siteUrl = u.Scheme + "://" + u.Host
        }

        log.Printf(siteUrl)

        resp, err := http.Get(address)
        if err != nil {
                w.WriteHeader(200)

        } else {
                defer resp.Body.Close()
                body, err := ioutil.ReadAll(resp.Body)
                if err !=nil {

                }

                if siteUrl != "" {

                  s := string(body[:])

                  s = strings.Replace(s, "<img src=\"/", "<img src=\"" + serverIP + siteUrl + "/" , -1)
                  s = strings.Replace(s, "url('/", "url('" + serverIP + siteUrl + "/" , -1)


                  s = strings.Replace(s, "href=\"/", "href=\"" + serverIP + siteUrl + "/" , -1)


                  body = []byte(s)


                }

                w.Header().Set("GO!", "Unleashed!")
                w.WriteHeader(200)
                w.Write(body)
        }

}
