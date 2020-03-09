package entity

import (
	"time"
)

type LogEntity struct {
	Time *time.Time
	Message string
}

// Set time
func (log *LogEntity) SetTime(datetime *time.Time) *LogEntity {
	log.Time = datetime

	return log
}

// Set message
func (log *LogEntity) SetMessage(message string) *LogEntity {
	log.Message = message

	return log
}

// Get time
func (log *LogEntity) GetTime() *time.Time {
	return log.Time
}

// Get message
func (log *LogEntity) GetMessage() string {
	return log.Message
}
