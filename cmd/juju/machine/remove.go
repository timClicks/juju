// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package machine

import (
	"time"

	"github.com/juju/cmd"
	"github.com/juju/errors"
	"github.com/juju/gnuflag"
	"gopkg.in/juju/names.v2"

	"github.com/juju/juju/api"
	"github.com/juju/juju/api/machinemanager"
	"github.com/juju/juju/apiserver/params"
	jujucmd "github.com/juju/juju/cmd"
	"github.com/juju/juju/cmd/juju/block"
	"github.com/juju/juju/cmd/modelcmd"
)

// NewRemoveCommand returns a command used to remove a specified machine.
func NewRemoveCommand() cmd.Command {
	return modelcmd.Wrap(&removeCommand{})
}

// removeCommand causes an existing machine to be destroyed.
type removeCommand struct {
	baseMachinesCommand
	apiRoot      api.Connection
	machineAPI   RemoveMachineAPI
	MachineIds   []string
	KeepInstance bool

	// Force carries out all steps required to remove a
	// machine, without halting when any of the intermediate
	// steps encounters an error.
	Force bool

	// NoWait instructs Juju that it should proceed from
	// step to step in the remove process as quickly as
	// possible, without waiting for intermediate steps
	// to be finished.
	//
	// Ignored when Force is false.
	NoWait bool
}

const destroyMachineDoc = `
Machines are specified by their numbers, which may be retrieved from the
output of ` + "`juju status`." + `
Machines responsible for the model cannot be removed.
Machines running units or containers can be removed using the '--force'
option; this will also remove those units and containers without giving
them an opportunity to shut down cleanly.

Examples:

- Remove machine 5:

    juju remove-machine 5

  This returns an error if the machine has live units or containers.
  Either stop them manually and retry the command, or add the --force flag.

- Remove machine 5 with running units or containers:

    juju remove-machine 5 --force

  To prevent Juju from waiting until each step has successfully completed
  before proceeding to the next one, add --no-wait to --force.

- Remove machine 5 with multiple live units or containers as fast as
  possible:

    juju remove-machine 5 --force --no-wait

  The --no-wait flag avoids waiting for one step to finish before proceeding
  to the next one.

- Remove machine 7 from the Juju model but do not stop
  the corresponding cloud instance:

    juju remove-machine 7 --keep-instance

  Once removed from the model, Juju will no longer have access to the 
  cloud instance. To shut it down, use the provider's console directly.

See also:
    add-machine
`

// Info implements Command.Info.
func (c *removeCommand) Info() *cmd.Info {
	return jujucmd.Info(&cmd.Info{
		Name:    "remove-machine",
		Args:    "<machine number> ...",
		Purpose: "Removes one or more machines from a model.",
		Doc:     destroyMachineDoc,
	})
}

// SetFlags implements Command.SetFlags.
func (c *removeCommand) SetFlags(f *gnuflag.FlagSet) {
	c.ModelCommandBase.SetFlags(f)
	f.BoolVar(&c.Force, "force", false, "Completely remove a machine and all its dependencies")
	f.BoolVar(&c.NoWait, "no-wait", false, "Avoid waiting for each step to complete through the removal process. Has no effect unless --force is specified.")
	f.BoolVar(&c.KeepInstance, "keep-instance", false, "Do not stop the running cloud instance")
}

func (c *removeCommand) Init(args []string) error {
	if len(args) == 0 {
		return errors.Errorf("no machines specified")
	}
	for _, id := range args {
		if !names.IsValidMachine(id) {
			return errors.Errorf("invalid machine id %q", id)
		}
	}
	c.MachineIds = args
	return nil
}

type RemoveMachineAPI interface {
	DestroyMachines(machines ...string) ([]params.DestroyMachineResult, error)
	ForceDestroyMachines(machines ...string) ([]params.DestroyMachineResult, error)
	DestroyMachinesWithParams(force, keep bool, maxWait *time.Duration, machines ...string) ([]params.DestroyMachineResult, error)
	Close() error
}

// TODO(axw) 2017-03-16 #1673323
// Drop this in Juju 3.0.
type removeMachineAdapter struct {
	*api.Client
}

