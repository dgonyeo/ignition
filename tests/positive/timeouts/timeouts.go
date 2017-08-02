// Copyright 2017 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package timeouts

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/coreos/ignition/tests/types"
)

var (
	respondDelayServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Hold the connection open for 11 seconds, then return
		time.Sleep(time.Second * 11)
		w.WriteHeader(http.StatusOK)
	}))

	lastResponse           time.Time
	respondThrottledServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (lastResponse != time.Time{}) && time.Since(lastResponse) > time.Second*4 {
			// Only respond successfully if it's been more than 4 seconds since
			// the last attempt
			w.WriteHeader(http.StatusOK)
			return
		}
		lastResponse = time.Now()
		w.WriteHeader(http.StatusInternalServerError)
	}))
)

func Increase_HTTP_Response_Headers_Timeout() types.Test {
	name := "Increase HTTP Response Headers Timeout"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	var mntDevices []types.MntDevice
	config := fmt.Sprintf(`{
		"ignition": {
			"version": "2.1.0",
			"timeouts": {
				"httpResponseHeaders": 12
			}
		},
		"storage": {
		    "files": [
			    {
					"filesystem": "root",
					"path": "/foo/bar",
					"contents": {
						"source": %q
					}
				}
			]
		}
	}`, respondDelayServer.URL)
	out[0].Partitions.AddFiles("ROOT", []types.File{
		{
			Node: types.Node{
				Name: "bar",
				Path: "foo",
			},
			Contents: []string{""},
		},
	})

	return types.Test{name, in, out, mntDevices, config}
}

func Confirm_HTTP_Backoff_Works() types.Test {
	name := "Confirm HTTP Backoff Works"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	var mntDevices []types.MntDevice
	config := fmt.Sprintf(`{
		"ignition": {
			"version": "2.1.0"
		},
		"storage": {
		    "files": [
			    {
					"filesystem": "root",
					"path": "/foo/bar",
					"contents": {
						"source": %q
					}
				}
			]
		}
	}`, respondThrottledServer.URL)
	out[0].Partitions.AddFiles("ROOT", []types.File{
		{
			Node: types.Node{
				Name: "bar",
				Path: "foo",
			},
			Contents: []string{""},
		},
	})

	return types.Test{name, in, out, mntDevices, config}
}
