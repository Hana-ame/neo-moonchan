package nft

import (
	tools "github.com/Hana-ame/neo-moonchan/Tools"
	"github.com/gin-gonic/gin"
)

func conv2Edge(url string) (string, error) {
	output, err := tools.Command("/home/lumin/miniconda3/bin/py", "cov_to_edge.py", "https://upload.moonchan.xyz/api/01LLWEUU7IDGWDORVQTRB3ZBUAHWUZUT4C/ss_03209bd4cde06cec229f73f084efabbe62373bd7.1920x1080.jpg")
	return output, err
}

func Conv2Edge(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(400, gin.H{
			"error": "url is required",
		})
		return
	}
	output, err := conv2Edge(url)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"url": output,
	})
}
