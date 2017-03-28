// Copyright 2015 CoreOS, Inc.
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

// The storage stage is responsible for partitioning disks, creating RAID
// arrays, formatting partitions, writing files, writing systemd units, and
// writing network units.

package luks

import (
	"fmt"
	"net/url"

	"github.com/coreos/ignition/config/types"
	"github.com/coreos/ignition/internal/exec/stages"
	"github.com/coreos/ignition/internal/exec/util"
	"github.com/coreos/ignition/internal/log"
	"github.com/coreos/ignition/internal/resource"
	"github.com/coreos/ignition/internal/systemd"
	"github.com/martinjungblut/cryptsetup"

	"golang.org/x/net/context"
)

const (
	name = "luks"
)

func init() {
	stages.Register(creator{})
}

type creator struct{}

func (creator) Create(logger *log.Logger, client *resource.HttpClient, root string) stages.Stage {
	return &stage{
		Util: util.Util{
			DestDir: root,
			Logger:  logger,
		},
		client: client,
	}
}

func (creator) Name() string {
	return name
}

type stage struct {
	util.Util

	client *resource.HttpClient
}

func (stage) Name() string {
	return name
}

func (s stage) Run(config types.Config) bool {
	if err := s.activateLukss(config); err != nil {
		s.Logger.Crit("activate luks device failed: %v", err)
		return false
	}

	return true
}

// waitOnDevices waits for the devices enumerated in devs as a logged operation
// using ctxt for the logging and systemd unit identity.
func (s stage) waitOnDevices(devs []string, ctxt string) error {
	if err := s.LogOp(
		func() error { return systemd.WaitOnDevices(devs, ctxt) },
		"waiting for devices %v", devs,
	); err != nil {
		return fmt.Errorf("failed to wait on %s devs: %v", ctxt, err)
	}

	return nil
}

// createDeviceAliases creates device aliases for every device in devs.
func (s stage) createDeviceAliases(devs []string) error {
	for _, dev := range devs {
		target, err := util.CreateDeviceAlias(dev)
		if err != nil {
			return fmt.Errorf("failed to create device alias for %q: %v", dev, err)
		}
		s.Logger.Info("created device alias for %q: %q -> %q", dev, util.DeviceAlias(dev), target)
	}

	return nil
}

// waitOnDevicesAndCreateAliases simply wraps waitOnDevices and createDeviceAliases.
func (s stage) waitOnDevicesAndCreateAliases(devs []string, ctxt string) error {
	if err := s.waitOnDevices(devs, ctxt); err != nil {
		return err
	}

	if err := s.createDeviceAliases(devs); err != nil {
		return err
	}

	return nil
}

// activateLukss activates any LUKS partitions described in config.Storage.Luks
func (s stage) activateLukss(config types.Config) error {
	lukss := make([]types.Luks, 0, len(config.Storage.Luks))
	for _, luks := range config.Storage.Luks {
		if len(luks.Device) > 0 {
			lukss = append(lukss, luks)
		}
	}

	if len(lukss) == 0 {
		return nil
	}
	s.Logger.PushPrefix("activateLuks")
	defer s.Logger.PopPrefix()

	devs := []string{}
	for _, luks := range lukss {
		devs = append(devs, string(luks.Device))
	}

	if err := s.waitOnDevicesAndCreateAliases(devs, "luks"); err != nil {
		return err
	}

	for _, luks := range lukss {
		keyUrl, err := url.Parse(luks.Key)
		if err != nil {
			return fmt.Errorf("Unable to parse key URL: %v", err)
		}

		key, err := resource.Fetch(s.Logger, s.client, context.Background(), *keyUrl)
		if err != nil {
			return fmt.Errorf("Unable to fetch key from %s: %v", keyUrl, err)
		}

		err, device := cryptsetup.Init(string(luks.Device))
		if err != nil {
			return fmt.Errorf("Unable to initialise cryptsetup: %v", err)
		}

		err = device.Activate(luks.Name, cryptsetup.CRYPT_ANY_SLOT, string(key), 0)
		if err != nil {
			return fmt.Errorf("Error activating LUKS device %s: %v", string(luks.Device), err)
		}
	}

	return nil
}
