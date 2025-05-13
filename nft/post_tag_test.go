package nft

import (
	"fmt"
	"testing"
)

func TestSetPostTag(t *testing.T) {
	var err error
	err = addTagToPost("114485557672017920", "tag")
	fmt.Println(err)
	err = addTagToPost("114485557671821312", "tag2")
	fmt.Println(err)

}

func TestPatchOwner(t *testing.T) {
	var err error
	err = patchOwnerOfPost("114496263078674432", "456", 999, true, false)
	fmt.Println(err)
	err = patchOwnerOfPost("114496263078674432", "999", -999, false, true)
	fmt.Println(err)
}

func TestChangeOwner(t *testing.T) {
	var err error
	err = changeOwnerOfPost("114496263078674432", "testuser")
	fmt.Println(err)
}
