package api

import (
	"github.com/gin-gonic/gin"
	"net/http/httputil"
	"net/url"
	"strings"
	"fmt"
)

type EsProxyController struct {
}

func NewEsProxyController(s *gin.RouterGroup) *EsProxyController {
	ctl := &EsProxyController{}

	s.Any("/es/*param", ctl.reverseProxy)

	return ctl
}

// type Prox struct {
// 	target *url.URL
// 	proxy *httputil.ReverseProxy
// }

func (pc *EsProxyController) reverseProxy(c *gin.Context) {

	url, _ := url.Parse("http://localhost:9200/")

	// if err != nil {
	// 	fmt.Printf("\n\n%v\n\n", err)
	// } else {
	// 	fmt.Printf("\n\n%v\n\n", url)
	// }

	proxy := httputil.NewSingleHostReverseProxy(url)

	c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/api/es/")

	fmt.Printf("URL > %v\n", c.Request.URL)
	
	proxy.ServeHTTP(c.Writer, c.Request)
}
