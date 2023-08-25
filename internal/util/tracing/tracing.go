package tracing

import (
	"context"
	"fmt"
	"net/http"

	"github.com/erkanzileli/nrfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type ContextTag string

func CreateContextWithTransaction(c *fiber.Ctx) context.Context {
	txn := nrfiber.FromContext(c)
	if txn != nil {
		txn.SetName(fmt.Sprintf("%s %s", c.Request().Header.Method(), c.Route().Path))
	}
	return context.WithValue(c.Context(), ContextTag("nrTransaction"), txn)
}

func GetTransactionFromContext(c context.Context) *newrelic.Transaction {
	res, ok := c.Value(ContextTag("nrTransaction")).(*newrelic.Transaction)
	if ok && res != nil {
		return res
	}
	return newrelic.FromContext(c)
}

func RecordLog(c context.Context, severity string, msg string) {
	txn := GetTransactionFromContext(c)
	if txn == nil {
		txn = newrelic.FromContext(c)
	}
	if txn != nil {
		txn.RecordLog(newrelic.LogData{
			Severity: severity,
			Message:  msg,
		})
	}
}

func StartSegment(c context.Context, name string) *newrelic.Segment {
	txn := GetTransactionFromContext(c)
	if txn != nil {
		return txn.StartSegment(name)
	}
	return nil
}

func StartExternalSegmentFromRequest(c context.Context, r *http.Request) *newrelic.ExternalSegment {
	txn := GetTransactionFromContext(c)
	if txn != nil {
		return newrelic.StartExternalSegment(txn, r)
	}
	return nil
}
