// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package instancepoller

import (
	"testing"

	gc "gopkg.in/check.v1"

	"github.com/juju/juju/state"
)

//go:generate go run github.com/golang/mock/mockgen -package instancepoller -destination package_mock_test.go github.com/juju/juju/apiserver/facades/controller/instancepoller StateLinkLayerDevice,StateLinkLayerDeviceAddress

func TestPackage(t *testing.T) {
	gc.TestingT(t)
}

type Patcher interface {
	PatchValue(ptr, value interface{})
}

func PatchState(p Patcher, st StateInterface) {
	p.PatchValue(&getState, func(*state.State, *state.Model) StateInterface {
		return st
	})
}
