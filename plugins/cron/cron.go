// Package cron implements an scheduled event source
package cron

import (
	"fmt"

	"github.com/apex/apex/function"
	"github.com/apex/apex/sources/cron"
)

func init() {
	function.RegisterPlugin("cron", &Plugin{})
}

// Plugin implementation.
type Plugin struct{}

// PostDeploy adds the Cron Scheduled Event.
func (p *Plugin) PostDeploy(fn *function.Function) error {
	if &fn.Sources == nil || fn.Sources.Schedule == "" {
		return nil
	}
	return p.addCron(fn)
}

// addCron builds the Cron Event.
func (p *Plugin) addCron(fn *function.Function) error {
	config, err := fn.GetConfig()

	if err != nil {
		return err
	}

	event := &cron.Cron{
		Name:              fmt.Sprintf("cron_%s", fn.FunctionName),
		Description:       fmt.Sprintf("cron for lambda function %s", fn.FunctionName),
		Expression:        fn.Sources.Schedule,
		FunctionName:      fn.FunctionName,
		FunctionArn:       *config.Configuration.FunctionArn,
		CloudWatchService: fn.CloudWatchEvents,
		LambdaService:     fn.Service,
	}

	err = event.AddSchedule()
	return err
}
