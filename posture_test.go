// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"testing"

	"tailscale.com/util/syspolicy/pkey"
	"tailscale.com/util/syspolicy/ptype"
	"tailscale.com/util/syspolicy/rsop"
	"tailscale.com/util/syspolicy/setting"
)

func TestDeviceSerialNumber(t *testing.T) {
	// Tailscale registers its policy settings outside of tests only.
	if err := setting.SetDefinitionsForTest(t,
		setting.NewDefinition(pkey.DeviceSerialNumber, setting.DeviceSetting, setting.StringValue),
		setting.NewDefinition(pkey.PostureChecking, setting.DeviceSetting, setting.PreferenceOptionValue),
	); err != nil {
		t.Fatal(err)
	}
	policy, err := rsop.PolicyFor(setting.DeviceScope)
	if err != nil {
		t.Fatal(err)
	}
	for _, serial := range []string{"F2LXK0Q1", "HQNC7X3422"} {
		if err := setDeviceSerialNumber(serial); err != nil {
			t.Fatal(err)
		}
		snapshot, err := policy.Reload()
		if err != nil {
			t.Fatal(err)
		}
		if got := snapshot.Get(pkey.DeviceSerialNumber); got != serial {
			t.Errorf("DeviceSerialNumber = %v, want %v", got, serial)
		}
		if got := snapshot.Get(pkey.PostureChecking); got != ptype.AlwaysByPolicy {
			t.Errorf("PostureChecking = %v, want always", got)
		}
	}
}
