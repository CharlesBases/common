package middleware

import (
	"bytes"
	"html/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/urfave/negroni"

	seelog "common/log"
)

func Seelog() gin.HandlerFunc {
	defer seelog.Flush()
	return func(c *gin.Context) {
		start := time.Now()
		log0 := negroni.LoggerEntry{
			StartTime: start.Format(LoggerDefaultDateFormat),
			Status:    c.Writer.Status(),
			Duration:  time.Since(start),
			Hostname:  c.Request.Host,
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Request:   c.Request,
		}
		buff := &bytes.Buffer{}
		t := template.Must(template.New("negroni_parser").Parse(LoggerDefaultFormat))
		t.Execute(buff, log0)
		seelog.Info(buff.String())
	}
}

// LoggerDefaultFormat is the format logged used by the default Logger instance.
var LoggerDefaultFormat = "{{.StartTime}} | {{.Status}} | {{.Duration}} | {{.Hostname}} | {{.Method}} {{.Path}}"

// LoggerDefaultDateFormat is the format used for date by the default Logger instance.
var LoggerDefaultDateFormat = "2006-01-02 15:04:05.000"
