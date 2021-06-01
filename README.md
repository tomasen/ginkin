# ginkin

Setup HTTP API server and command-line interface at the same time by mix Gin and Kingpin.

## Example

By following example. 

- `go run main.go` will start http api server handle requests to "http://:3000/version".
- `go run main.go version` will print out version info directly.
- `go run main.go user/list` will print out user list.
- `go run main.go user/:user john` will print out user list.
- `go run main.go user/add "['jane']"` will add users.

```go
func main() {
    apis := ginkin.APIHandlers{
    	"version":    {"GET", DescribeVersion, "print version info"}, 
    	"user/list":  {"POST", UserList, "print version info"},
        "user/:user": {"GET", DescribeUser, "print user info"}, 
        "user/add":   {"PUT", AssUsers, "add users"},
    }
    
    // prepare gin engine
    router := gin.Default()
    
    // add other flag or more command to kingpin 
    kingpin.HelpFlag.Short('h')
    
    gk := &ginkin.GinKin{
        APIHandlers:     apis,
        ServeGinFunc:    ServeGin,
        CLIFallbackFunc: CLIFallback,
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

func AssUsers(c *gin.Context) {
	var users []string
	if err := c.ShouldBindJSON(&users); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, users)
}

func DescribeVersion(c *gin.Context) {
	c.JSON(http.StatusOK, "0.1")
}

func ServeGin(router *gin.Engine) {
    router.Run(":3000")
}

func CLIFallback(cmd string) {
    log.Println("handled command line action:", cmd)
}
```