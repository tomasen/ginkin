# ginkin

Build command line interface and http api server at the same time by mix Gin and Kingpin.

## Example

```go
func main() {
    apis := ginkin.APIHandlers{
      "version": {"GET", DescribeVersion, "print version ifno"},
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

func ServeGin(router *gin.Engine) {
    router.Run(":3000")
}

func CLIFallback(cmd string) {
    log.Println("handled command line action:", cmd)
}
```