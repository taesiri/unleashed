package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

const serverIP = "http://127.0.0.1:8000/"
const serverPort = ":8000"
const encryptedLink = "http://127.0.0.1:8000/enc/"

var linksRegexp = regexp.MustCompile("\"(http|https)://([a-zA-Z0-9+&%=#.(){};:,.<>_+?|\\\\/\\-]*)\"")

var implicitLinksRegexp = regexp.MustCompile("\"/([a-zA-Z0-9+&=%#.(){};:,.<>_+?|\\\\/\\-]*)\"") // href=//resource
var implicitLinks2Regexp = regexp.MustCompile("\\(/[a-zA-Z0-9+&=%#.{};:,.<>_+?|\\\\/\\-]*\\)")  //Url(/resouce)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	cookieJar, _ := cookiejar.New(nil)

	http.HandleFunc("/", hello)

	http.HandleFunc("/enc/", func(w http.ResponseWriter, r *http.Request) {
		theUrl := strings.SplitN(r.URL.Path, "/", 3)[2]

		r.ParseForm()
		log.Println(r.PostForm)

		sDec, _ := base64.StdEncoding.DecodeString(theUrl)
		sDec, _ = base64.StdEncoding.DecodeString(string(sDec[:]))

		address := string(sDec[:])
		log.Printf(address)

		client := &http.Client{
			Jar: cookieJar,
		}

		var resp *http.Response
		if len(r.PostForm) == 0 {
			resp, _ = client.Get(address)
		} else {
			resp, _ = client.PostForm(address, r.PostForm)
		}

		body, err := ioutil.ReadAll(resp.Body)
		check(err)

		s := string(body[:])

		output := linksRegexp.FindAllString(s, -1)
		for _, link := range output {
			secureLink := base64.StdEncoding.EncodeToString([]byte(link[1 : len(link)-1]))
			secureLink = base64.StdEncoding.EncodeToString([]byte(secureLink))

			s = strings.Replace(s, link[1:len(link)-1], encryptedLink+secureLink, -1)
		}

		u, err := url.Parse(address)
		check(err)

		if u.Host != "" {
			siteUrl := u.Scheme + "://" + u.Host
			log.Println(siteUrl)

			output = implicitLinksRegexp.FindAllString(s, -1)
			for _, link := range output {

				fullLink := siteUrl + "//" + link[1:len(link)-1]

				secureLink := base64.StdEncoding.EncodeToString([]byte(fullLink))
				secureLink = base64.StdEncoding.EncodeToString([]byte(secureLink))

				s = strings.Replace(s, link[1:len(link)-1], encryptedLink+secureLink, -1)
			}

			output = implicitLinks2Regexp.FindAllString(s, -1)
			for _, link := range output {

				fullLink := siteUrl + "//" + link[1:len(link)-1]

				secureLink := base64.StdEncoding.EncodeToString([]byte(fullLink))
				secureLink = base64.StdEncoding.EncodeToString([]byte(secureLink))

				s = strings.Replace(s, link[1:len(link)-1], encryptedLink+secureLink, -1)
			}
		}

		body = []byte(s)

		w.Header().Set("GO!", "Unleashed!")
		w.WriteHeader(200)
		w.Write(body)

	})

	http.ListenAndServe(serverPort, nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	dat, err := ioutil.ReadFile("home.html")
	check(err)
	w.Write(dat)
}
