package action

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestWebfinger(t *testing.T) {
	// id, err := FetchWebfinger("nanakananoka@mstdn.jp")
	id, err := FetchWebfinger("haruurara@wxw.moe")
	fmt.Println(id, err)
}

// func TestFetch(t *testing.T) {
// 	resp, err := myfetch.Fetch("GET", "https://mstdn.jp/.well-known/webfinger?resource=acct:nanakananoka@mstdn.jp", nil, nil)
// 	log.Println(err)
// 	resp.Body
// }

func TestUser(t *testing.T) {
	fmt.Println(3)
	user, err := FetchUser("https://wxw.moe/users/HaruUrara")
	fmt.Println(err)
	j, _ := json.Marshal(user)
	fmt.Println(string(j))
}
