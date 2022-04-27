package logApi_v1

import (
	"gorm.io/gorm"
	"time"
)

type LogRequest struct {
	Message string `json:"message" xml:"message"`
	Service string `json:"service" xml:"service"`
	Type    string `json:"type" xml:"type"`
}

type LogEntry struct {
	gorm.Model
	Message string
	Service string
	Type    string
}

type LogResponse struct {
	Message string    `json:"message" xml:"message"`
	Service string    `json:"service" xml:"service"`
	Type    string    `json:"type" xml:"type"`
	Created time.Time `json:"created" xml:"created"`
}

func (r *LogRequest) toEntry() LogEntry {
	return LogEntry{
		Message: r.Message,
		Service: r.Service,
		Type:    r.Type,
		//Created: time.Now().UTC(),
		//Id:      time.Now().UnixNano(),
	}
}

func (r *LogEntry) toResponse() LogResponse {
	return LogResponse{
		Message: r.Message,
		Service: r.Service,
		Type:    r.Type,
		Created: r.CreatedAt.UTC(),
	}
}

func EntriesToResponses(requests []LogEntry) []LogResponse {
	var result = make([]LogResponse, len(requests))
	for idx, i := range requests {
		result[idx] = i.toResponse()
	}
	return result
}
