package main

import (
	"cloudlogger/service"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"os"
	"strings"
)

const version = "1.0.0"
var (
	app = kingpin.New("clowdlogger", "Send log data to AWS CloudWatch.")
	group = app.Flag("group", "Log group name").Short('g').Required().String()
	stream = app.Flag("stream", "Log stream name").Short('s').Required().String()
	log = app.Arg("log", "Log text.").String()
)

func main() {
	app.Version(version)
	_, err := app.Parse(os.Args[1:])
	if !errors.Is(err, nil) {
		app.FatalUsage(err.Error())
	}

	status, err := os.Stdin.Stat()
	if !errors.Is(err, nil) {
		app.FatalUsage(err.Error())
	}

	if status.Size() > 0 {
		data, err := ioutil.ReadAll(os.Stdin)
		if !errors.Is(err, nil) {
			app.FatalUsage(err.Error())
		}
		*log = string(data)
	}

	if *log == "" {
		app.FatalUsage("Options: log are required.")
	}

	sess := session.Must(session.NewSession(&aws.Config{}))
	token, err := service.GetNextForwardToken(sess, *group, *stream)
	if !errors.Is(err, nil) {
		app.FatalUsage(err.Error())
	}

	err = service.Send(
		sess,
		*group,
		*stream,
		token,
		strings.Split(*log, "\n"))
	if !errors.Is(err, nil) {
		app.FatalUsage(err.Error())
	}
}
