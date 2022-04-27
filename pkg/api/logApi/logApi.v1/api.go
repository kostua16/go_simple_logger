package logApi_v1

import (
	"github.com/go-chi/chi"
	"github.com/kostua16/go_simple_logger/pkg/db"
	"github.com/kostua16/go_simple_logger/pkg/logger"
	"github.com/kostua16/go_simple_logger/pkg/webUtils"
	"net/http"
)

var log = logger.CreateLogger("logApi")

func CreateApi(connection *db.Connection) (func(r chi.Router), error) {
	service, err := NewLogService(connection)

	if err != nil {
		return nil, err
	}

	router := func(r chi.Router) {
		r.Get("/", func(resp http.ResponseWriter, request *http.Request) {
			writeErr := webUtils.WriteJson(resp, EntriesToResponses(service.GetLogs()))
			if writeErr != nil {
				log.Errorf("/api/logApi.v1/logs - failed to write response: %v", writeErr)
				resp.WriteHeader(http.StatusInternalServerError)
				_, _ = resp.Write([]byte(writeErr.Error()))
			} else {
				resp.WriteHeader(http.StatusOK)
			}
		})
		r.Put("/", func(resp http.ResponseWriter, request *http.Request) {
			logRequest := &LogRequest{}
			readErr := webUtils.ReadJson(request, logRequest)
			if readErr != nil {
				log.Errorf("/api/logApi.v1/logs - failed to read new log request: %v", readErr)
				resp.WriteHeader(http.StatusInternalServerError)
				_, _ = resp.Write([]byte(readErr.Error()))
			} else {
				service.AddLog(logRequest.toEntry())
				resp.WriteHeader(http.StatusOK)
			}
		})
		r.Delete("/", func(resp http.ResponseWriter, request *http.Request) {
			service.CleanLogs()
			resp.WriteHeader(http.StatusOK)
		})
	}

	return router, nil
}
