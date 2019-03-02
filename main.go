/*
Copyright 2019 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package main

import (
	"context"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/gravitational/fakeiot/pkg/client"
	"github.com/gravitational/fakeiot/pkg/runner"
	"github.com/gravitational/fakeiot/pkg/utils"

	"github.com/gravitational/kingpin"
	"github.com/gravitational/trace"
	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
)

func main() {
	utils.InitLogger(log.InfoLevel)
	err := run()
	if err != nil {
		log.Errorf("Fake IOT program has exited with error: %v", err)
		os.Exit(255)
	}
	log.Info("Fake IOT program run successfully.")
}

func run() error {
	app := kingpin.New("fakeiot", "Fake IOT device simulator.")
	token := app.Flag("token", "Bearer token.").Required().String()
	targetURL := app.Flag("url", "URL to emit metrics to.").Required().URL()

	certPath := app.Flag("ca-cert", "Path to PEM-encoded file with trusted root CA certificate.").String()
	debug := app.Flag("debug", "Turn on debug logging.").Bool()

	test := app.Command("test", "Run compliance tests")

	run := app.Command("run", "Run simulation")
	runPeriod := run.Flag("period", "Sim duration").Default("30s").Duration()
	runAccountID := run.Flag("account-id", "Account ID").Default(uuid.New()).String()
	runFreq := run.Flag("freq", "Sim frequency").Default("1s").Duration()
	runUsers := run.Flag("users", "Run users").Default("100").Int()

	// parse CLI commands+flags:
	command, err := app.Parse(os.Args[1:])
	if err != nil {
		return trace.Wrap(err)
	}

	// apply -d flag:
	if *debug {
		utils.InitLogger(log.DebugLevel)
	}

	config := client.Config{
		URL:         *targetURL,
		BearerToken: *token,
	}
	if *certPath != "" {
		certData, err := ioutil.ReadFile(*certPath)
		if err != nil {
			return trace.ConvertSystemError(err)
		}
		cert, err := utils.ParseCertificatePEM(certData)
		if err != nil {
			return trace.Wrap(err)
		}
		config.CACert = cert
	}

	client, err := client.New(config)
	if err != nil {
		return trace.Wrap(err)
	}
	r := runner.New(runner.Config{Client: client})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		exitSignals := make(chan os.Signal, 1)
		signal.Notify(exitSignals, syscall.SIGTERM, syscall.SIGINT)

		select {
		case sig := <-exitSignals:
			log.Debugf("Got signal: %v", sig)
			cancel()
		}
	}()

	switch command {
	case test.FullCommand():
		return r.RunTests(ctx)
	case run.FullCommand():
		return r.RunSimulation(ctx, runner.Simulation{
			AccountID: *runAccountID,
			Users:     *runUsers,
			Freq:      *runFreq,
			Period:    *runPeriod,
		})
	}
	return trace.BadParameter("command %q is not supported yet", command)
}
