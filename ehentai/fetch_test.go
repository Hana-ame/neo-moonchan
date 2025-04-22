package ehentai

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	myfetch "github.com/Hana-ame/neo-moonchan/Tools/my_fetch"
)

func TestPost(t *testing.T) {
	resp, err := myfetch.Fetch(http.MethodPost, "https://exhentai.org/archiver.php?gid=3325143&token=12179143e2", map[string][]string{"Cookie": {"ipb_member_id=3096156; ipb_pass_hash=7689f8ad48f6576453620a52ffa238e1; sl=dm_1; sk=939cq4lc9fccsgkh6kj8gpiyby6a; hath_perks=m1-0266e6109f; igneous=r52c4sygccj3w51gx"}}, strings.NewReader("dltype=org&dlcheck=Download+Original+Archive"))
	fmt.Println(err)
	body, err := io.ReadAll(resp.Body)
	fmt.Println(err)
	fmt.Println(string(body))

}
