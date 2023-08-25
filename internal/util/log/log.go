package log

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

type ContextTag string

var client *zap.SugaredLogger

// ctxTagValues is a list of keys that are used to extract values from context and then printed to WithContext log functions
var ctxTagValues = []ContextTag{
	"scoringId",
	"requestId",
}

func init() {
	osValue := os.Getenv("ENV")
	log, _ := zap.NewProduction()
	if osValue == "" || osValue == "local" {
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		log, _ = cfg.Build(zap.AddCaller(), zap.AddCallerSkip(1))
		fmt.Println("Zap started in dev mode")
	}

	client = log.Sugar()
}

func Info(args ...interface{}) {
	client.Info(args...)
}

func Infof(template string, args ...interface{}) {
	client.Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	client.Infow(msg, keysAndValues...)
}

func InfoWithContext(ctx context.Context, args ...interface{}) {
	ctxTags := getTagsFromContext(ctx)
	newArgs := append([]interface{}{ctxTags}, args...)
	client.Info(newArgs...)
}

func InfofWithContext(ctx context.Context, template string, args ...interface{}) {
	ctxTags := getTagsFromContext(ctx)
	client.Infof(ctxTags+template, args...)
}

func Warn(args ...interface{}) {
	client.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	client.Warnf(template, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	client.Warnw(msg, keysAndValues...)
}

func WarnWithContext(ctx context.Context, args ...interface{}) {
	ctxTags := getTagsFromContext(ctx)
	newArgs := append([]interface{}{ctxTags}, args...)
	client.Warn(newArgs...)
}

func WarnfWithContext(ctx context.Context, template string, args ...interface{}) {
	ctxTags := getTagsFromContext(ctx)
	client.Warnf(ctxTags+template, args...)
}

func Error(args ...interface{}) {
	client.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	client.Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	client.Errorw(msg, keysAndValues...)
}

func ErrorWithContext(ctx context.Context, args ...interface{}) {
	ctxTags := getTagsFromContext(ctx)
	newArgs := append([]interface{}{ctxTags}, args...)
	client.Error(newArgs...)
}

func ErrorfWithContext(ctx context.Context, template string, args ...interface{}) {
	ctxTags := getTagsFromContext(ctx)
	client.Errorf(ctxTags+template, args...)
}

func Debug(args ...interface{}) {
	client.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	client.Debugf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	client.Debugw(msg, keysAndValues...)
}

func DebugWithContext(ctx context.Context, args ...interface{}) {
	ctxTags := getTagsFromContext(ctx)
	newArgs := append([]interface{}{ctxTags}, args...)
	client.Debug(newArgs...)
}

func DebugfWithContext(ctx context.Context, template string, args ...interface{}) {
	ctxTags := getTagsFromContext(ctx)
	client.Debugf(ctxTags+template, args...)
}

func Panic(args ...interface{}) {
	client.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	client.Panicf(template, args...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	client.Panicw(msg, keysAndValues...)
}

func PanicWithContext(ctx context.Context, args ...interface{}) {
	ctxTags := getTagsFromContext(ctx)
	newArgs := append([]interface{}{ctxTags}, args...)
	client.Panic(newArgs...)
}

func PanicfWithContext(ctx context.Context, template string, args ...interface{}) {
	ctxTags := getTagsFromContext(ctx)
	client.Panicf(ctxTags+template, args...)
}

// getTagsFromContext extract the values from context that are in the map ctxTagValues and returns a string of them wrapped in parentheses
func getTagsFromContext(ctx context.Context) string {
	var tags []string
	for _, k := range ctxTagValues {
		val := ctx.Value(k)
		if val != nil {
			tags = append(tags, fmt.Sprintf("[%s: %v]", k, ctx.Value(k)))
		}
	}
	return fmt.Sprintf("%s ", strings.Join(tags, ""))
}
