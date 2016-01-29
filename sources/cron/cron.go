//go:generate mockgen -destination mock/cloudwatcheventsiface.go github.com/aws/aws-sdk-go/service/cloudwatchevents/cloudwatcheventsiface CloudWatchEventsAPI

package cron

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents/cloudwatcheventsiface"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
)

type failedEntryCountError struct {
	s string
}

func (e *failedEntryCountError) Error() string {
	return e.s
}

type Cron struct {
	Name              string
	Description       string
	Expression        string
	FunctionName      string
	FunctionArn       string
	CloudWatchService cloudwatcheventsiface.CloudWatchEventsAPI
	LambdaService     lambdaiface.LambdaAPI
}

func (c *Cron) put() (string, error) {
	res, err := c.CloudWatchService.PutRule(&cloudwatchevents.PutRuleInput{
		Description:        aws.String(c.Description),
		Name:               aws.String(c.Name),
		ScheduleExpression: aws.String(c.Expression),
		State:              aws.String("ENABLED"),
	})
	if err != nil {
		return "", err
	}
	return *res.RuleArn, err
}

func (c *Cron) connect(ruleName string) error {
	res, err := c.CloudWatchService.PutTargets(&cloudwatchevents.PutTargetsInput{
		Rule: aws.String(ruleName),
		Targets: []*cloudwatchevents.Target{
			&cloudwatchevents.Target{
				Arn: aws.String(c.FunctionArn),
				Id:  aws.String(fmt.Sprintf("%x", time.Now().Unix())),
			},
		},
	})
	if err != nil {
		return err
	}
	if *res.FailedEntryCount > 0.0 {
		return &failedEntryCountError{"failed entry count > 0"}
	}
	return err
}

func (c *Cron) allow() error {
	_, err := c.LambdaService.AddPermission(&lambda.AddPermissionInput{
		StatementId:  aws.String(fmt.Sprintf("allowscheduledevent%x", time.Now().Unix())),
		Action:       aws.String("lambda:InvokeFunction"),
		Principal:    aws.String("events.amazonaws.com"),
		FunctionName: &c.FunctionName,
	})

	return err
}

func (c *Cron) AddSchedule() error {
	// 1st Put schedule event to CloudWatch Events
	_, err := c.put()
	if err != nil {
		return err
	}
	// 2nd Connect the schedule event rule with our function target
	err = c.connect(c.Name)
	if err != nil {
		return err
	}
	// 3rd Update permissions, see: http://docs.aws.amazon.com/lambda/latest/dg/with-scheduled-events.html
	err = c.allow()

	return err
}
