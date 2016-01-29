package cron

import (
	"os"
	"testing"

	"github.com/apex/apex/function"
	"github.com/apex/apex/mock"
	"github.com/apex/apex/sources"
	"github.com/apex/apex/sources/cron/mock"
	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetHandler(discard.New())
}

func TestPlugin_Run_buildHook(t *testing.T) {
	p := &Plugin{}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cloudwatchService := mock_cloudwatcheventsiface.NewMockCloudWatchEventsAPI(mockCtrl)
	lambdaService := mock_lambdaiface.NewMockLambdaAPI(mockCtrl)

	cloudwatchService.EXPECT().PutRule(&cloudwatchevents.PutRuleInput{
		Description:        aws.String("Cron Description"),
		Name:               aws.String("Cron Name"),
		ScheduleExpression: aws.String("rate(5 minutes)"),
		State:              aws.String("ENABLED"),
	}).Return(nil, nil)

	f := &function.Function{
		Log:          log.Log,
		Path:         os.TempDir(),
		FunctionName: "testfunction",
		FunctionArn:  "arn:function",
		Service:      lambdaService,
		CloudWatch:   cloudwatchService,
		Config: function.Config{
			Sources: sources.Sources{
				Schedule: "rate(5 minutes)",
			},
		},
	}

	err := p.Run(function.DeployHook, f)
	assert.NoError(t, err)
}
