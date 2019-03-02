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

package metric

import (
	"encoding/json"
	"time"
)

// Metric is a metric send by the iot device
// every time user logs into it
type Metric struct {
	// AccountID is a unique UUID identifying the account
	AccountID string `json:"account_id"`
	// UserID is a unique ID identityfing the user
	// activity
	UserID string `json:"user_id"`
	// Timestamp is a time as recorded by the device
	Timestamp time.Time `json:"timestamp"`
}

// String returns debug-friendly representation of the metric
func (m *Metric) String() string {
	data, err := json.Marshal(m)
	if err != nil {
		return err.Error()
	}
	return string(data)
}
