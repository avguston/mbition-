package apiresponse

import (
	"github.com/google/go-cmp/cmp"
	"io/ioutil"
	"testing"
)

func TestGetJson(t *testing.T) {
	users := &UserJson{}
	d, _ := ioutil.ReadFile("../../test-data/test_user.json")
	users.GetJson(d)
	testUserSet := UserJson{
		Page:       1,
		PerPage:    6,
		Total:      1,
		TotalPages: 1,
		UsersData: []UserData{
			UserData{
				Id:        1,
				Email:     "george.bluth@reqres.in",
				FirstName: "George",
				LastName:  "Bluth",
				Avatar:    "https://reqres.in/img/faces/1-image.jpg",
			},
		},
		Support: Metadata{
			Url:  "https://reqres.in/#support-heading",
			Text: "To keep ReqRes free, contributions towards server costs are appreciated!",
		},
	}
	if !cmp.Equal(users, &testUserSet) {
		t.Errorf("handler returned unexpected user set: got %v want %v",
			testUserSet, users)
	}
}