func (a removeMachineAdapter) DestroyMachines(machines ...string) ([]params.DestroyMachineResult, error) {
	return a.destroyMachines(a.Client.DestroyMachines, machines)
}

func (a removeMachineAdapter) ForceDestroyMachines(machines ...string) ([]params.DestroyMachineResult, error) {
	return a.destroyMachines(a.Client.ForceDestroyMachines, machines)
}

func (a removeMachineAdapter) DestroyMachinesWithParams(force, keep bool, maxWait *time.Duration, machines ...string) ([]params.DestroyMachineResult, error) {
	return a.destroyMachines(a.Client.ForceDestroyMachines, machines)
}

func (a removeMachineAdapter) destroyMachines(f func(...string) error, machines []string) ([]params.DestroyMachineResult, error) {
	if err := f(machines...); err != nil {
		return nil, err
	}
	results := make([]params.DestroyMachineResult, len(machines))
	for i := range results {
		results[i].Info = &params.DestroyMachineInfo{}
	}
	return results, nil
}

func (c *removeCommand) getAPIRoot() (api.Connection, error) {
	if c.apiRoot != nil {
		return c.apiRoot, nil
	}
	return c.NewAPIRoot()
}

func (c *removeCommand) getRemoveMachineAPI() (RemoveMachineAPI, error) {
	root, err := c.getAPIRoot()
	if err != nil {
		return nil, err
	}
	if root.BestFacadeVersion("MachineManager") < 4 && c.KeepInstance {
		return nil, errors.New("this version of Juju doesn't support --keep-instance")
	}
	if root.BestFacadeVersion("MachineManager") >= 3 && c.machineAPI == nil {
		return machinemanager.NewClient(root), nil
	}
	if c.machineAPI != nil {
		return c.machineAPI, nil
	}
	return removeMachineAdapter{root.Client()}, nil
}

// Run implements Command.Run.
func (c *removeCommand) Run(ctx *cmd.Context) error {
	client, err := c.getRemoveMachineAPI()
	if err != nil {
		return err
	}
	defer client.Close()

	var maxWait *time.Duration
	if c.NoWait {
		if !c.Force {
			return errors.NotValidf("--no-wait without --force")
		}
	} else { // sorry about the double negative
		oneMinute := time.Minute
		maxWait = &oneMinute
	}

	var results []params.DestroyMachineResult
	if c.KeepInstance {
		results, err = client.DestroyMachinesWithParams(c.Force, c.KeepInstance, maxWait, c.MachineIds...)
	} else {
		destroy := client.DestroyMachines
		if c.Force {
			destroy = client.ForceDestroyMachines
		}
		results, err = destroy(c.MachineIds...)
	}
	if err := block.ProcessBlockedError(err, block.BlockRemove); err != nil {
		return err
	}

	anyFailed := false
	for i, id := range c.MachineIds {
		result := results[i]
		if result.Error != nil {
			anyFailed = true
			ctx.Infof("removing machine %s failed: %s", id, result.Error)
			continue
		}
		if c.KeepInstance {
			ctx.Infof("removing machine %s (but retaining cloud instance)", id)
		} else {
			ctx.Infof("removing machine %s", id)
		}
		for _, entity := range result.Info.DestroyedUnits {
			unitTag, err := names.ParseUnitTag(entity.Tag)
			if err != nil {
				logger.Warningf("%s", err)
				continue
			}
			ctx.Infof("- will remove %s", names.ReadableString(unitTag))
		}
		for _, entity := range result.Info.DestroyedStorage {
			storageTag, err := names.ParseStorageTag(entity.Tag)
			if err != nil {
				logger.Warningf("%s", err)
				continue
			}
			ctx.Infof("- will remove %s", names.ReadableString(storageTag))
		}
		for _, entity := range result.Info.DetachedStorage {
			storageTag, err := names.ParseStorageTag(entity.Tag)
			if err != nil {
				logger.Warningf("%s", err)
				continue
			}
			ctx.Infof("- will detach %s", names.ReadableString(storageTag))
		}
	}

	if anyFailed {
		return cmd.ErrSilent
	}
	return nil
}
