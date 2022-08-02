package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	ar "github.com/apiresponse"
	"github.com/google/go-cmp/cmp"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func TestHelloHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(helloHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `<h1>Hello user</h1>`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestUserHandler(t *testing.T) {
	s := httptest.NewServer(handlers())
	defer s.Close()
	// userApi = apiUser{apiUrl: fmt.Sprintf("%s/reqres", s.URL)}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/reqres", s.URL), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	testHtml, _ := ioutil.ReadFile("test-data/test.html")
	expected := string(testHtml)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func handlers() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/reqres", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		d, _ := ioutil.ReadFile("test-data/test_user.json")
		if _, err := fmt.Fprint(w, string(d)); err != nil {
			log.Println(err)
		}
	})

	return r
}

func TestGetUsers(t *testing.T) {
	testUserSet := ar.UserJson{
		Page:       1,
		PerPage:    6,
		Total:      1,
		TotalPages: 1,
		UsersData: []ar.UserData{
			ar.UserData{
				Id:        1,
				Email:     "george.bluth@reqres.in",
				FirstName: "George",
				LastName:  "Bluth",
				Avatar:    "https://reqres.in/img/faces/1-image.jpg",
			},
		},
		Support: ar.Metadata{
			Url:  "https://reqres.in/#support-heading",
			Text: "To keep ReqRes free, contributions towards server costs are appreciated!",
		},
	}
	s := httptest.NewServer(handlers())
	defer s.Close()
	users := &ar.UserJson{}
	getUsers(users, fmt.Sprintf("%s/reqres", s.URL))
	if !cmp.Equal(users, &testUserSet) {
		t.Errorf("handler returned unexpected user set: got %v want %v",
			testUserSet, users)
	}
}
