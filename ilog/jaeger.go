package ilog

type jaegerLog loggingT

// 只能通过ILog派生过来
func JaegerLog() jaegerLog {
	return (jaegerLog)(logging)
}

//
func (l jaegerLog) Error(msg string) {
	logging.print(errorLog, msg)
}

func (l jaegerLog) Infof(format string, args ...interface{}) {
	logging.printf(infoLog, format, args...)
}
