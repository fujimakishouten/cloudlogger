package service

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func RetrieveLogStreams(session *session.Session, group string) (*cloudwatchlogs.DescribeLogStreamsOutput, error) {
	svc := cloudwatchlogs.New(session)
	input := cloudwatchlogs.DescribeLogStreamsInput{
		Descending: aws.Bool(true),
		Limit: aws.Int64(1),
		LogGroupName: aws.String(group),
	}

	result, err := svc.DescribeLogStreams(&input)
	if !errors.Is(err, nil) {
		if aerr, ok := err.(awserr.Error); ok {
			return nil, aerr;
		}

		return nil, err;
	}

	return result, nil
}

func Send(session *session.Session, group string, stream string, token string, logs []string) error {
	events := []*cloudwatchlogs.InputLogEvent{}
	for _, log := range logs {
		events = append(events, &cloudwatchlogs.InputLogEvent{
			Message:   aws.String(log),
			Timestamp: aws.Int64(time.Now().UnixNano() / int64(time.Millisecond)),
		})
	}

	svc := cloudwatchlogs.New(session)
	input := cloudwatchlogs.PutLogEventsInput{
		LogEvents: events,
		LogGroupName: aws.String(group),
		LogStreamName: aws.String(stream),
	}
	if token != "" {
		input.SequenceToken = aws.String(token)
	}

	_, err := svc.PutLogEvents(&input)
	if !errors.Is(err, nil) {
		if aerr, ok := err.(awserr.Error); ok {
			return aerr;
		}

		return err;
	}

	return nil
}

func GetNextForwardToken(session *session.Session, group string, stream string) (string, error) {
	result, err := RetrieveLogStreams(session, group)
	if !errors.Is(err, nil) {
		return "", err
	}

	for _, log := range result.LogStreams {
		if *log.LogStreamName == stream {
			return *log.UploadSequenceToken, nil
		}
	}

	return "", nil
}
