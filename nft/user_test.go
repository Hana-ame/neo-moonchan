package nft

import (
	"fmt"
	"testing"
)

func TestSetUser(t *testing.T) {
	var err error
	err = setUserProfile("testuser", map[string]any{"test": "test"})
	fmt.Println(err)
	setUserProfile("testuser", map[string]any{"test": "test2"})
	fmt.Println(err)
}

func TestGetUser(t *testing.T) {
	var err error
	profile, err := getUserProfile("testuser")
	fmt.Println(err)
	fmt.Println(profile)
}
