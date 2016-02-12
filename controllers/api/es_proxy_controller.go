package api

import (
	"github.com/gin-gonic/gin"
	// "fmt"
	"net/http/httputil"
	"net/url"
)

type EsProxyController struct {

}
func NewEsProxyController(s *gin.RouterGroup) *EsProxyController {
	ctl := &EsProxyController{}

	s.Any("/es", ctl.reverseProxy)

	return ctl;
}

// type Prox struct {
// 	target *url.URL
// 	proxy *httputil.ReverseProxy
// }

func (pc *EsProxyController) reverseProxy(c *gin.Context) {

	url, _ := url.Parse("http://localhost:9200/collectd-*/")

	// if err != nil {
	// 	fmt.Printf("\n\n%v\n\n", err)
	// } else {
	// 	fmt.Printf("\n\n%v\n\n", url)
	// }

	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.ServeHTTP(c.Writer, c.Request)
}