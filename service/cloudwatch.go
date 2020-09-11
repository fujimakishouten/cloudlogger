package service

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"

	"github.com/fujimakishouten/cloudlogger/repository"
)

func CreateLogGroup(svc *cloudwatchlogs.CloudWatchLogs, group string) error {
	input := cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(group),
	}

	_, err := svc.CreateLogGroup(&input)
	if !errors.Is(err, nil) {
		return err
	}

	return nil
}

func DescribeLogGroup(svc *cloudwatchlogs.CloudWatchLogs, group string) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
	var describe func(token string) (*cloudwatchlogs.DescribeLogGroupsOutput, error)
	describe = func(token string) (*cloudwatchlogs.DescribeLogGroupsOutput, error) {
		input := cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: aws.String(group),
		}
		if token != "" {
			input.NextToken = aws.String(token)
		}

		result, err := svc.DescribeLogGroups(&input)
		if !errors.Is(err, nil) {
			if aerr, ok := err.(awserr.Error); ok {
				return nil, aerr
			}
			return nil, err
		}

		for _, log := range result.LogGroups {
			if *log.LogGroupName == group {
				return result, nil
			}
		}

		if result.NextToken != nil {
			return describe(*result.NextToken)
		}

		return nil, nil
	}

	return describe("")
}

func EnsureLogGroup(svc *cloudwatchlogs.CloudWatchLogs, group string) error {
	result, err := DescribeLogGroup(svc, group)
	if !errors.Is(err, nil) {
		return err
	}
	if result == nil {
		return CreateLogGroup(svc, group)
	}

	return nil
}

func CreateLogStream(svc *cloudwatchlogs.CloudWatchLogs, group string, stream string) error {
	input := cloudwatchlogs.CreateLogStreamInput{
		LogGroupName: aws.String(group),
		LogStreamName: aws.String(stream),
	}

	_, err := svc.CreateLogStream(&input)
	if !errors.Is(err, nil) {
		return err
	}

	return nil
}

func DescribeLogStream(svc *cloudwatchlogs.CloudWatchLogs, group string, stream string) (*cloudwatchlogs.DescribeLogStreamsOutput, error) {
	var describe func(token string) (*cloudwatchlogs.DescribeLogStreamsOutput, error)
	describe = func(token string) (*cloudwatchlogs.DescribeLogStreamsOutput, error) {
		input := cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName: aws.String(group),
			LogStreamNamePrefix: aws.String(stream),
		}
		if token != "" {
			input.NextToken = aws.String(token)
		}

		result, err := svc.DescribeLogStreams(&input)
		if !errors.Is(err, nil) {
			if aerr, ok := err.(awserr.Error); ok {
				return nil, aerr
			}
			return nil, err
		}

		for _, log := range result.LogStreams {
			if *log.LogStreamName == stream {
				return result, nil
			}
		}

		if result.NextToken != nil {
			return describe(*result.NextToken)
		}

		return nil, nil
	}

	return describe("")
}

func EnsureLogStream(svc *cloudwatchlogs.CloudWatchLogs, group string, stream string) error {
	result, err := DescribeLogStream(svc, group, stream)
	if !errors.Is(err, nil) {
		return err
	}

	if result == nil {
		return CreateLogStream(svc, group, stream)
	}

	return nil
}

func GetSequenceToken(svc *cloudwatchlogs.CloudWatchLogs, group string, stream string) (string, error) {
	result, err := DescribeLogStream(svc, group, stream)
	if !errors.Is(err, nil) {
		return "", err
	}


	for _, log := range result.LogStreams {
		if *log.LogStreamName == stream {
			if log.UploadSequenceToken == nil {
				return "", nil
			}

			return *log.UploadSequenceToken, nil
		}
	}

	return "", nil
}


func Send(svc *cloudwatchlogs.CloudWatchLogs, group string, stream string, logs *repository.LogRepository) error {
	events := []*cloudwatchlogs.InputLogEvent{}
	for _, log := range *logs.GetLogs() {
		events = append(events, &cloudwatchlogs.InputLogEvent{
			Message:   aws.String(log.GetMessage()),
			Timestamp: aws.Int64(log.GetTime().UnixNano() / int64(time.Millisecond)),
		})
	}

	input := cloudwatchlogs.PutLogEventsInput{
		LogEvents: events,
		LogGroupName: aws.String(group),
		LogStreamName: aws.String(stream),
	}
	token, err := GetSequenceToken(svc, group, stream)
	if !errors.Is(err, nil) {
		if aerr, ok := err.(awserr.Error); ok {
			return aerr
		}

		return err
	}
	if token != "" {
		input.SequenceToken = aws.String(token)
	}

	_, err = svc.PutLogEvents(&input)
	if !errors.Is(err, nil) {
		if aerr, ok := err.(awserr.Error); ok {
			return aerr;
		}

		return err;
	}

	return nil
}
