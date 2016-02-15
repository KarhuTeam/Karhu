package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gotoolz/env"
	"net/http/httputil"
	"net/url"
	"strings"
)

type EsProxyController struct {
	url   *url.URL
	proxy *httputil.ReverseProxy
}

func NewEsProxyController(s *gin.RouterGroup) *EsProxyController {
	ctl := &EsProxyController{}

	var err error
	ctl.url, err = url.Parse(env.GetDefault("ES_HOST", "http://localhost:9200/"))
	if err != nil {
		panic(err)
	}
	ctl.proxy = httputil.NewSingleHostReverseProxy(ctl.url)

	s.Any("/es/*param", ctl.reverseProxy)

	return ctl
}

func (espc *EsProxyController) reverseProxy(c *gin.Context) {

	c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/api/es")

	espc.proxy.ServeHTTP(c.Writer, c.Request)
}
