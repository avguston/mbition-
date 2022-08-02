package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"text/template"

	ar "github.com/apiresponse"
	log "github.com/sirupsen/logrus"
)

const (
	APIURL    = "https://reqres.in/api/users"
	INDEXPAGE = "templates/index.html"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	log.Info("helloHendler /")
	fmt.Fprintf(w, "<h1>Hello user</h1>")
}

func userHandler(w http.ResponseWriter, _ *http.Request) {
	log.Info("userHendle /reqres")
	users := &ar.UserJson{}
	getUsers(users, APIURL)
	t := renderHtml(INDEXPAGE)
	t.Execute(w, users)
}

func getUsers(j *ar.UserJson, api_url string) error {
	res, err := http.Get(api_url)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Invalide server response")
		return err
	}
	body, err := io.ReadAll(res.Body)
	// TODO res.Body could whatever but not a JSON
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Cannot read a res.Body")
		return err
	}
	userPage, err := j.GetJson(body)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Invalide json format")
		return err
	}
	msg := fmt.Sprintf("parcing %sth page", getPageNum(api_url))
	log.WithFields(log.Fields{"user_list": userPage.UsersData}).Infof(msg)
	if userPage.Page < userPage.TotalPages {
		userPage.Page += 1
		getUsers(j, fmt.Sprintf("%s?page=%d", APIURL, userPage.Page))
	}
	return nil
}

func getPageNum(uri string) string {
	//TODO Never trust bloody users they'd send anything exept expecting data
	fullUrl, err := url.Parse(uri)
	if err != nil {
		log.Error(err)
	}
	queryMark, err := url.ParseQuery(fullUrl.RawQuery)
	if err != nil {
		log.Error(err)
	}
	if v, ok := queryMark["page"]; ok {
		return v[0]
	}
	return "1"
}

func renderHtml(path string) *template.Template {
	t := template.New(filepath.Base(path))
	t.Funcs(template.FuncMap{"mod": func(i, j int) bool { return i%j == 0 }})
	t.Funcs(template.FuncMap{"inc": func(i, j int) int { return i + j }})
	t.ParseFiles(path)
	return t
}

func main() {
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/reqres", userHandler)
	log.Info("server is up on http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
