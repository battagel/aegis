package logger

type Logger interface {
	Errorln(string)
	Infoln(string)
	Debugln(string)
	Warnln(string)
	Errorw(string, ...interface{})
	Infow(string, ...interface{})
	Debugw(string, ...interface{})
	Warnw(string, ...interface{})
	Sync() error
}
