package logger

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func LoggerInit() {
	// 日志作为JSON而不是默认的ASCII格式器.
	log.SetFormatter(&log.JSONFormatter{})

	// 输出到标准输出,可以是任何io.Writer
	log.SetOutput(os.Stdout)

	// 只记录xx级别或以上的日志
	log.SetLevel(log.TraceLevel)
}
