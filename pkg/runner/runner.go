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

package runner

import (
	"context"
	"time"

	"github.com/gravitational/fakeiot/pkg/client"
	"github.com/gravitational/fakeiot/pkg/metric"
	"github.com/pborman/uuid"

	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
)

// Config is a runner configuration
type Config struct {
	// Client is a configured client
	Client *client.Client
}

// New returns a new instance of runner
func New(config Config) *Runner {
	return &Runner{
		Config: config,
		Entry:  log.WithFields(log.Fields{trace.Component: "runner"}),
	}
}

// Runner is Fake IOT device runner
type Runner struct {
	// Config is a device configuration
	Config
	*log.Entry
}

// Simulation is a IOT device simulation
type Simulation struct {
	// Period is a period to run simulation for
	Period time.Duration
	// Freq is a frequency of activity
	Freq time.Duration
	// AccountID is an account ID to emit
	AccountID string
	// Users is a total amount of different users to simulate accross
	// all accounts. Simulation will generate distinct user IDs
	// and emit them in a rotating loop
	Users int
}

// RunSimulation runs test simulation
func (r *Runner) RunSimulation(ctx context.Context, sim Simulation) error {
	r.Infof("Starting simulation with %v users over %v period with %v frequency.", sim.Users, sim.Period, sim.Freq)
	tctx, cancel := context.WithTimeout(ctx, sim.Period+time.Second)
	defer cancel()

	var users []string
	for i := 0; i < sim.Users; i++ {
		users = append(users, uuid.New())
	}

	t := time.NewTicker(sim.Freq)
	defer t.Stop()
	var i int
	for {
		user := users[i]
		i = (i + 1) % len(users)
		metric := metric.Metric{
			AccountID: sim.AccountID,
			UserID:    user,
			Timestamp: time.Now().UTC(),
		}
		err := r.Client.Send(tctx, metric)
		if err != nil {
			return trace.Wrap(err)
		}
		r.Infof("Sent %v.", metric.String())
		select {
		case <-time.After(sim.Period):
			return nil
		case <-tctx.Done():
			return nil
		case <-t.C:

		}
	}
}

// RunTests runs a series of compliance tests
// to make sure that server performs as expected
func (r *Runner) RunTests(ctx context.Context) error {
	r.Info("Starting compliance tests.")
	var errors []error
	for _, tc := range []func(ctx context.Context) error{
		r.SendOK,
		r.SendEmptyMetric,
		r.SendCorruptedMetric,
		func(ctx context.Context) error { return r.SendBogusAuth(ctx, "") },
		func(ctx context.Context) error { return r.SendBogusAuth(ctx, uuid.New()) },
	} {
		err := tc(ctx)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return trace.NewAggregate(errors...)
}

// SendOK sends OK metric and expects the metric to succeed
func (r *Runner) SendOK(ctx context.Context) error {
	r.Debugf("[TEST] Sending OK request...")
	err := r.Client.Send(ctx, metric.Metric{
		AccountID: ValidTestAccountID,
		UserID:    ValidTestUserID,
		Timestamp: time.Now().UTC(),
	})
	if err != nil {
		r.Errorf("[FAIL] Sending OK request.")
		return trace.BadParameter("expected success, have received %q error from server", err)
	}
	r.Infof("[PASS] Sending OK request.")
	return nil
}

// SendBogusAuth sends a request with bogus authentication and expects to fail
func (r *Runner) SendBogusAuth(ctx context.Context, bearerToken string) error {
	r.Debugf("[TEST] Sending invalid authentication with bearer auth token %q.", bearerToken)
	client, err := client.New(client.Config{
		URL:    r.Client.URL,
		CACert: r.Client.CACert,
	})
	if err != nil {
		return trace.Wrap(err)
	}
	err = client.Send(ctx, metric.Metric{
		AccountID: ValidTestAccountID,
		UserID:    ValidTestUserID,
		Timestamp: time.Now().UTC(),
	})
	if err == nil {
		r.Errorf("[FAIL] Sending bogus request.")
		return trace.BadParameter("expected access denied HTTP error, have received OK from server")
	}
	if !trace.IsAccessDenied(err) {
		r.Errorf("[FAIL] Sending bogus request.")
		return trace.BadParameter("expected access denied HTTP error, have received %q error from server", err)
	}
	r.Infof("[PASS] Sending bogus request.")
	return nil
}

// SendEmptyMetric sends empty metric and expects an error
func (r *Runner) SendEmptyMetric(ctx context.Context) error {
	r.Debugf("[TEST] Sending empty request...")
	err := r.Client.Send(ctx, metric.Metric{})
	if err == nil {
		r.Errorf("[FAIL] Sending empty request.")
		return trace.BadParameter("runner sent an empty metric and received OK response, we expected HTTP Bad request error.")
	}
	if !trace.IsBadParameter(err) {
		log.Warningf("runner sent an empty metric and received an error %v, however bad request response would have been better", err)
	}
	r.Infof("[PASS] Sending empty request.")
	return nil
}

// SendCorruptedMetric sends corrupted metric
func (r *Runner) SendCorruptedMetric(ctx context.Context) error {
	r.Debugf("[TEST] Sending corrupted request...")
	err := r.Client.SendCorruptedData(ctx)
	if err == nil {
		r.Errorf("[FAIL] Sending corrupted request.")
		return trace.BadParameter("runner sent non-JSON metric and received OK response, we expected HTTP Bad request error.")
	}
	if !trace.IsBadParameter(err) {
		log.Warningf("runner sent non-JSON metric and received an error %v, however bad request response would have been better", err)
	}
	r.Infof("[PASS] Sending corrupted request.")
	return nil
}

const (
	// ValidTestAccountID is a valid test account that should be recognized
	// by the system
	ValidTestAccountID = "testacct-0000-0000-0000-000000000000"

	// ValidTestUserID is a valid test user id that should be recognized by the system
	ValidTestUserID = "testuser-0000-0000-0000-000000000000"
)
