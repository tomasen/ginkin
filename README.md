# ginkin

Build command line interface and http api server at the same time by mix Gin and Kingpin.

## Example

By following example. 

- `go run main.go` will start http api server handle requests to "http://:3000/version".
- `go run main.go version` will print out version info directly.

```go
func main() {
    apis := ginkin.APIHandlers{
      "version": {"GET", DescribeVersion, "print version info"},
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