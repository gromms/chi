package middleware

// Ported from Goji's middleware, source:
// https://github.com/zenazn/goji/tree/master/web/middleware

import (
	"bytes"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/pressly/chi"
	"golang.org/x/net/context"
)

// Recoverer is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible.
//
// Recoverer prints a request ID if one is provided.
func Recoverer(next chi.Handler) chi.Handler {
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		reqID := GetReqID(ctx)

		defer func() {
			if err := recover(); err != nil {
				printPanic(reqID, err)
				debug.PrintStack()
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTPC(ctx, w, r)
	}

	return chi.HandlerFunc(fn)
}

func printPanic(reqID string, err interface{}) {
	var buf bytes.Buffer

	if reqID != "" {
		cW(&buf, bBlack, "[%s] ", reqID)
	}
	cW(&buf, bRed, "panic: %+v", err)

	log.Print(buf.String())
}
