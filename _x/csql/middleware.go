package csql

import (
	"bufio"
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/tusharsoni/copper/chttp"

	"github.com/tusharsoni/copper/clogger"
	"gorm.io/gorm"
)

// DBTxnMiddleware provides a middleware that wraps the http request in a database transaction. If the response status
// code is between 100-399, the transaction is committed, else a rollback is performed.
type DBTxnMiddleware interface {
	WrapInTxn(next http.Handler) http.Handler
}

type dbTxnMiddleware struct {
	db     *gorm.DB
	logger clogger.Logger
}

func NewDBTxnMiddleware(db *gorm.DB, logger clogger.Logger) chttp.MiddlewareFunc {
	mw := dbTxnMiddleware{
		db:     db,
		logger: logger,
	}

	return mw.WrapInTxn
}

func NewDBTxnMiddlewareFx(db *gorm.DB, logger clogger.Logger) chttp.GlobalMiddlewareFuncResult {
	return chttp.GlobalMiddlewareFuncResult{
		GlobalMiddlewareFunc: NewDBTxnMiddleware(db, logger),
	}
}

func (m *dbTxnMiddleware) WrapInTxn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		txn := m.db.Begin()
		ctx := context.WithValue(r.Context(), connCtxKey, txn)

		rw := &txnrw{
			internal: w,
			db:       txn,
			logger:   m.logger,
		}

		next.ServeHTTP(rw, r.WithContext(ctx))

		err := rw.commitIfNeeded()
		if err != nil {
			m.logger.Error("Failed to commit db transaction", err)
			return
		}
	})
}

type txnrw struct {
	internal http.ResponseWriter
	db       *gorm.DB
	logger   clogger.Logger

	didCommit bool
}

func (w *txnrw) Header() http.Header {
	return w.internal.Header()
}

func (w *txnrw) Write(b []byte) (int, error) {
	return w.internal.Write(b)
}

func (w *txnrw) WriteHeader(statusCode int) {
	if statusCode >= 400 {
		w.didCommit = true
		w.db.Rollback()
		w.internal.WriteHeader(statusCode)
		return
	}

	err := w.commitIfNeeded()
	if err != nil {
		w.logger.Error("Failed to commit db transaction", err)
		w.internal.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.internal.WriteHeader(statusCode)
}

func (w *txnrw) commitIfNeeded() error {
	if w.didCommit {
		return nil
	}

	w.didCommit = true
	return w.db.Commit().Error
}

func (w *txnrw) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.internal.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("internal response writer is not http.Hijacker")
	}

	return h.Hijack()
}
