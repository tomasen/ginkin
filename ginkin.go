package ginkin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
)

type APIHandler struct {
	HTTPMethod  string
	Handler     gin.HandlerFunc
	Help        string
}

type APIHandlers map[string]APIHandler
type ServeGinFunc func(engine *gin.Engine)
type CLIFallbackFunc func(cmd string)

type GinKin struct {
	APIHandlers
	ServeGinFunc
	CLIFallbackFunc
}

var UnderCommandLine bool

// Run Gin Server or process command line
// relativePath should end with "/"
func (gk *GinKin) Run(router *gin.Engine, relativePath string, middleware ...gin.HandlerFunc)  {
	kingpin.Command("start/server", "").Default()

	v1api := router.Group(relativePath)
	v1api.Use(middleware...)
	payloads := map[string]*string{}
	for cmd, v := range gk.APIHandlers {
		v1api.Handle(v.HTTPMethod, relativePath + cmd, v.Handler)
		c := kingpin.Command(cmd, v.Help)
		payloads[cmd] = c.Arg("payload", "").Default("").String()
	}

	cmd  := kingpin.Parse()
	switch cmd {
	case "start/server":
		// start gin server
		gk.ServeGinFunc(router)
	default:
		// process command line
		UnderCommandLine = true

		handler, exist := gk.APIHandlers[cmd]
		if !exist {
			if gk.CLIFallbackFunc != nil {
				gk.CLIFallbackFunc(cmd)
			} else {
				log.Println("unhandled the command line action:", cmd)
			}
			break
		}

		router.NoMethod()
		payload := payloads[cmd]
		w := httptest.NewRecorder()
		var buf io.Reader = nil
		if len(*payload) > 0 {
			re := regexp.MustCompile(`\:[a-zA-Z0-9]+`)
			if re.MatchString(cmd) {
				// TODO: by this way we can only support one param
				// might need support more
				cmd = re.ReplaceAllString(cmd, *payload)
			} else {
				buf = bytes.NewReader([]byte(*payload))
			}
		}
		req, err := http.NewRequest(handler.HTTPMethod, relativePath + cmd, buf)
		if err != nil {
			log.Fatalln("fail to create request", err)
		}
		router.ServeHTTP(w, req)

		// pretty print response
		dst := &bytes.Buffer{}
		if err := json.Indent(dst, w.Body.Bytes(), "", "  "); err != nil {
			log.Println("json.Indent error:", err)
			fmt.Println(w.Body.String())
			return
		}

		fmt.Println(dst.String())
	}
}