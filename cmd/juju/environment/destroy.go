// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for infos.

package environment

import (
	"fmt"

	"github.com/juju/cmd"
	"github.com/juju/errors"
	"launchpad.net/gnuflag"

	"github.com/juju/juju/apiserver/params"
	jujucmd "github.com/juju/juju/cmd"
	"github.com/juju/juju/cmd/envcmd"
	"github.com/juju/juju/cmd/juju/block"
	"github.com/juju/juju/environs"
	"github.com/juju/juju/environs/configstore"
)

func NewDestroyCommand() cmd.Command {
	return envcmd.Wrap(
		&destroyCommand{},
		envcmd.EnvSkipDefault,
		envcmd.EnvSkipFlags,
	)
}

// destroyCommand destroys the specified environment.
type destroyCommand struct {
	envcmd.EnvCommandBase
	envName   string
	assumeYes bool
	api       DestroyEnvironmentAPI
}

var destroyDoc = `Destroys the specified model`
var destroyEnvMsg = `
WARNING! This command will destroy the %q model.
This includes all machines, services, data and other resources.

Continue [y/N]? `[1:]

// DestroyEnvironmentAPI defines the methods on the modelmanager
// API that the destroy command calls. It is exported for mocking in tests.
type DestroyEnvironmentAPI interface {
	Close() error
	DestroyModel() error
}

// Info implements Command.Info.
func (c *destroyCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "destroy-model",
		Args:    "<model name>",
		Purpose: "terminate all machines and other associated resources for a non-controller model",
		Doc:     destroyDoc,
	}
}

// SetFlags implements Command.SetFlags.
func (c *destroyCommand) SetFlags(f *gnuflag.FlagSet) {
	f.BoolVar(&c.assumeYes, "y", false, "Do not ask for confirmation")
	f.BoolVar(&c.assumeYes, "yes", false, "")
}

// Init implements Command.Init.
func (c *destroyCommand) Init(args []string) error {
	switch len(args) {
	case 0:
		return errors.New("no model specified")
	case 1:
		c.envName = args[0]
		c.SetEnvName(c.envName)
		return nil
	default:
		return cmd.CheckEmpty(args[1:])
	}
}

func (c *destroyCommand) getAPI() (DestroyEnvironmentAPI, error) {
	if c.api != nil {
		return c.api, nil
	}
	return c.NewAPIClient()
}

// Run implements Command.Run
func (c *destroyCommand) Run(ctx *cmd.Context) error {
	store, err := configstore.Default()
	if err != nil {
		return errors.Annotate(err, "cannot open model info storage")
	}

	cfgInfo, err := store.ReadInfo(c.envName)
	if err != nil {
		return errors.Annotate(err, "cannot read model info")
	}

	// Verify that we're not destroying a controller
	apiEndpoint := cfgInfo.APIEndpoint()
	if apiEndpoint.ServerUUID != "" && apiEndpoint.EnvironUUID == apiEndpoint.ServerUUID {
		return errors.Errorf("%q is a controller; use 'juju destroy-controller' to destroy it", c.envName)
	}

	if !c.assumeYes {
		fmt.Fprintf(ctx.Stdout, destroyEnvMsg, c.envName)

		if err := jujucmd.UserConfirmYes(ctx); err != nil {
			return errors.Annotate(err, "model destruction")
		}
	}

	// Attempt to connect to the API.  If we can't, fail the destroy.
	api, err := c.getAPI()
	if err != nil {
		return errors.Annotate(err, "cannot connect to API")
	}
	defer api.Close()

	// Attempt to destroy the environment.
	err = api.DestroyModel()
	if err != nil {
		return c.handleError(errors.Annotate(err, "cannot destroy model"))
	}

	return environs.DestroyInfo(c.envName, store)
}

func (c *destroyCommand) handleError(err error) error {
	if err == nil {
		return nil
	}
	if params.IsCodeOperationBlocked(err) {
		return block.ProcessBlockedError(err, block.BlockDestroy)
	}
	logger.Errorf(`failed to destroy model %q`, c.envName)
	return err
}
