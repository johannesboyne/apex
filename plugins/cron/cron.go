// Package cron implements an scheduled event source
package cron

import (
	"log"

	"github.com/apex/apex/function"
	"github.com/apex/apex/sources/cron"
)

func init() {
	function.RegisterPlugin("cron", &Plugin{})
}

type Plugin struct{}

func (p *Plugin) Run(hook function.Hook, fn *function.Function) error {
	if hook == function.DeployHook && &fn.Sources != nil && &fn.Sources.Schedule != nil {
		return p.addCron(fn)
	}
	return nil
}

func (p *Plugin) addCron(fn *function.Function) error {
	log.Println("=== > add cron configuration")
	event := &cron.Cron{
		Name:              "Cron_" + fn.FunctionName,
		Description:       "Cron_" + fn.Description,
		Expression:        fn.Sources.Schedule,
		FunctionName:      fn.FunctionName,
		FunctionArn:       fn.FunctionArn,
		CloudWatchService: fn.CloudWatch,
		LambdaService:     fn.Service,
	}

	err := event.AddSchedule()
	return err
}
