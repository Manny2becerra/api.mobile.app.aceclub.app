package Logging

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	constants "api-mobile-app/src/api/common/constants"
	errorHandling "api-mobile-app/src/api/common/errorHandling"

	"github.com/google/uuid"
)

type requestInformation struct {
	timeStamp string

	level string

	requestId string

	httpMethod string

	path string

	query *string

	statusCode int

	duration_ms int64

	bytesIn int64

	clientIp string

	userId *string

	error *errorHandling.ProblemDetails
}
type Logger struct {
	Log *slog.Logger
}

func NewLogger() *Logger {

	return &Logger{
		Log: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

}

func (l *Logger) LogRequest(r *http.Request, errorValue *errorHandling.ProblemDetails) {
	var requestInfo requestInformation

	if errorValue == nil {
		// No error case
		requestInfo = requestInformation{
			timeStamp:   time.Now().GoString(),
			level:       "info",
			requestId:   uuid.New().String(),
			httpMethod:  r.Method,
			path:        r.URL.Path,
			statusCode:  200,
			duration_ms: time.Since(r.Context().Value(constants.StartTimeKey).(time.Time)).Milliseconds(),
			bytesIn:     r.ContentLength,
			clientIp:    r.RemoteAddr,
		}
	} else {
		// Error case
		problemDetails := errorValue
		requestInfo = requestInformation{
			timeStamp:   time.Now().GoString(),
			level:       "info",
			requestId:   uuid.New().String(),
			httpMethod:  r.Method,
			path:        r.URL.Path,
			statusCode:  problemDetails.Status,
			duration_ms: time.Since(r.Context().Value(constants.StartTimeKey).(time.Time)).Milliseconds(),
			bytesIn:     r.ContentLength,
			clientIp:    r.RemoteAddr,
			error:       problemDetails,
		}
	}

	attrs := []slog.Attr{
		slog.String("timeStamp", requestInfo.timeStamp),
		slog.String("request_id", requestInfo.requestId),
		slog.String("method", requestInfo.httpMethod),
		slog.String("path", requestInfo.path),
		slog.String("status_code", strconv.Itoa(requestInfo.statusCode)),
		slog.String("duration_ms", strconv.FormatInt(requestInfo.duration_ms, 10)),
		slog.String("bytes_in", strconv.FormatInt(requestInfo.bytesIn, 10)),
		slog.String("clientIp", requestInfo.clientIp),
	}

	if requestInfo.userId != nil {
		attrs = append(attrs, slog.String("user_Id", *requestInfo.userId))
	}

	if requestInfo.query != nil {
		attrs = append(attrs, slog.String("query", *requestInfo.query))
	}

	if requestInfo.error != nil {
		attrs = append(attrs,
			slog.String("error_type", requestInfo.error.ErrorType),
			slog.String("error_title", requestInfo.error.Title),
			slog.String("error_detail", requestInfo.error.Detail),
		)

		if requestInfo.error.Instance != nil {
			attrs = append(attrs, slog.String("error_instance", *requestInfo.error.Instance))
		}
	}

	l.Log.LogAttrs(r.Context(), slog.LevelInfo, "http_request", attrs...)
}

// the type http.Handler is an interface which specifies any type that has function ServeHttp(w, r)
func (l *Logger) Logging(next http.Handler) http.Handler {

	// http.HandlerFunc is a function that returns a handler. it allows us to turn any function into an http handler

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		ctx := context.WithValue(r.Context(), constants.StartTimeKey, start)

		r = r.WithContext(ctx)

		// were calling the original handler's ServeHTTP and calling its function passing in the w, r from the handler we created
		next.ServeHTTP(w, r)

	})
}

func (l *Logger) StartUpAPI() {
	var logo string = "\033[38;2;1;223;230m" + `
//              _____  ______    _____  _      _    _  ____  ` + "\033[38;2;1;200;207m" + `
//       /\    / ____||  ____|  / ____|| |    | |  | ||  _ \ ` + "\033[38;2;1;177;184m" + `
//      /  \  | |     | |__    | |     | |    | |  | || |_) |` + "\033[38;2;1;154;161m" + `
//     / /\ \ | |     |  __|   | |     | |    | |  | ||  _ < ` + "\033[38;2;1;131;138m" + `
//    / ____ \| |____ | |____  | |____ | |____| |__| || |_) |` + "\033[38;2;1;108;115m" + `
//   /_/    \_\\_____||______|  \_____||______|\____/ |____/ ` + "\033[38;2;1;85;92m" + `
//                                                           ` + "\033[38;2;1;62;69m" + `
//                                                           ` + "\033[0m"

	fmt.Println(logo)
	l.Log.Info("API Started Successfully")
}
