package action

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	myfetch "github.com/Hana-ame/neo-moonchan/Tools/my_fetch"
)

func TestFollow(t *testing.T) {
	fmt.Println(os.Getenv("HTTPS_PROXY"))
	pc := myfetch.NewProxyClient(os.Getenv("HTTPS_PROXY"))
	myfetch.DefaultFetcher = myfetch.NewFetcher(nil, myfetch.NewClientPool([]*http.Client{pc}))
	myfetch.SetDefaultHeader(http.Header{"User-Agent": []string{"MyFetch/1.1.0"}})
	// resp, _ := myfetch.Fetch("GET", "https://echo.free.beeceptor.com", http.Header{"Accept": []string{"application/ld+json", "application/activity+json"}, "User-Agent": []string{"NeoMoonchan/1.0.9 (https://moonchan.xyz)"}, "Cache-Control": []string{"no-cache"}}, nil)
	// resp, _ := myfetch.Fetch("GET", "https://echo.free.beeceptor.com/", nil, nil)
	// defer resp.Body.Close()
	// tools.WriteReaderToFile("echo.json", resp.Body)
	fmt.Println("no  11111111")
	err := Follow("https://mstdn.work.gd/users/nanakananoka", "https://mstdn.jp/users/nanakananoda")
	fmt.Println(err)
	// resp, _ := myfetch.Fetch("GET", "https://mstdn.jp/users/nanakananoda", http.Header{"Accept": []string{"application/ld+json", "application/activity+json"}, "User-Agent": []string{"NeoMoonchan/1.0.9 (https://moonchan.xyz)"}, "Cache-Control": []string{"no-cache"}}, nil)
	// defer resp.Body.Close()
	// tools.WriteReaderToFile("echo.json", resp.Body)
	// o := orderedmap.New()
	// for k, v := range resp.Header {
	// 	o.Set(k, v)
	// }
	// o.SortKeys(sort.Strings)
	// for _, k := range o.Keys() {
	// 	fmt.Println(k, o.GetOrDefault(k, []string{}).([]string))
	// }

}

// 2025/03/20 18:04:58 {"@context":"https://www.w3.org/ns/activitystreams","id":"https://mstdn.jp/e01c9e84-6498-45a4-8b52-76c46c83b1a1","type":"Follow","actor":"https://mstdn.jp/users/nanakananoda","object":"https://mstdn.work.gd/users/nanakananoka"}
// 2025/03/20 18:04:58 {"Accept-Encoding":["gzip"],"Connection":["close"],"Content-Length":["227"],"Content-Type":["application/activity+json"],"Date":["Thu, 20 Mar 2025 10:05:00 GMT"],"Digest":["SHA-256=JIok7tlw7kz2CsYaI9gn4wDEgn2J4y30G+86O55jG4s="],"Signature":["keyId=\"https://mstdn.jp/users/nanakananoda#main-key\",algorithm=\"rsa-sha256\",headers=\"(request-target) host date digest content-type\",signature=\"R4FxEzUmFwfeFZAu9t83tKJ+TmN1heX9epDAeaZHVQ+CD+2kDOsR9zYz+U82Aj+MJ2qKukh4mnY/w+UpeW/800I8BXIE3dYm6XpbboEVsk3lQV7YLYtCyWQJRqdNh+YmhohhIVzhGxIrMphTEiCBNRHNzrcea8Oqerre33+YwtvoJ5UJwqpvsXGwcqfRNJ1l73muhPnDvCWyirgANfYq1Kp9Uf1h7RuollfWVbreuBuBI4FCYOjBJ5JWlaAJFBqXUtlwsaPYY+Xtq9wB3F4ZScQqFBvqUqleSu0JuCV7eN+dBjNBotRH7caGZGFIfNlB/c9iFb+FuM+vBpIG0AML4Q==\""],"User-Agent":["http.rb/5.1.1 (Mastodon/4.1.20; +https://mstdn.jp/)"]}
// [GIN] 2025/03/20 - 18:04:58 | 200 |   29.308747ms |       127.0.0.1 | POST     "/users/nanakananoka/inbox"
// 2025/03/20 18:05:27 {"@context":"https://www.w3.org/ns/activitystreams","id":"https://mstdn.work.gd/5e7b0184-6651-47a0-8fd6-647f8f279c62","type":"Follow","actor":"https://mstdn.work.gd/users/nanakananoka","object":"https://mstdn.work.gd/users/nanakananoka","status":"pending"}
// 2025/03/20 18:05:27 {"Accept-Encoding":["gzip"],"Connection":["close"],"Content-Length":["256"],"Content-Type":["application/activity+json"],"Date":["Thu, 20 Mar 2025 10:05:25 GMT"],"Digest":["SHA-256=5sMbSdr7ElCtqpK4SWsGl60rrNHVcCToooxQu+eC788="],"Signature":["keyId=\"https://mstdn.work.gd/users/nanakananoka#main-key\",algorithm=\"rsa-sha256\",headers=\"(request-target) host date digest content-type\",signature=\"jjg6oz/kOqhYbM2fqrmnBqccsr2PIx4geHKelzx5WgEiZNAjSUOtqWd6vUgbmyBcinJx2llghC6vfsIUoGhXf9blNelUSkdTcH4UAUfj+r7hlrv8lN8V6+8RBU3KG5TXxQARV0XYqOS8uQxUVskV/lHNYjGDdZMQf8HKI+71LWvToTDE26lvB3UDsG5ybS39tth4/wVovalsgn0P7fYLCAE9MbPW4mB/s75UxGs/G3judvLQ40LWlvC0zOflHgZaFqXopMltUPEIixHzy/TqJJQdZ8UJgL3CE9nUftr/yImanw4YuoW/nISDu2dTrj1nekoVxHgybQIht3BrZu4orV7pLM0w7rLQU1/eqVUDkIAZO8p3zh1XPjTgI65oEl4XhqUu6+DAAgy7jBo4r6okpfIk8MGLWqGNJDR4SUaAoJS+DFg02vQ2pCUfpYjOsSc+YIGFBYodoXVS5cmyJpw5nGUtdP40ct9zUHjun13bCugmKKjh+VhH1NVF5J0SABvcXRjyusiUGwZRF9Ko7KeZLuqRC6ahM9s+lshy8egWh/EjjqBISu6+J9plFDzjaEBxAcwHNcR1JTdOrYRUMCgyx9/ka8gsi1HtP6l28sqWiVrxcTfAu+76x3E8jViOgczOAWdnBxusbI2zBuUp0Tmcz+M6+HHXAJdOUBDvIsuq2hE=\""],"User-Agent":["Go-http-client/2.0"]}

