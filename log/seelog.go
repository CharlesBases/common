package log

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/cihub/seelog"
)

const defaultSeelogConfig = "./config/seelog.xml"

func init() {
	logger, err := seelog.LoggerFromConfigAsFile(defaultSeelogConfig)
	if err != nil {
		seelog.Warnf("load seelog config error: %v, using defualt seelog config.", err)
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
	logger.SetAdditionalStackDepth(1)
	seelog.UseLogger(logger)

	flush()
}

func Trace(vs ...interface{}) {
	seelog.Trace(vs...)
}

func Tracef(format string, params ...interface{}) {
	seelog.Tracef(format, params...)
}

func Debug(vs ...interface{}) {
	seelog.Debug(vs...)
}

func Debugf(format string, params ...interface{}) {
	seelog.Debugf(format, params...)
}

func Info(vs ...interface{}) {
	seelog.Info(vs...)
}

func Infof(format string, params ...interface{}) {
	seelog.Infof(format, params...)
}

func Warn(vs ...interface{}) {
	seelog.Warn(vs...)
}

func Warnf(format string, params ...interface{}) {
	seelog.Warnf(format, params...)
}

func Error(vs ...interface{}) {
	seelog.Error(vs...)
}

func Errorf(format string, params ...interface{}) {
	seelog.Errorf(format, params...)
}

func Fatal(vs ...interface{}) {
	seelog.Critical(append(vs, "\n", string(debug.Stack()))...)
}

func Fatalf(format string, params ...interface{}) {
	seelog.Criticalf(fmt.Sprintf("%s\n%s", format, string(debug.Stack())), params...)
}

func Critical(vs ...interface{}) {
	seelog.Critical(append(vs, "\n", string(debug.Stack()))...)
}

func Criticalf(format string, params ...interface{}) {
	seelog.Criticalf(fmt.Sprintf("%s\n%s", format, string(debug.Stack())), params...)
}

func flush() {
	go func() {
		s := make(chan os.Signal)
		signal.Notify(s, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGSTOP, syscall.SIGKILL, syscall.SIGTERM)
		<-s
		seelog.Flush()
	}()
}
