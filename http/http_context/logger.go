package http_context

import (
	"fmt"

	"github.com/cihub/seelog"

	"github.com/CharlesBases/common/algo"
)

type logger struct {
	trace string
}

const defaultSeelogConfig = "./config/seelog.xml"

func init() {
	logger, err := seelog.LoggerFromConfigAsFile(defaultSeelogConfig)
	if err != nil {
		seelog.Warn("seelog config error: ", err, " using defualt seelog config.")
		logger, _ = seelog.LoggerFromConfigAsString(`
			<?xml version="1.0" encoding="utf-8" ?>
			<seelog levels="debug,info,warn,error,critical">
				<outputs formatid="main">
					<filter levels="warn">
						<console formatid="main"/>
					</filter>
					<filter levels="info">
						<console formatid="info"/>
					</filter>
					<filter levels="debug">
						<console formatid="debug"/>
					</filter>
					<filter levels="error,critical">
						<console formatid="error"/>
					</filter>
					<rollingfile formatid="main" type="date" filename="./log/seelog.log" datepattern="2006-01-02" maxrolls="30" namemode="prefix"/>
				</outputs>
				<formats>
					<format id="main"  format="[%Date(2006-01-02 15:04:05.000)][%LEV] %Func %File:%Line ==&gt; %Msg%n"/>
					<format id="info"  format="%EscM(32)[%Date(2006-01-02 15:04:05.000)][%LEV] %Func %File:%Line ==&gt; %Msg%n%EscM(0)"/>
					<format id="debug" format="%EscM(36)[%Date(2006-01-02 15:04:05.000)][%LEV] %Func %File:%Line ==&gt; %Msg%n%EscM(0)"/>
					<format id="error" format="%EscM(31)[%Date(2006-01-02 15:04:05.000)][%LEV] %Func %File:%Line ==&gt; %Msg%n%EscM(0)"/>
				</formats>
			</seelog>`)
	} else {
		seelog.Infof("using seelog configed from %s", defaultSeelogConfig)
	}
	logger.SetAdditionalStackDepth(2)
	seelog.ReplaceLogger(logger)
}

func newlogger() *logger {
	return &logger{trace: fmt.Sprintf("[TarceID: %s] ", algo.GetTraceID())}
}

func (httpContext *httpContext) Debug(vs ...interface{}) {
	seelog.Debug(append([]interface{}{httpContext.trace}, vs...)...)
}

func (httpContext *httpContext) Debugf(format string, vs ...interface{}) {
	seelog.Debugf(fmt.Sprintf("%s%s", httpContext.trace, format), vs...)
}

func (httpContext *httpContext) Info(vs ...interface{}) {
	seelog.Info(append([]interface{}{httpContext.trace}, vs...)...)
}

func (httpContext *httpContext) Infof(format string, vs ...interface{}) {
	seelog.Infof(fmt.Sprintf("%s%s", httpContext.trace, format), vs...)
}

func (httpContext *httpContext) Warn(vs ...interface{}) {
	seelog.Warn(append([]interface{}{httpContext.trace}, vs...)...)
}

func (httpContext *httpContext) Warnf(format string, vs ...interface{}) {
	seelog.Warnf(fmt.Sprintf("%s%s", httpContext.trace, format), vs...)
}

func (httpContext *httpContext) Error(vs ...interface{}) {
	seelog.Error(append([]interface{}{httpContext.trace}, vs...)...)
}

func (httpContext *httpContext) Errorf(format string, vs ...interface{}) {
	seelog.Errorf(fmt.Sprintf("%s%s", httpContext.trace, format), vs...)
}

func (httpContext *httpContext) Flush() {
	seelog.Flush()
}
