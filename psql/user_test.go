package psql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	tools "github.com/Hana-ame/neo-moonchan/Tools"
	"github.com/Hana-ame/neo-moonchan/Tools/db"
	"github.com/Hana-ame/neo-moonchan/Tools/orderedmap"
)

func TestSaveUser(t *testing.T) {

	s := `
	{
		"@context": [
			"https://www.w3.org/ns/activitystreams",
			"https://w3id.org/security/v1",
			{
				"manuallyApprovesFollowers": "as:manuallyApprovesFollowers",
				"toot": "http://joinmastodon.org/ns#",
				"featured": {
					"@id": "toot:featured",
					"@type": "@id"
				},
				"featuredTags": {
					"@id": "toot:featuredTags",
					"@type": "@id"
				},
				"alsoKnownAs": {
					"@id": "as:alsoKnownAs",
					"@type": "@id"
				},
				"movedTo": {
					"@id": "as:movedTo",
					"@type": "@id"
				},
				"schema": "http://schema.org#",
				"PropertyValue": "schema:PropertyValue",
				"value": "schema:value",
				"discoverable": "toot:discoverable",
				"Device": "toot:Device",
				"Ed25519Signature": "toot:Ed25519Signature",
				"Ed25519Key": "toot:Ed25519Key",
				"Curve25519Key": "toot:Curve25519Key",
				"EncryptedMessage": "toot:EncryptedMessage",
				"publicKeyBase64": "toot:publicKeyBase64",
				"deviceId": "toot:deviceId",
				"claim": {
					"@type": "@id",
					"@id": "toot:claim"
				},
				"fingerprintKey": {
					"@type": "@id",
					"@id": "toot:fingerprintKey"
				},
				"identityKey": {
					"@type": "@id",
					"@id": "toot:identityKey"
				},
				"devices": {
					"@type": "@id",
					"@id": "toot:devices"
				},
				"messageFranking": "toot:messageFranking",
				"messageType": "toot:messageType",
				"cipherText": "toot:cipherText",
				"suspended": "toot:suspended",
				"focalPoint": {
					"@container": "@list",
					"@id": "toot:focalPoint"
				}
			}
		],
		"id": "https://fedi.moonchan.xyz/users/nanakananoka",
		"type": "Person",
		"following": "https://fedi.moonchan.xyz/users/nanakananoka/following",
		"followers": "https://fedi.moonchan.xyz/users/nanakananoka/followers",
		"inbox": "https://fedi.moonchan.xyz/users/nanakananoka/inbox",
		"outbox": "https://fedi.moonchan.xyz/users/nanakananoka/outbox",
		"featured": "https://fedi.moonchan.xyz/users/nanakananoka/collections/featured",
		"featuredTags": "https://fedi.moonchan.xyz/users/nanakananoka/collections/tags",
		"preferredUsername": "nanakananoka",
		"name": "ななかなのか？！",
		"summary": "<p>头像是性癖<br />头图是自己</p>",
		"url": "https://fedi.moonchan.xyz/@nanakananoka",
		"manuallyApprovesFollowers": false,
		"discoverable": false,
		"published": "2024-05-18T00:00:00Z",
		"devices": "https://fedi.moonchan.xyz/users/nanakananoka/collections/devices",
		"publicKey": {
			"id": "https://fedi.moonchan.xyz/users/nanakananoka#main-key",
			"owner": "https://fedi.moonchan.xyz/users/nanakananoka",
			"publicKeyPem": "-----BEGIN PUBLIC KEY-----\nMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA0FneAc0XPrwzvskqqujJ\n2Wf0bLAPkCVcjAC33zmomviI3Cyf144oKAxJNsp/yciWXw9J5PNuXm5LGe7ZN7Qb\nMWExRAbv09o8iahoFUtcqpoMxpDiiqCTYs35BhudqZqKtO9hteCV5taR2oxDhbqu\n7FDY9hV98CuVJ5iMT0TOQwCuFaiy354B3lxuZFZOST8pLDvo6kuM4mEaxikqNlCJ\ntoylHR97aqYUoRYyw9fq+CVBb61HERGrH577rrRegyEoirf78xChrPLXVaPHK+Qf\nhFy9UuiITBVIH4Na2e4V05Y5ay2jeFR5QNLygU/UxGAbKGD2Nd2a29EnDpOEnZGZ\nPgiTSF33Gy31H8JSb0f0OVFeHQ2ECPuaZoKBvGmnnXyOZX3cxTLuqN/NWmRm3Skp\nwOPAUQ86E4elY/cNlvVTJfgG7EuOF+N4SX23tJyenajy4kjxLeuu3A8DkxT/zswO\nWnrlL7hQeNbnf3U8ZEX8n2XDCjgwhUF0KiewCpNSz7NOqGHPC1ELY8ntilZ2sDEz\n/ofVwBUPC0peVo0WcwjmkFp9w6mZ50qxNQY1Jg0AEJM/16OKW7YetmZMEjqAXQ7l\njw/rSWjmbgk0/d/8VeurtU8YWHi8aTvSIkiWxp4eNmWLTtC+4cVW/ZTMwDbiTpJK\nhhSuyqemFuFshy/HUOkaYJkCAwEAAQ==\n-----END PUBLIC KEY-----\n"
		},
		"privateKeyPem":"-----BEGIN PRIVATE KEY-----\n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\n+03FBt4grFlYR36/fru5BEnuJlXcLYc=\n-----END PRIVATE KEY-----\n",
		"tag": [],
		"attachment": [],
		"endpoints": {
			"sharedInbox": "https://fedi.moonchan.xyz/inbox"
		},
		"icon": {
			"type": "Image",
			"mediaType": "image/png",
			"url": "https://media.mstdn.jp/accounts/avatars/112/461/152/626/039/703/original/f4a814f00c18cd26.png"
		},
		"image": {
			"type": "Image",
			"mediaType": "image/jpeg",
			"url": "https://media.mstdn.jp/accounts/headers/112/461/152/626/039/703/original/9f5bab767f21a76b.jpg"
		}
	}`
	o := orderedmap.New()
	json.Unmarshal([]byte(s), &o)
	fmt.Println(o.Keys())
	tools.MoveToFirstInPlace(o.Keys(), "@context")
	fmt.Println(o.Keys())
	e := db.Exec(func(tx *sql.Tx) error {
		e := SaveUser(tx, "https://fedi.moonchan.xyz/users/nanakananoka", o)
		fmt.Println(e)
		return tx.Commit()
	})
	fmt.Println("??", e)
}

