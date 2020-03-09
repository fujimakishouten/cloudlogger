package repository

import (
	"errors"
	"regexp"
	"time"

	"cloudlogger/entity"
)

type LogRepository struct {
	RegExp *regexp.Regexp
	Format string
	Logs []*entity.LogEntity
}

// Set RegExp
func (logs *LogRepository) SetRegExp(re *regexp.Regexp) *LogRepository {
	logs.RegExp = re

	return logs
}

// Set format
func (logs *LogRepository) SetFormat(format string) *LogRepository {
	logs.Format = format

	return logs
}

// Add log
func (logs *LogRepository) Append(log string) (*LogRepository, error) {
	var err error
	datetime := time.Now()
	matches := logs.GetRegexp().FindStringSubmatch(log)
	if len(matches) > 1 {
		datetime, err = time.Parse(logs.GetFormat(), matches[1])
		if !errors.Is(err, nil) {
			return nil, err
		}
	}

	logs.Logs = append(logs.Logs, &entity.LogEntity{
		Time: &datetime,
		Message: log,
	})

	return logs, nil
}

// Get regexp
func (logs *LogRepository) GetRegexp() *regexp.Regexp {
	return logs.RegExp
}

// Get format
func (logs *LogRepository) GetFormat() string {
	return logs.Format
}

// Get logs
func (logs *LogRepository) GetLogs() *[]*entity.LogEntity {
	return &logs.Logs
}
