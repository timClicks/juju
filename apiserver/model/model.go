// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package model

import (
	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/state"
)

func init() {
	common.RegisterStandardFacade("Model", 1, NewModelAPI)
}

// ModelAPI implements the API used by the machine model worker.
type ModelAPI struct {
	*common.EnvironWatcher
	*ModelTools
}

// NewModelAPI creates a new instance of the Model API.
func NewModelAPI(st *state.State, resources *common.Resources, authorizer common.Authorizer) (*ModelAPI, error) {
	return &ModelAPI{
		EnvironWatcher: common.NewEnvironWatcher(st, resources, authorizer),
		ModelTools:     NewEnvironTools(st, authorizer),
	}, nil
}
