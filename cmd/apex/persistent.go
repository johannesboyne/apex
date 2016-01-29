package main

import (
	"os"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/spf13/cobra"

	"github.com/apex/apex/dryrun"
	"github.com/apex/apex/project"
)

type persistentValues struct {
	Chdir    string
	DryRun   bool
	Env      []string
	LogLevel string
	Verbose  bool
	Yes      bool
	Profile  string

	session *session.Session
	project *project.Project
}

func (pv *persistentValues) noopRun(*cobra.Command, []string) {}

func (pv *persistentValues) preRun(c *cobra.Command, args []string) {
	if l, err := log.ParseLevel(pv.LogLevel); err == nil {
		log.SetLevel(l)
	}

	cfg := aws.NewConfig()

	if pv.Profile != "" {
		cfg = cfg.WithCredentials(credentials.NewSharedCredentials("", pv.Profile))
	}

	pv.session = session.New(cfg)

	pv.project = &project.Project{
		Log:  log.Log,
		Path: ".",
	}

	if pv.DryRun {
		log.SetLevel(log.WarnLevel)
		pv.project.Service = dryrun.New(pv.session)
		pv.project.Concurrency = 1
	} else {
		pv.project.Service = lambda.New(pv.session)
		pv.project.CloudWatchEvents = cloudwatchevents.New(pv.session)
	}

	if pv.Chdir != "" {
		if err := os.Chdir(pv.Chdir); err != nil {
			log.Fatalf("error: %s", err)
		}
	}

	if err := pv.project.Open(); err != nil {
		log.Fatalf("error opening project: %s", err)
	}
}

var pv persistentValues
