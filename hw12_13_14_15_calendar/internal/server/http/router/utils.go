package router

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

func splitPath(s string) (tail, head string) {
	s = path.Clean("/" + s)
	i := strings.Index(s[1:], "/") + 1
	if i <= 0 {
		return "/", s[1:]
	}

	return s[1:i], s[i:]
}

func parseInt(s string) (value int, tail string, err error) {
	s = path.Clean("/" + s)
	i := strings.Index(s[1:], "/") + 1
	if i <= 0 {
		value, err = strconv.Atoi(s[1:])
		return
	}

	value, err = strconv.Atoi(s[1:i])
	tail = s[i:]
	return
}

func respond(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			serverError(w, err)
		}
	}
}

func serverError(w http.ResponseWriter, err error) {
	zap.L().Error("http: internal server error", zap.Error(err))
	respond(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
}

func clientError(w http.ResponseWriter, r *http.Request, code int, err error) {
	zap.L().Warn("http: client error",
		zap.Int("code", code),
		zap.String("method", r.Method),
		zap.String("uri", r.RequestURI),
		zap.String("remote_addr", r.RemoteAddr),
		zap.Error(err),
	)
	respond(w, code, map[string]string{"error": err.Error()})
}
