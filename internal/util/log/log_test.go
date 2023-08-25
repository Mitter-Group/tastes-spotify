package log

import (
	"context"
	"testing"
)

func TestInfo(t *testing.T) {
	Info("test", "test2", "test3")
	Infof("test 1:%s 2:%s 3:%s", "test2", "test3", "test4")
	Infow("REQUEST", "method", "GET", "url", "http://localhost:8080/api/v1/cards/123")
	ctx := context.Background()
	ctx = context.WithValue(ctx, "scoringId", 123456)
	ctx = context.WithValue(ctx, "requestId", "abcdef")
	InfofWithContext(ctx, "InfoF test 1:%s 2:%s 3:%s", "test2", "test3", "test4")
	InfoWithContext(ctx, "Info test", "test2", "test3")
}

func TestWarn(t *testing.T) {
	Warn("test", "test2", "test3")
	Warnf("test 1:%s 2:%s 3:%s", "test2", "test3", "test4")
	Warnw("REQUEST", "method", "GET", "url", "http://localhost:8080/api/v1/cards/123")
}

func TestError(t *testing.T) {
	Error("test", "test2", "test3")
	Errorf("test 1:%s 2:%s 3:%s", "test2", "test3", "test4")
	Errorw("REQUEST", "method", "GET", "url", "http://localhost:8080/api/v1/cards/123")
}

func TestDebug(t *testing.T) {
	Debug("test", "test2", "test3")
	Debugf("test 1:%s 2:%s 3:%s", "test2", "test3", "test4")
	Debugw("REQUEST", "method", "GET", "url", "http://localhost:8080/api/v1/cards/123")
}
