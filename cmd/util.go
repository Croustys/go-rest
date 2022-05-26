package main

import "github.com/gin-gonic/gin"

func set_request_provider(c *gin.Context) {
	//Building query beacuse gothic doesn't handle gin request
	q := c.Request.URL.Query()
	q.Add("provider", c.Param("provider"))
	c.Request.URL.RawQuery = q.Encode()
}
