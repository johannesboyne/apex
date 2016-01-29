// Package cron implements an scheduled event source
package cron

import (
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
	if &fn.Sources != nil && &fn.Sources.Schedule != nil {
		return p.addCron(fn)
	}
	return nil
}

// addCron builds the Cron Event.
func (p *Plugin) addCron(fn *function.Function) error {
	event := &cron.Cron{
		Name:              "Cron_" + fn.FunctionName,
		Description:       "Cron_" + fn.FunctionName,
		Expression:        fn.Sources.Schedule,
		FunctionName:      fn.FunctionName,
		FunctionArn:       fn.FunctionArn,
		CloudWatchService: fn.CloudWatchEvents,
		LambdaService:     fn.Service,
	}

	err := event.AddSchedule()
	return err
}
