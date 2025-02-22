package inbox

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Hana-ame/neo-moonchan/Tools/orderedmap"
	"github.com/gin-gonic/gin"
)

func Inbox(c *gin.Context) {
	name := c.Param("name")
	host := c.GetHeader("Host")

	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	o := orderedmap.New()
	if err := json.Unmarshal(b, &o); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// for debug
	log.Println(name, host)
	log.Println(string(b))
	h, _ := json.Marshal(c.Request.Header)
	log.Println(string(h))
	j, _ := json.Marshal(o)
	log.Println(string(j))

	// err = verify(c, b)
	if err != nil {
		log.Println(err)
		// c.JSON(http.StatusUnauthorized, err.Error())
		// return
	}
	// if err := handler.Inbox(o, name, host, err); err != nil {
	// 	log.Println(err)
	// 	c.JSON(http.StatusUnauthorized, err.Error())
	// 	return
	// }

	c.AbortWithStatus(http.StatusOK)
}
