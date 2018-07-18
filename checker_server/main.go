/*
This is a simple server that returns the real ip behind the proxy,
used for checking whether a proxy is high anonymous or not
*/

package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/checkproxy", func(c *gin.Context) {
		/*
			ClientIP implements a best effort algorithm to return the real client IP,
			it parses X-Real-IP & X-Forwarded-For in order to work properly with reverse-proxies.
			Use X-Forwarded-For before X-Real-Ip as nginx uses X-Real-Ip with the proxy's IP.
		*/
		realIP := c.ClientIP()
		c.String(200, realIP)
	})

	r.Run(":80") // listen and serve on 0.0.0.0:8080
}
