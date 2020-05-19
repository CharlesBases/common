package middleware

import (
	"bytes"
	"net/http"
	"text/template"
	"time"

	"github.com/urfave/negroni"

	"github.com/CharlesBases/common/http/http_context"
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
	start := time.Now()
	httpContext := http_context.NewContext()
	defer httpContext.Flush()
	next(rw, r.WithContext(httpContext))
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
	httpContext.Info(buff.String())
}