// 2025/03/20 18:09:53 {"@context":"https://www.w3.org/ns/activitystreams","id":"https://mstdn.work.gd/b0f9f1d6-0553-4213-a041-70c9f4cf7969","type":"Follow","actor":"https://mstdn.work.gd/users/nanakananoka","object":"https://mstdn.work.gd/users/nanakananoka2","status":"pending"}
// 2025/03/20 18:09:53 {"Accept-Encoding":["gzip"],"Connection":["close"],"Content-Length":["257"],"Content-Type":["application/activity+json"],"Date":["Thu, 20 Mar 2025 10:09:52 GMT"],"Digest":["SHA-256=XVyCsw93EpBnsCSF8c0AUmKbWxeSA5osm5wYSEgq+Dk="],"Signature":["keyId=\"https://mstdn.work.gd/users/nanakananoka#main-key\",algorithm=\"rsa-sha256\",headers=\"(request-target) host date digest content-type\",signature=\"t1gsI0i2fHvaV4+sCovKh+BTgtTi8N/gYocLx9DKM44oUsBiUQgZ0Smqn8+GSPuIuUmW0ii04bcCkqnWHH91MFYD6oQ7/TYPgAfLXliQjOkBE1EcYSdGM0DNyJmqiwX565Xlx/iRJ5cBEilOpm26vkUbFgsi/YUzNJ+wXIkhDDZpLvLSL2JG5tBk1XIU+N0HWvBax79YATh5M/V8IW0EL5223sf7m8SKIrLM/XwP/BeIhjb5iTdovna11KLyAixDW0XsfmZt4+mmzVf61VFpML5l68AylTjjQ1+1unHrW5RCGYSf3q8AcUEUw3qw5DYSR5NjJzJnhRKTvh177XNhYjuzUFcCqPCPEr+z2AdjOcGDhec+j3v0osLzamZP2z7GySPsSfyrtspBlX7sHjDNyCM7RL/Xxf5pWvfUe6U75UKcQVFtaYZS0ghFd1LdJ6DJiJ3+BOi1syIbqvNAg6nQQUZ63ugVNq9sftKI3VBW5K/6XFWNujo1kaCRq++GRLWmz4/cWsMA3I0E+dfKMXsdknhfpzQktkY4yzUmojlbHdd5Qa2bEjFRGYb3+nZ2VtAXXQ6ObQozI6uUlJkD/LcJbclkSW0QFx4Jx1mVLLxmBVGZ4hnGcpqp4DYNsoJvXwntG8Z1LCA19uORaQwihP8aW+uPzbJMs/1KatJxSgsgmR8=\""],"User-Agent":["Go-http-client/2.0"]}
