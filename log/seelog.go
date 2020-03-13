package log

import (
	"fmt"
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
	logger.SetAdditionalStackDepth(1)
	seelog.ReplaceLogger(logger)
}

func Trace(vs ...interface{}) {
	values := make([]interface{}, 0)
	for _, v := range vs {
		values = append(values, v, " ")
	}
	seelog.Trace(values...)
}

func Tracef(format string, params ...interface{}) {
	seelog.Tracef(format, params...)
}

func Debug(vs ...interface{}) {
	values := make([]interface{}, 0)
	for _, v := range vs {
		values = append(values, v, " ")
	}
	seelog.Debug(values...)
}

func Debugf(format string, params ...interface{}) {
	seelog.Debugf(format, params...)
}

func Info(vs ...interface{}) {
	values := make([]interface{}, 0)
	for _, v := range vs {
		values = append(values, v, " ")
	}
	seelog.Info(values...)
}

func Infof(format string, params ...interface{}) {
	seelog.Infof(format, params...)
}

func Warn(vs ...interface{}) {
	values := make([]interface{}, 0)
	for _, v := range vs {
		values = append(values, v, " ")
	}
	seelog.Warn(values...)
}

func Warnf(format string, params ...interface{}) {
	seelog.Warnf(format, params...)
}

func Error(vs ...interface{}) {
	values := make([]interface{}, 0)
	for _, v := range vs {
		values = append(values, v, " ")
	}
	seelog.Error(values...)
}

func Errorf(format string, params ...interface{}) {
	seelog.Errorf(format, params...)
}

func Critical(vs ...interface{}) {
	values := make([]interface{}, 0)
	for _, v := range vs {
		values = append(values, v, " ")
	}
	values = append(values, "\n", string(debug.Stack()))
	seelog.Critical(values...)
}

func Criticalf(format string, params ...interface{}) {
	seelog.Criticalf(fmt.Sprintf("%s\n%s", format, string(debug.Stack())), params...)
}

func Flush() {
	seelog.Flush()
}
