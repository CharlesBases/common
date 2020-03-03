package log

import (
	"runtime/debug"

	"github.com/cihub/seelog"
)

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
					<format id="main" format="[%Date(2006-01-02 15:04:05.000)][%LEV] %Func %File:%Line ==&gt; %Msg%n"/>
					<format id="info" format="%EscM(32)[%Date(2006-01-02 15:04:05.000)][%LEV] %Func %File:%Line ==&gt; %Msg%n%EscM(0)"/>
					<format id="debug" format="%EscM(36)[%Date(2006-01-02 15:04:05.000)][%LEV] %Func %File:%Line ==&gt; %Msg%n%EscM(0)"/>
					<format id="error" format="%EscM(31)[%Date(2006-01-02 15:04:05.000)][%LEV] %Func %File:%Line ==&gt; %Msg%n%EscM(0)"/>
				</formats>
			</seelog>`)
	} else {
		seelog.Infof("using seelog configed from %s", defaultSeelogConfig)
	}
	logger.SetAdditionalStackDepth(0)
	seelog.ReplaceLogger(logger)
}

func Trace(v ...interface{}) {
	seelog.Trace(v...)
}

func Tracef(format string, params ...interface{}) {
	seelog.Tracef(format, params...)
}

func Debug(v ...interface{}) {
	seelog.Debug(v...)
}

func Debugf(format string, params ...interface{}) {
	seelog.Debugf(format, params...)
}

func Info(v ...interface{}) {
	seelog.Info(v...)
}

func Infof(format string, params ...interface{}) {
	seelog.Infof(format, params...)
}

func Warn(v ...interface{}) {
	seelog.Warn(v...)
}

func Warnf(format string, params ...interface{}) {
	seelog.Warnf(format, params...)
}

func Error(v ...interface{}) {
	seelog.Error(v...)
}

func Errorf(format string, params ...interface{}) {
	seelog.Errorf(format, params...)
}

func Fatal(v ...interface{}) {
	seelog.Critical(v...)
}

func Fatalf(format string, params ...interface{}) {
	seelog.Criticalf(format, params...)
}

func Critical(v ...interface{}) {
	v = append(v, string(debug.Stack()))
	seelog.Critical(v...)
}

func Criticalf(format string, params ...interface{}) {
	format += string(debug.Stack())
	seelog.Criticalf(format, params...)
}

func Flush() {
	seelog.Flush()
}
