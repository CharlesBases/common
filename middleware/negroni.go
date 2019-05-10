package middleware

import (
	"bytes"
	"html/template"
	"net/http"
	"time"

	seelog "common/log"

	"github.com/urfave/negroni"
)

var (
	DefaultDateFormat = "2006-01-02 15:04:05.000"
	DefaultFormat     = "{{.StartTime}} | {{.Status}} | {{.Duration}} | {{.Hostname}} | {{.Method}} {{.Path}}"
)

type logger struct {
	dateFormat string
	template   *template.Template
}

func NegroniLogger() *logger {
	return &logger{
		dateFormat: DefaultDateFormat,
		template:   template.Must(template.New("negroni_parser").Parse(DefaultFormat)),
	}
}

func (l *logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer seelog.Flush()
	start := time.Now()
	next(rw, r)
	res := rw.(negroni.ResponseWriter)
	log := negroni.LoggerEntry{
		StartTime: start.Format(l.dateFormat),
		Status:    res.Status(),
		Duration:  time.Since(start),
		Hostname:  r.Host,
		Method:    r.Method,
		Path:      r.URL.Path,
		Request:   r,
	}
	buff := &bytes.Buffer{}
	l.template.Execute(buff, log)
	seelog.Info(buff.String())
}
