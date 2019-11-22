package weather

import (
	"net/http"
	"time"

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

func (al *WebApi) AccessLog(ha http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := NewLoggingResponseWriter(w)
		st := time.Now()
		ha.ServeHTTP(lrw, r)
		rt := time.Now().Sub(st)
		glog.Infof("Timestamp:%d Method:%s URI:%s StatusCode:%d ResponseTime:%d RemoteAddr:%s", st, r.Method, r.RequestURI, lrw.statusCode, rt/time.Millisecond, r.RemoteAddr)
	})
}
