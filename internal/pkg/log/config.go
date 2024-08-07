package log

import "gopkg.in/natefinch/lumberjack.v2"

var lumlog = &lumberjack.Logger{
	// Filename:   "app.log",
	MaxSize:    10,  // megabytes
	MaxBackups: 3,   // number of log files
	MaxAge:     365, // days
	Compress:   true,
}
