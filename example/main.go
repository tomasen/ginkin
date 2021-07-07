package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tomasen/ginkin"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"net/http"
)

func main() {
	apis := map[string]ginkin.APIHandler{
		"version":        {"GET", DescribeVersion, "print version info"},
		"user/list":      {"POST", UserList, "list users"},
		"user/:user":     {"GET", DescribeUser, "print user info"},
		"user/:user#del": {"DELETE", DeleteUser, "print user info"},
		"user/add":       {"PUT", AddUsers, "add users"},
	}

	// prepare gin engine
	router := gin.Default()

	// add other flag or more command to kingpin
	kingpin.HelpFlag.Short('h')

	gk := &ginkin.GinKin{
		APIs:     apis,
		Start:    ServeGin,
		Fallback: CLIFallback,
	}
	gk.Run(router, "/")
}

func UserList(c *gin.Context) {
	c.JSON(http.StatusOK, "john")
}

func DescribeUser(c *gin.Context) {
	user, exist := c.Params.Get("user")
	if !exist {
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, user)
}

func AddUsers(c *gin.Context) {
	var users []string
	if err := c.ShouldBindJSON(&users); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func DeleteUser(c *gin.Context) {
	user, _ := c.Params.Get("user")

	log.Println("deleting:", user)
	c.Status(http.StatusOK)
}

func DescribeVersion(c *gin.Context) {
	c.JSON(http.StatusOK, "0.1")
}

func ServeGin(router *gin.Engine) {
	router.Run(":3000")
}

func CLIFallback(cmd string) {
	log.Println("unhandled command line action:", cmd)
}
