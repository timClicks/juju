// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package application

import (
	"github.com/juju/cmd"
	"github.com/juju/errors"
	"github.com/juju/juju/api"
	"github.com/juju/juju/api/annotations"
	"github.com/juju/juju/api/application"
	"github.com/juju/juju/api/base"
	apicharms "github.com/juju/juju/api/charms"
	"github.com/juju/juju/api/modelconfig"
	"github.com/juju/juju/charmstore"
	"github.com/juju/juju/cmd/modelcmd"
	"github.com/juju/juju/jujuclient"
	"github.com/juju/juju/resource/resourceadapters"
	"github.com/juju/romulus"
	"gopkg.in/juju/charmrepo.v3"
)

// NewDeployCommandForTest returns a command to deploy applications intended to be used only in tests.
func NewDeployCommandForTest(newAPIRoot func() (DeployAPI, error), steps []DeployStep) modelcmd.ModelCommand {
	deployCmd := &DeployCommand{
		Steps:      steps,
		NewAPIRoot: newAPIRoot,
	}
	if newAPIRoot == nil {
		deployCmd.NewAPIRoot = func() (DeployAPI, error) {
			apiRoot, err := deployCmd.ModelCommandBase.NewAPIRoot()
			if err != nil {
				return nil, errors.Trace(err)
			}
			bakeryClient, err := deployCmd.BakeryClient()
			if err != nil {
				return nil, errors.Trace(err)
			}
			controllerAPIRoot, err := deployCmd.NewControllerAPIRoot()
			if err != nil {
				return nil, errors.Trace(err)
			}
			csURL, err := getCharmStoreAPIURL(controllerAPIRoot)
			if err != nil {
				return nil, errors.Trace(err)
			}
			mURL, err := deployCmd.getMeteringAPIURL(controllerAPIRoot)
			if err != nil {
				return nil, errors.Trace(err)
			}
			cstoreClient := newCharmStoreClient(bakeryClient, csURL).WithChannel(deployCmd.Channel)

			return &deployAPIAdapter{
				Connection:        apiRoot,
				apiClient:         &apiClient{Client: apiRoot.Client()},
				charmsClient:      &charmsClient{Client: apicharms.NewClient(apiRoot)},
				applicationClient: &applicationClient{Client: application.NewClient(apiRoot)},
				modelConfigClient: &modelConfigClient{Client: modelconfig.NewClient(apiRoot)},
				charmstoreClient:  &charmstoreClient{&charmstoreClientShim{cstoreClient}},
				annotationsClient: &annotationsClient{Client: annotations.NewClient(apiRoot)},
				charmRepoClient:   &charmRepoClient{charmrepo.NewCharmStoreFromClient(cstoreClient)},
				plansClient:       &plansClient{planURL: mURL},
			}, nil
		}
	}
	return modelcmd.Wrap(deployCmd)
}

//type CharmAdder interface {
//	AddLocalCharm(*charm.URL, charm.Charm, bool) (*charm.URL, error)
//	AddCharm(*charm.URL, params.Channel, bool) error
//	AddCharmWithAuthorization(*charm.URL, params.Channel, *macaroon.Macaroon, bool) error
//	AuthorizeCharmstoreEntity(*charm.URL) (*macaroon.Macaroon, error)
//}

// noopCharmAdder for the purposes of testing the deploy command,
// we'll assume that charms/bundles are added
type noopCharmAdder struct {
	*deployAPIAdapter
}

//func (c *noopCharmAdder) AddCharm(id *charm.URL, channel params.Channel, force bool) error {
//	logger.Infof("AddCharm(id: %v, ... (arg ... what should we do!!?))", id)
//	return nil
//}
//
//func (c *noopCharmAdder) AddLocalCharm(id *charm.URL, charmData charm.Charm, force bool) (*charm.URL, error) {
//	logger.Infof("AddLocalCharm(id: %v, ...)", id)
//	return id, nil
//}
//
//func (c *noopCharmAdder) AddCharmWithAuthorization(id *charm.URL, channel params.Channel, auth *macaroon.Macaroon,force bool) error {
//	logger.Infof("AddCharmWithAuthorization(id: %v, channel: %s)", id, channel)
//	return nil
//}
//
//func (c *noopCharmAdder) AuthorizeCharmstoreEntity(id *charm.URL) (*macaroon.Macaroon, error) {
//	logger.Infof("AuthorizeCharmstoreEntity(id: %v)")
//	return (*macaroon.Macaroon)(nil), nil
//}


