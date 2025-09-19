package action

import (
	"bytes"
	"crypto"
	"database/sql"
	"net/http"
	"time"

	"github.com/Hana-ame/httpsig"
	tools "github.com/Hana-ame/neo-moonchan/Tools"
	db "github.com/Hana-ame/neo-moonchan/Tools/db/pq"
	myfetch "github.com/Hana-ame/neo-moonchan/Tools/my_fetch"
	"github.com/Hana-ame/neo-moonchan/psql"
)

// pubKeyOwner is ID with out '#main-key'
func FetchWithSign(
	id string,
	method string, endpoint string, body []byte,
) (
	*http.Response, error,
) {
	var privateKey crypto.PrivateKey
	err := db.Exec(func(tx *sql.Tx) error {
		user, err := psql.ReadUser(tx, id)
		if err != nil {
			return err
		}

		privateKeyPem, err := tools.Extract[string](user, "privateKeyPem")
		if err != nil {
			return err
		}
		privateKey, err = tools.ParsePrivateKey([]byte(privateKeyPem))
		if err != nil {
			return err
		}
		// 这里不写，在别的地方写入之后提示这里去fetch
		return tx.Commit()
	})
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(method, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	r.Header.Set("host", r.URL.Host)
	r.Header.Set("date", time.Now().UTC().Format(http.TimeFormat))
	r.Header.Set("content-type", "application/activity+json")

	err = Sign(privateKey, id+"#main-key", r, body)
	if err != nil {
		return nil, err
	}

	return myfetch.DefaultFetcher.Do(r)

}

// usage:
// fill the inputs
func Sign(privateKey crypto.PrivateKey, pubKeyID string, r *http.Request, body []byte) error {
	prefs := []httpsig.Algorithm{httpsig.RSA_SHA256}
	digestAlgorithm := httpsig.DigestSha256
	// The "Date" and "Digest" headers must already be set on r, as well as r.URL.
	headersToSign := []string{httpsig.RequestTarget, "host", "date", "digest", "content-type"}
	signer, chosenAlgo, err := httpsig.NewSigner(prefs, digestAlgorithm, headersToSign, httpsig.Signature, 1<<16)
	// log.Println(chosenAlgo)
	_ = chosenAlgo
	if err != nil {
		return err
	}
	// To sign the digest, we need to give the signer a copy of the body...
	// ...but it is optional, no digest will be signed if given "nil"
	// body := []byte{}
	// log.Println(string(body))

	// If r were a http.ResponseWriter, call SignResponse instead.
	return signer.SignRequest(privateKey, pubKeyID, r, body)
}
