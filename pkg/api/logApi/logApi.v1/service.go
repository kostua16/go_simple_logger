package logApi_v1

import (
	"github.com/kostua16/go_simple_logger/pkg/db"
	"time"
)

type LogService struct {
	db *db.Connection
}

func (s *LogService) GetLogs() []LogEntry {
	var logs []LogEntry
	s.db.API().Find(&logs)
	return logs
}

func (s *LogService) CleanLogs() {
	tx := s.db.API().Delete(&LogEntry{}, "created_at < ?", time.Now().Add(-1*time.Hour))
	defer tx.Commit()
}

func (s *LogService) AddLog(log LogEntry) {
	tx := s.db.API().Create(&log)
	defer tx.Commit()
}

func NewLogService(db *db.Connection) (*LogService, error) {
	migErr := db.API().AutoMigrate(&LogEntry{})
	if migErr != nil {
		return nil, migErr
	}
	service := &LogService{
		db: db,
	}
	return service, nil
}
