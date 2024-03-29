package middleware

import (
	"bytes"
	"fmt"
	"github.com/shammishailaj/taskfactory/pkg/utils"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type Log struct {
	l         *log.Logger
	printBody bool
}

func NewLog(l *log.Logger, printBody bool) *Log {
	return &Log{
		l:         l,
		printBody: printBody,
	}
}

func (l *Log) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := utils.NewUtils(l.l)
		w.Header().Set("X-Backend-Server", u.Hostname())
		entry := fmt.Sprintf("%s %s", r.Method, r.URL.RequestURI())
		if l.printBody {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body.Close()

			r.Body = ioutil.NopCloser(bytes.NewReader(b))

			entry = fmt.Sprintf("%s %s", entry, string(b))
		}

		l.l.Info(entry)

		next.ServeHTTP(w, r)
	})
}
