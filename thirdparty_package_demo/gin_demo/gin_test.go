package gin_demo

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
)

func TestBasicEx(t *testing.T) {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	r.Run(":8088")
}

func TestLoadHTML(t *testing.T) {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	r.Run(":8088")
}

func TestJSONP(t *testing.T) {
	r := gin.Default()
	r.GET("/JSONP", func(c *gin.Context) {
		data := map[string]interface{}{
			"foo": "bar",
		}

		// /JSONP?callback=x
		// output: x({\"foo\":\"bar\"})
		c.JSONP(http.StatusOK, data)
	})

	r.Run(":8088")
}

type LoginForm struct {
	User     string `form:"user" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func TestMultipart(t *testing.T) {
	r := gin.Default()
	r.POST("/login", func(c *gin.Context) {
		var form LoginForm
		if c.ShouldBind(&form) == nil {
			if form.User == "user" && form.Password == "password" {
				c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			}
		}
	})
	r.Run(":8088")
}

func TestFormParse(t *testing.T) {
	r := gin.Default()
	r.POST("/post_form", func(c *gin.Context) {
		msg := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous")

		c.JSON(http.StatusOK, gin.H{
			"status":  "posted",
			"message": msg,
			"nick":    nick,
		})
	})
	r.Run(":8088")
}
