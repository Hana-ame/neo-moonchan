package register

import (
	"fmt"
	"testing"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
)

func TestHash(t *testing.T) {
	fmt.Println(tools.Hash("123456", "tsuki"))
}
