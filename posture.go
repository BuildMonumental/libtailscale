// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause

package main

import "C"

import (
	"sync"

	_ "tailscale.com/feature/posture"
	_ "tailscale.com/feature/syspolicy"
	"tailscale.com/util/syspolicy"
	"tailscale.com/util/syspolicy/pkey"
	"tailscale.com/util/syspolicy/rsop"
	"tailscale.com/util/syspolicy/setting"
)

// devicePolicy is the device-scope system policy of the process. It holds the
// serial number the host app sets, and turns posture checking on so the node
// reports that serial number when the control plane asks for it.
var devicePolicy struct {
	mu  sync.Mutex
	reg *rsop.StoreRegistration
}

// deviceStore is a policy store with the device's serial number. Apps on iOS
// cannot read the serial number themselves; device management gives it to
// them, so the host app passes it on.
type deviceStore struct {
	serialNumber string
}

func (s deviceStore) ReadString(key pkey.Key) (string, error) {
	switch key {
	case pkey.DeviceSerialNumber:
		return s.serialNumber, nil
	case pkey.PostureChecking:
		return "always", nil
	}
	return "", setting.ErrNotConfigured
}

func (deviceStore) ReadUInt64(pkey.Key) (uint64, error) { return 0, setting.ErrNotConfigured }
func (deviceStore) ReadBoolean(pkey.Key) (bool, error)  { return false, setting.ErrNotConfigured }
func (deviceStore) ReadStringArray(pkey.Key) ([]string, error) {
	return nil, setting.ErrNotConfigured
}

//export TsnetSetDeviceSerialNumber
func TsnetSetDeviceSerialNumber(str *C.char) C.int {
	if err := setDeviceSerialNumber(C.GoString(str)); err != nil {
		return -1
	}
	return 0
}

func setDeviceSerialNumber(serialNumber string) (err error) {
	devicePolicy.mu.Lock()
	defer devicePolicy.mu.Unlock()
	store := deviceStore{serialNumber: serialNumber}
	if devicePolicy.reg == nil {
		devicePolicy.reg, err = syspolicy.RegisterStore("libtailscale", setting.DeviceScope, store)
	} else {
		devicePolicy.reg, err = devicePolicy.reg.ReplaceStore(store)
	}
	return err
}
