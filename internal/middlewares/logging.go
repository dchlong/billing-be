package middlewares

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"

	"github.com/dchlong/billing-be/pkg/logger"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func SetupLog(ilogger logger.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ctx = ilogger.NewContext(ctx)
		log := ilogger.GetLogger(ctx)
		c.Request = c.Request.WithContext(ctx)
		path := c.Request.URL.Path
		buf, _ := io.ReadAll(c.Request.Body)
		requestBody := string(buf)
		readCloser := io.NopCloser(bytes.NewBuffer(buf))
		// We have to create a new Buffer and transfer it to request body again.
		c.Request.Body = readCloser
		log.Infof("Request to %s: payload: %s", path, requestBody)

		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}

		c.Writer = blw
		c.Next()
		log.Infof("Response: %s", blw.body.String())
	}
}
