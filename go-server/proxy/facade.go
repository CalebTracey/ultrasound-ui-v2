package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

const clientURL = "http://localhost:6088"
const serverURL = "http://localhost:6080"

type ProxyFacade interface {
	//Client() *httputil.ReverseProxy
	Server() *httputil.ReverseProxy
}

func NewService() Service {
	return Service{}
}

type Service struct{}

//func (s Service) Client() *httputil.ReverseProxy {
//	remote, err := url.Parse(clientURL)
//	if err != nil {
//		panic(err)
//	}
//	proxy, proxyErr := newProxy(clientURL)
//	if proxyErr != nil {
//		panic(proxyErr)
//	}
//
//	return proxy
//}

func (s Service) Server() *httputil.ReverseProxy {
	origin, _ := url.Parse(clientURL)

	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = "http"
		req.URL.Host = origin.Host
	}
	proxy := &httputil.ReverseProxy{Director: director}

	//if proxyErr != nil {
	//	panic(proxyErr)
	//}
	return proxy
}

//
//func newProxy(targetHost string) (*httputil.ReverseProxy, error) {
//	hostUrl, err := url.Parse(targetHost)
//	logrus.Infoln("URL: " + hostUrl.RawPath)
//	if err != nil {
//		return nil, err
//	}
//
//	return httputil.NewSingleHostReverseProxy(hostUrl), nil
//}
