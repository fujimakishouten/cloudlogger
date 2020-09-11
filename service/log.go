package service

import (
	"errors"
	"regexp"

	"github.com/fujimakishouten/cloudlogger/repository"
)

// Load log
func Load(re string, format string, logs []string) (*repository.LogRepository, error) {
	matcher, err := regexp.Compile(re)
	if !errors.Is(err, nil) {
		return nil, err
	}

	result := repository.LogRepository{
		RegExp: matcher,
		Format: format,
		Logs:   nil,
	}
	for _, log := range logs {
		_, err := result.Append(log)
		if !errors.Is(err, nil) {
			return nil, err
		}
	}

	return &result, nil
}
