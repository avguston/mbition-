package apiresponse

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type UserData struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
}

type Metadata struct {
	Url  string `json:"url"`
	Text string `json:"text"`
}

type UserJson struct {
	Page       int        `json:"page"`
	PerPage    int        `json:"per_page"`
	Total      int        `json:"total"`
	TotalPages int        `json:"total_pages"`
	UsersData  []UserData `json:"data"`
	Support    Metadata   `json:"support"`
}

func (jj *UserJson) GetJson(r []byte) (*UserJson, error) {
	// TODO naming is shit
	// TODO what are you going to do if a list of users will be tonish
	jj.TotalPages = 1
	jj.Page = 1
	j := &UserJson{}
	// err := json.NewDecoder(r).Decode(j)
	err := json.Unmarshal(r, j)
	if err != nil {
		log.WithFields(log.Fields{"origin_msg": err}).Error("Invalide json format")
		return nil, err
	}
	jj.PerPage += j.PerPage
	jj.Total += j.Total
	jj.Support = j.Support
	for _, user := range j.UsersData {
		jj.UsersData = append(jj.UsersData, user)
	}
	return j, nil
}

// func main() {
// 	users := &UserJson{}
// 	d, _ := ioutil.ReadFile("user1.json")
// 	users.getJson(d)
// 	d, _ = ioutil.ReadFile("user2.json")
// 	users.getJson(d)
// 	fmt.Println("==========================")
// 	fmt.Println(users)
// }