// NewDeployCommandForTest2 returns a command to deploy applications intended to be used only in tests
// that do not use gomock.
func NewDeployCommandForTest2(charmstore charmstoreForDeploy, charmrepo *charmstore.Repository) modelcmd.ModelCommand {
	deployCmd := &DeployCommand{
		// TOOD (tsm) allow []DeployStep to be a (varargs?) parameter
		Steps: []DeployStep{
			&RegisterMeteredCharm{
				PlanURL:      romulus.DefaultAPIRoot,
				RegisterPath: "/plan/authorize",
				QueryPath:    "/charm",
			},
			&ValidateLXDProfileCharm{},
		},
	}

	deployCmd.NewAPIRoot = func() (DeployAPI, error) {
		if charmstore == nil {
			return nil, errors.NotValidf("charmstore argument must be supplied")
		}

		if charmrepo == nil {
			return nil, errors.NotValidf("charmrepo argument must be supplied")
		}

		apiRoot, err := deployCmd.ModelCommandBase.NewAPIRoot()
		if err != nil {
			return nil, errors.Trace(err)
		}
		controllerAPIRoot, err := deployCmd.NewControllerAPIRoot()
		if err != nil {
			return nil, errors.Trace(err)
		}
		mURL, err := deployCmd.getMeteringAPIURL(controllerAPIRoot)
		if err != nil {
			return nil, errors.Trace(err)
		}
		charmstore := charmstore.WithChannel(deployCmd.Channel)

		// ch := testcharms.Repo.CharmArchive("/tmp", "nope")
		// curl := charm.MustParseURL(
		// 	fmt.Sprintf("local:quantal/%s-%d", ch.Meta().Name, ch.Revision()),
		// )
		// api.AddCharm(info)
		// api.State().AddCharm(info)
		//}

		deployAdapter := &deployAPIAdapter{
			Connection:        apiRoot,
			apiClient:         &apiClient{Client: apiRoot.Client()},
			charmsClient:      &charmsClient{Client: apicharms.NewClient(apiRoot)},
			applicationClient: &applicationClient{Client: application.NewClient(apiRoot)},
			modelConfigClient: &modelConfigClient{Client: modelconfig.NewClient(apiRoot)},
			charmstoreClient:  &charmstoreClient{charmstore},
			annotationsClient: &annotationsClient{Client: annotations.NewClient(apiRoot)},
			charmRepoClient:   &charmRepoClient{charmrepo},
			plansClient:       &plansClient{planURL: mURL},
		}
		return &noopCharmAdder{deployAdapter}, nil
	}

	return modelcmd.Wrap(deployCmd)
}

func NewUpgradeCharmCommandForTest(
	store jujuclient.ClientStore,
	apiOpen api.OpenFunc,
	deployResources resourceadapters.DeployResourcesFunc,
	resolveCharm ResolveCharmFunc,
	newCharmAdder NewCharmAdderFunc,
	newCharmClient func(base.APICallCloser) CharmClient,
	newCharmUpgradeClient func(base.APICallCloser) CharmAPIClient,
	newModelConfigGetter func(base.APICallCloser) ModelConfigGetter,
	newResourceLister func(base.APICallCloser) (ResourceLister, error),
	charmStoreURLGetter func(base.APICallCloser) (string, error),
) cmd.Command {
	cmd := &upgradeCharmCommand{
		DeployResources:       deployResources,
		ResolveCharm:          resolveCharm,
		NewCharmAdder:         newCharmAdder,
		NewCharmClient:        newCharmClient,
		NewCharmUpgradeClient: newCharmUpgradeClient,
		NewModelConfigGetter:  newModelConfigGetter,
		NewResourceLister:     newResourceLister,
		CharmStoreURLGetter:   charmStoreURLGetter,
	}
	cmd.SetClientStore(store)
	cmd.SetAPIOpen(apiOpen)
	return modelcmd.Wrap(cmd)
}

// NewResolvedCommandForTest returns a ResolvedCommand with the api provided as specified.
func NewResolvedCommandForTest(applicationResolveAPI applicationResolveAPI, clientAPI clientAPI, store jujuclient.ClientStore) modelcmd.ModelCommand {
	cmd := &resolvedCommand{applicationResolveAPI: applicationResolveAPI, clientAPI: clientAPI}
	cmd.SetClientStore(store)
	return modelcmd.Wrap(cmd)
}

// NewAddUnitCommandForTest returns an AddUnitCommand with the api provided as specified.
func NewAddUnitCommandForTest(api applicationAddUnitAPI, store jujuclient.ClientStore) modelcmd.ModelCommand {
	cmd := &addUnitCommand{api: api}
	cmd.SetClientStore(store)
	return modelcmd.Wrap(cmd)
}

// NewAddUnitCommandForTest returns an AddUnitCommand with the api provided as specified as well as overrides the refresh function.
func NewAddUnitCommandForTestWithRefresh(api applicationAddUnitAPI, store jujuclient.ClientStore, refreshFunc func(jujuclient.ClientStore, string) error) modelcmd.ModelCommand {
	cmd := &addUnitCommand{api: api}
	cmd.SetClientStore(store)
	cmd.SetModelRefresh(refreshFunc)
	return modelcmd.Wrap(cmd)
}

