package cron

import (
	"fmt"
	"testing"
	"time"

	"github.com/apex/apex/event_sources/cron/mock"
	"github.com/apex/apex/mock"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/golang/mock/gomock"
)

func TestPut(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	cloudwatchService := mock_cloudwatcheventsiface.NewMockCloudWatchEventsAPI(mockCtrl)
	lambdaService := mock_lambdaiface.NewMockLambdaAPI(mockCtrl)

	cloudwatchService.EXPECT().PutRule(&cloudwatchevents.PutRuleInput{
		Description:        aws.String("Test"),
		Name:               aws.String("Name"),
		ScheduleExpression: aws.String("rate(5 minutes)"),
		State:              aws.String("ENABLED"),
	}).Return(&cloudwatchevents.PutRuleOutput{
		RuleArn: aws.String("eventsourcearn"),
	}, nil)
	cloudwatchService.EXPECT().PutTargets(&cloudwatchevents.PutTargetsInput{
		Rule: aws.String("Name"),
		Targets: []*cloudwatchevents.Target{
			&cloudwatchevents.Target{
				Arn: aws.String("functionarn"),
				Id:  aws.String(fmt.Sprintf("%x", time.Now().Unix())),
			},
		},
	}).Return(&cloudwatchevents.PutTargetsOutput{
		FailedEntryCount: aws.Int64(0.0),
	}, nil)

	lambdaService.EXPECT().AddPermission(&lambda.AddPermissionInput{
		StatementId:  aws.String("allowscheduledevent"),
		Action:       aws.String("lambda:InvokeFunction"),
		Principal:    aws.String("events.amazonaws.com"),
		FunctionName: aws.String("functionname"),
	}).Return(&lambda.AddPermissionOutput{
		Statement: aws.String("jsonresponse"),
	}, nil)

	cronEvent := &Cron{
		Name:              "Name",
		Description:       "Test",
		Expression:        "rate(5 minutes)",
		FunctionName:      "functionname",
		FunctionArn:       "functionarn",
		CloudWatchService: cloudwatchService,
		LambdaService:     lambdaService,
	}
	err := cronEvent.AddSchedule()
	if err != nil {
		t.Errorf("error")
	}
}
