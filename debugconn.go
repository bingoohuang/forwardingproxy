package fproxy

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"net"
)

func netDirection(dest, src net.Conn) string {
	return fmt.Sprintf("%s-%s-%s-%s", src.RemoteAddr(), src.LocalAddr(), dest.LocalAddr(), dest.RemoteAddr())
}

func newDebugReadCloser(r io.ReadCloser, direction string, logger *zap.Logger) io.ReadCloser {
	return &debugReadCloser{
		ReadCloser: r,
		Direction:  direction,
		Logger:     logger,
	}
}

type debugReadCloser struct {
	io.ReadCloser
	Direction string
	*zap.Logger
	n int
}

func (r *debugReadCloser) Read(p []byte) (n int, err error) {
	n, err = r.ReadCloser.Read(p)
	if n > 0 {
		r.n += n
		r.Logger.Debug("transferred",
			zap.String("direction", r.Direction),
			zap.Int("bytes", n),
			zap.Int("total", r.n))
	}

	if err != nil {
		r.Logger.Debug("read error",
			zap.String("direction", r.Direction),
			zap.Int("total", r.n),
			zap.String("error", err.Error()))
	}

	return n, err
}
