package util

import (
	"github.com/elazarl/goproxy"
	"net/http"
)

// RunProxy ..
func RunProxy(ch chan error) {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	ch <- http.ListenAndServe(":8080", proxy)
}