func TestAddSaveUser(t *testing.T) {

	nanakananoka := "nanakananoka3"

	s := `
	{
		"@context": [
			"https://www.w3.org/ns/activitystreams",
			"https://w3id.org/security/v1",
			{
				"manuallyApprovesFollowers": "as:manuallyApprovesFollowers",
				"toot": "http://joinmastodon.org/ns#",
				"featured": {
					"@id": "toot:featured",
					"@type": "@id"
				},
				"featuredTags": {
					"@id": "toot:featuredTags",
					"@type": "@id"
				},
				"alsoKnownAs": {
					"@id": "as:alsoKnownAs",
					"@type": "@id"
				},
				"movedTo": {
					"@id": "as:movedTo",
					"@type": "@id"
				},
				"schema": "http://schema.org#",
				"PropertyValue": "schema:PropertyValue",
				"value": "schema:value",
				"discoverable": "toot:discoverable",
				"Device": "toot:Device",
				"Ed25519Signature": "toot:Ed25519Signature",
				"Ed25519Key": "toot:Ed25519Key",
				"Curve25519Key": "toot:Curve25519Key",
				"EncryptedMessage": "toot:EncryptedMessage",
				"publicKeyBase64": "toot:publicKeyBase64",
				"deviceId": "toot:deviceId",
				"claim": {
					"@type": "@id",
					"@id": "toot:claim"
				},
				"fingerprintKey": {
					"@type": "@id",
					"@id": "toot:fingerprintKey"
				},
				"identityKey": {
					"@type": "@id",
					"@id": "toot:identityKey"
				},
				"devices": {
					"@type": "@id",
					"@id": "toot:devices"
				},
				"messageFranking": "toot:messageFranking",
				"messageType": "toot:messageType",
				"cipherText": "toot:cipherText",
				"suspended": "toot:suspended",
				"focalPoint": {
					"@container": "@list",
					"@id": "toot:focalPoint"
				}
			}
		],
		"id": "https://fedi.moonchan.xyz/users/` + nanakananoka + `",
		"type": "Person",
		"following": "https://fedi.moonchan.xyz/users/` + nanakananoka + `/following",
		"followers": "https://fedi.moonchan.xyz/users/` + nanakananoka + `/followers",
		"inbox": "https://fedi.moonchan.xyz/users/` + nanakananoka + `/inbox",
		"outbox": "https://fedi.moonchan.xyz/users/` + nanakananoka + `/outbox",
		"featured": "https://fedi.moonchan.xyz/users/` + nanakananoka + `/collections/featured",
		"featuredTags": "https://fedi.moonchan.xyz/users/` + nanakananoka + `/collections/tags",
		"preferredUsername": "` + nanakananoka + `",
		"name": "ななかなのか？！",
		"summary": "<p>头像是性癖<br />头图是自己</p>",
		"url": "https://fedi.moonchan.xyz/@` + nanakananoka + `",
		"manuallyApprovesFollowers": false,
		"discoverable": false,
		"published": "2024-05-18T00:00:00Z",
		"devices": "https://fedi.moonchan.xyz/users/` + nanakananoka + `/collections/devices",
		"publicKey": {
			"id": "https://fedi.moonchan.xyz/users/` + nanakananoka + `#main-key",
			"owner": "https://fedi.moonchan.xyz/users/` + nanakananoka + `",
			"publicKeyPem": "-----BEGIN PUBLIC KEY-----\nMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA0FneAc0XPrwzvskqqujJ\n2Wf0bLAPkCVcjAC33zmomviI3Cyf144oKAxJNsp/yciWXw9J5PNuXm5LGe7ZN7Qb\nMWExRAbv09o8iahoFUtcqpoMxpDiiqCTYs35BhudqZqKtO9hteCV5taR2oxDhbqu\n7FDY9hV98CuVJ5iMT0TOQwCuFaiy354B3lxuZFZOST8pLDvo6kuM4mEaxikqNlCJ\ntoylHR97aqYUoRYyw9fq+CVBb61HERGrH577rrRegyEoirf78xChrPLXVaPHK+Qf\nhFy9UuiITBVIH4Na2e4V05Y5ay2jeFR5QNLygU/UxGAbKGD2Nd2a29EnDpOEnZGZ\nPgiTSF33Gy31H8JSb0f0OVFeHQ2ECPuaZoKBvGmnnXyOZX3cxTLuqN/NWmRm3Skp\nwOPAUQ86E4elY/cNlvVTJfgG7EuOF+N4SX23tJyenajy4kjxLeuu3A8DkxT/zswO\nWnrlL7hQeNbnf3U8ZEX8n2XDCjgwhUF0KiewCpNSz7NOqGHPC1ELY8ntilZ2sDEz\n/ofVwBUPC0peVo0WcwjmkFp9w6mZ50qxNQY1Jg0AEJM/16OKW7YetmZMEjqAXQ7l\njw/rSWjmbgk0/d/8VeurtU8YWHi8aTvSIkiWxp4eNmWLTtC+4cVW/ZTMwDbiTpJK\nhhSuyqemFuFshy/HUOkaYJkCAwEAAQ==\n-----END PUBLIC KEY-----\n"
		},
		"privateKeyPem":"-----BEGIN PRIVATE KEY-----\n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\n+03FBt4grFlYR36/fru5BEnuJlXcLYc=\n-----END PRIVATE KEY-----\n",
		"tag": [],
		"attachment": [],
		"endpoints": {
			"sharedInbox": "https://fedi.moonchan.xyz/inbox"
		},
		"icon": {
			"type": "Image",
			"mediaType": "image/png",
			"url": "https://media.mstdn.jp/accounts/avatars/112/461/152/626/039/703/original/f4a814f00c18cd26.png"
		},
		"image": {
			"type": "Image",
			"mediaType": "image/jpeg",
			"url": "https://media.mstdn.jp/accounts/headers/112/461/152/626/039/703/original/9f5bab767f21a76b.jpg"
		}
	}`
	o := orderedmap.New()
	json.Unmarshal([]byte(s), &o)
	fmt.Println(o.Keys())
	tools.MoveToFirstInPlace(o.Keys(), "@context")
	fmt.Println(o.Keys())
	e := db.Exec(func(tx *sql.Tx) error {
		e := SaveUser(tx, `https://fedi.moonchan.xyz/users/`+nanakananoka+``, o)
		fmt.Println(e)
		return tx.Commit()
	})
	fmt.Println("??", e)
}
