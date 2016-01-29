package logs

import (
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// Log implements log fetching and polling for CloudWatchLogs,
// and represents a single group.
type Log struct {
	Config
	GroupName string
	Log       log.Interface
	err       error
}

// Start consuming logs.
func (l *Log) Start() <-chan *Event {
	ch := make(chan *Event)
	go l.start(ch)
	return ch
}

// start consuming and exit if Follow is not enabled.
func (l *Log) start(ch chan<- *Event) {
	defer close(ch)

	l.Log.Debug("enter")
	defer l.Log.Debug("exit")

	var start = l.StartTime.UnixNano() / int64(time.Millisecond)
	var nextToken *string
	var err error

	for {
		l.Log.WithField("start", start).Debug("poll")
		nextToken, start, err = l.fetch(nextToken, start, ch)

		if err != nil {
			l.err = fmt.Errorf("log %q: %s", l.GroupName, err)
			break
		}

		if !l.Follow {
			break
		}

		time.Sleep(l.PollInterval)
	}
}

// fetch logs relative to the given token and start time. We ignore when the log group is not found.
func (l *Log) fetch(nextToken *string, start int64, ch chan<- *Event) (*string, int64, error) {
	res, err := l.Service.FilterLogEvents(&cloudwatchlogs.FilterLogEventsInput{
		LogGroupName:  &l.GroupName,
		FilterPattern: &l.FilterPattern,
		StartTime:     &start,
		NextToken:     nextToken,
	})

	if e, ok := err.(awserr.Error); ok {
		if e.Code() == "ResourceNotFoundException" {
			l.Log.Debug("not found")
			return nil, 0, nil
		}
	}

	if err != nil {
		return nil, 0, err
	}

	for _, event := range res.Events {
		start = *event.Timestamp + 1
		ch <- &Event{
			GroupName: l.GroupName,
			Message:   *event.Message,
		}
	}

	return res.NextToken, start, nil
}

// Err returns the first error, if any, during processing.
func (l *Log) Err() error {
	return l.err
}
