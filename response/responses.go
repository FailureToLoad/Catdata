package response

import (
	"context"
	"log/slog"
	"net/http"
)

func BadRequest(ctx context.Context, rw http.ResponseWriter, reason string, err error) {
	http.Error(rw, reason, http.StatusBadRequest)

	slogger := slog.With("reason", reason)
	if err != nil {
		slogger = slog.With(slog.Any("error", err))

	}
	slogger.Error("bad request")
}

func InternalServerError(ctx context.Context, rw http.ResponseWriter, reason string, err error) {
	http.Error(rw, reason, http.StatusInternalServerError)

	slogger := slog.With("reason", reason)
	if err != nil {
		slogger = slog.With(slog.Any("error", err))

	}
	slogger.Error("internal server error")
}
