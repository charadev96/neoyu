package middleware

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"connectrpc.com/connect"
	"github.com/rs/zerolog"
)

func NewConsoleLogger() zerolog.Logger {
	writer := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "15:04",
	}

	writer.FormatFieldName = func(i any) string {
		return fmt.Sprintf("\n  \033[36m%s: \033[0m", i)
	}

	writer.FormatErrFieldName = func(i any) string {
		return fmt.Sprintf("\n  \033[31m%s: \033[0m", i)
	}

	logger := zerolog.New(writer).
		With().
		Timestamp().
		Logger()

	return logger
}

func NewLogInterceptor(log *zerolog.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			proc := strings.TrimPrefix(req.Spec().Procedure, "/")
			parts := strings.SplitN(proc, "/", 2)

			var service, method string
			if len(parts) == 2 {
				service, method = parts[0], parts[1]
			}

			res, err := next(ctx, req)
			if err != nil {
				var errCon *connect.Error
				if errors.As(err, &errCon) {
					log.Error().
						Str("service", service).
						Str("method", method).
						Err(errCon).
						Msg("request")
					return res, err
				}
			}

			log.Info().
				Str("service", service).
				Str("method", method).
				Msg("request")

			return res, err
		}
	}
}
