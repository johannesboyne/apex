// Package cron implements an scheduled event source
package cron

import (
	"fmt"

	"github.com/apex/apex/function"
)

func init() {
	function.RegisterPlugin("cron", &Plugin{})
}

type Plugin struct{}

func (p *Plugin) Run(hook function.Hook, fn *function.Function) error {
	if hook != function.OpenHook {
		return nil
	}

	if hook == function.DeployHook {
		return p.addCron(fn)
	}
	return nil
}

func (p *Plugin) addCron(fn *function.Function) error {
	fn.Log.Debug("add cron configuration")
	fn.Log.Debug("FN")
	fn.Log.Debug(fmt.Sprintf("%+v", fn))
	//@TODO(jb): continue
	return nil
}
