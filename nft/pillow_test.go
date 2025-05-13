package nft

import (
	"fmt"
	"testing"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
)

func TestCommand(t *testing.T) {
	fmt.Println(2)
	output, err := tools.Command("/home/lumin/miniconda3/bin/py", "cov_to_edge.py", "https://upload.moonchan.xyz/api/01LLWEUU7IDGWDORVQTRB3ZBUAHWUZUT4C/ss_03209bd4cde06cec229f73f084efabbe62373bd7.1920x1080.jpg")
	fmt.Println(err)
	fmt.Println(err)
	fmt.Println(output)
	fmt.Println(output)
}