// NewRemoveUnitCommandForTest returns a RemoveUnitCommand with the api provided as specified.
func NewRemoveUnitCommandForTest(api RemoveApplicationAPI, store jujuclient.ClientStore) modelcmd.ModelCommand {
	cmd := &removeUnitCommand{api: api}
	cmd.SetClientStore(store)
	return modelcmd.Wrap(cmd)
}

type removeAPIFunc func() (RemoveApplicationAPI, int, error)

// NewRemoveApplicationCommandForTest returns a RemoveApplicationCommand.
func NewRemoveApplicationCommandForTest(f removeAPIFunc, store jujuclient.ClientStore) modelcmd.ModelCommand {
	c := &removeApplicationCommand{}
	c.newAPIFunc = f
	c.SetClientStore(store)
	return modelcmd.Wrap(c)
}

// NewAddRelationCommandForTest returns an AddRelationCommand with the api provided as specified.
func NewAddRelationCommandForTest(addAPI applicationAddRelationAPI, consumeAPI applicationConsumeDetailsAPI) modelcmd.ModelCommand {
	cmd := &addRelationCommand{addRelationAPI: addAPI, consumeDetailsAPI: consumeAPI}
	return modelcmd.Wrap(cmd)
}

// NewRemoveRelationCommandForTest returns an RemoveRelationCommand with the api provided as specified.
func NewRemoveRelationCommandForTest(api ApplicationDestroyRelationAPI, store jujuclient.ClientStore) modelcmd.ModelCommand {
	cmd := &removeRelationCommand{newAPIFunc: func() (ApplicationDestroyRelationAPI, error) {
		return api, nil
	}}
	cmd.SetClientStore(store)
	return modelcmd.Wrap(cmd)
}

// NewConsumeCommandForTest returns a ConsumeCommand with the specified api.
func NewConsumeCommandForTest(
	store jujuclient.ClientStore,
	sourceAPI applicationConsumeDetailsAPI,
	targetAPI applicationConsumeAPI,
) cmd.Command {
	c := &consumeCommand{sourceAPI: sourceAPI, targetAPI: targetAPI}
	c.SetClientStore(store)
	return modelcmd.Wrap(c)
}

// NewSetSeriesCommandForTest returns a SetSeriesCommand with the specified api.
func NewSetSeriesCommandForTest(
	seriesAPI setSeriesAPI,
	store jujuclient.ClientStore,
) modelcmd.ModelCommand {
	cmd := &setSeriesCommand{
		setSeriesClient: seriesAPI,
	}
	cmd.SetClientStore(store)
	return modelcmd.Wrap(cmd)
}

// NewSuspendRelationCommandForTest returns a SuspendRelationCommand with the api provided as specified.
func NewSuspendRelationCommandForTest(api SetRelationSuspendedAPI, store jujuclient.ClientStore) modelcmd.ModelCommand {
	cmd := &suspendRelationCommand{newAPIFunc: func() (SetRelationSuspendedAPI, error) {
		return api, nil
	}}
	cmd.SetClientStore(store)
	return modelcmd.Wrap(cmd)
}

// NewResumeRelationCommandForTest returns a ResumeRelationCommand with the api provided as specified.
func NewResumeRelationCommandForTest(api SetRelationSuspendedAPI, store jujuclient.ClientStore) modelcmd.ModelCommand {
	cmd := &resumeRelationCommand{newAPIFunc: func() (SetRelationSuspendedAPI, error) {
		return api, nil
	}}
	cmd.SetClientStore(store)
	return modelcmd.Wrap(cmd)
}

// NewRemoveSaasCommandForTest returns a RemoveSaasCommand with the api provided as specified.
func NewRemoveSaasCommandForTest(api RemoveSaasAPI, store jujuclient.ClientStore) modelcmd.ModelCommand {
	cmd := &removeSaasCommand{newAPIFunc: func() (RemoveSaasAPI, error) {
		return api, nil
	}}
	cmd.SetClientStore(store)
	return modelcmd.Wrap(cmd)
}

// NewScaleCommandForTest returns a ScaleCommand with the api provided as specified.
func NewScaleCommandForTest(api scaleApplicationAPI, store jujuclient.ClientStore) modelcmd.ModelCommand {
	cmd := &scaleApplicationCommand{newAPIFunc: func() (scaleApplicationAPI, error) {
		return api, nil
	}}
	cmd.SetClientStore(store)
	return modelcmd.Wrap(cmd)
}

func NewBundleDiffCommandForTest(api base.APICallCloser, charmStore BundleResolver, store jujuclient.ClientStore) modelcmd.ModelCommand {
	cmd := &bundleDiffCommand{
		_apiRoot:    api,
		_charmStore: charmStore,
	}
	cmd.SetClientStore(store)
	return modelcmd.Wrap(cmd)
}

func NewShowCommandForTest(api ApplicationsInfoAPI, store jujuclient.ClientStore) cmd.Command {
	cmd := &showApplicationCommand{newAPIFunc: func() (ApplicationsInfoAPI, error) {
		return api, nil
	}}
	cmd.SetClientStore(store)
	return modelcmd.Wrap(cmd)
}
