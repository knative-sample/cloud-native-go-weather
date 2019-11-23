package utils

import (
	"net/http"
	"time"

	"strings"

	"github.com/golang/glog"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func AccessLog(ha http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := NewLoggingResponseWriter(w)
		st := time.Now()
		ha.ServeHTTP(lrw, r)
		rt := time.Now().Sub(st)
		module := ""
		if strings.Contains(r.RequestURI, "/api/area/weather") {
			module = "DETAIL"
		}
		if strings.Contains(r.RequestURI, "/api/area/list") || strings.Contains(r.RequestURI, "/api/cities") {
			module = "CITY"
		}
		if strings.Contains(r.RequestURI, "/api/city/list") || strings.Contains(r.RequestURI, "/api/city/weather") {
			module = "WEB"
		}
		glog.Infof("[%s] Timestamp:%d Method:%s URI:%s StatusCode:%d ResponseTime:%d RemoteAddr:%s", module, st, r.Method, r.RequestURI, lrw.statusCode, rt/time.Millisecond, r.RemoteAddr)
	})
}
