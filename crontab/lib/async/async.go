package async

import log "github.com/sirupsen/logrus"

// GO 全局声明的异步 - 执行异步
func GO(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.WithFields(log.Fields{
					"GO": "err",
				}).Error(err)
			}
		}()

		f()
	}()
}
