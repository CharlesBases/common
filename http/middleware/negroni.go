package middleware

import (
	"bytes"
	"net/http"
	"text/template"
	"time"

	"github.com/urfave/negroni"

	"github.com/CharlesBases/common/log"
)

var (
	defaultDateFormat = "2006-01-02 15:04:05.000"
	defaultFormat     = "{{.StartTime}} | {{.Status}} | {{.Duration}} | {{.Hostname}} | {{.Method}} {{.Path}}"
	defaulttemplate   = template.Must(template.New("negroni_parser").Parse(defaultFormat))
)

type negroniLogger struct{}

func Negroni() *negroniLogger {
	return new(negroniLogger)
}

func (nl *negroniLogger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer log.Flush()

	start := time.Now()

	next(rw, r)

	writer := rw.(negroni.ResponseWriter)
	logger := negroni.LoggerEntry{
		StartTime: start.Format(defaultDateFormat),
		Status:    writer.Status(),
		Duration:  time.Since(start),
		Hostname:  r.Host,
		Method:    r.Method,
		Path:      r.URL.Path,
		Request:   r,
	}
	buff := new(bytes.Buffer)
	defaulttemplate.Execute(buff, logger)
	log.Info(buff.String())
}
