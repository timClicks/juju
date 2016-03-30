// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package user_test

import (
	"errors"
	"strings"

	"github.com/juju/cmd"
	"github.com/juju/names"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
	"gopkg.in/macaroon.v1"

	"github.com/juju/juju/cmd/juju/user"
	"github.com/juju/juju/juju"
	"github.com/juju/juju/jujuclient"
	coretesting "github.com/juju/juju/testing"
)

type LoginCommandSuite struct {
	BaseSuite
	mockAPI *mockLoginAPI
}

var _ = gc.Suite(&LoginCommandSuite{})

func (s *LoginCommandSuite) SetUpTest(c *gc.C) {
	s.BaseSuite.SetUpTest(c)
	s.mockAPI = &mockLoginAPI{}
}

func (s *LoginCommandSuite) run(c *gc.C, args ...string) (*cmd.Context, juju.NewAPIConnectionParams, error) {
	var argsOut juju.NewAPIConnectionParams
	cmd, _ := user.NewLoginCommandForTest(func(args juju.NewAPIConnectionParams) (user.LoginAPI, error) {
		argsOut = args
		// The account details are modified in place, so take a copy.
		accountDetails := *argsOut.AccountDetails
		argsOut.AccountDetails = &accountDetails
		return s.mockAPI, nil
	}, s.store)
	ctx := coretesting.Context(c)
	ctx.Stdin = strings.NewReader("sekrit\nsekrit\n")
	err := coretesting.InitCommand(cmd, args)
	if err != nil {
		return nil, argsOut, err
	}
	err = cmd.Run(ctx)
	return ctx, argsOut, err
}

func (s *LoginCommandSuite) TestInit(c *gc.C) {
	for i, test := range []struct {
		args        []string
		user        string
		generate    bool
		errorString string
	}{
		{
		// no args is fine
		}, {
			args:     []string{"foobar"},
			user:     "foobar",
			generate: true,
		}, {
			args:        []string{"--foobar"},
			errorString: "flag provided but not defined: --foobar",
		}, {
			args:        []string{"foobar", "extra"},
			errorString: `unrecognized args: \["extra"\]`,
		},
	} {
		c.Logf("test %d", i)
		wrappedCommand, command := user.NewLoginCommandForTest(nil, s.store)
		err := coretesting.InitCommand(wrappedCommand, test.args)
		if test.errorString == "" {
			c.Check(command.User, gc.Equals, test.user)
		} else {
			c.Check(err, gc.ErrorMatches, test.errorString)
		}
	}
}

func (s *LoginCommandSuite) assertStorePassword(c *gc.C, user, pass string) {
	details, err := s.store.AccountByName("testing", user)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(details.Password, gc.Equals, pass)
}

func (s *LoginCommandSuite) assertStoreMacaroon(c *gc.C, user string, mac *macaroon.Macaroon) {
	details, err := s.store.AccountByName("testing", user)
	c.Assert(err, jc.ErrorIsNil)
	if mac == nil {
		c.Assert(details.Macaroon, gc.Equals, "")
		return
	}
	macaroonJSON, err := mac.MarshalJSON()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(details.Macaroon, gc.Equals, string(macaroonJSON))
}

func (s *LoginCommandSuite) TestLogin(c *gc.C) {
	context, args, err := s.run(c)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(coretesting.Stdout(context), gc.Equals, "")
	c.Assert(coretesting.Stderr(context), gc.Equals, `
password: 
type password again: 
You are now logged in to "testing" as "current-user@local".
`[1:],
	)
	s.assertStorePassword(c, "current-user@local", "")
	s.assertStoreMacaroon(c, "current-user@local", fakeLocalLoginMacaroon(names.NewUserTag("current-user@local")))
	c.Assert(args.AccountDetails, jc.DeepEquals, &jujuclient.AccountDetails{
		User:     "current-user@local",
		Password: "sekrit",
	})
}

func (s *LoginCommandSuite) TestLoginNewUser(c *gc.C) {
	context, args, err := s.run(c, "new-user")
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(coretesting.Stdout(context), gc.Equals, "")
	c.Assert(coretesting.Stderr(context), gc.Equals, `
password: 
type password again: 
You are now logged in to "testing" as "new-user@local".
`[1:],
	)
	s.assertStorePassword(c, "new-user@local", "")
	s.assertStoreMacaroon(c, "new-user@local", fakeLocalLoginMacaroon(names.NewUserTag("new-user@local")))
	c.Assert(args.AccountDetails, jc.DeepEquals, &jujuclient.AccountDetails{
		User:     "new-user@local",
		Password: "sekrit",
	})
}

func (s *LoginCommandSuite) TestLoginFail(c *gc.C) {
	s.mockAPI.SetErrors(errors.New("failed to do something"))
	_, _, err := s.run(c)
	c.Assert(err, gc.ErrorMatches, "failed to create a temporary credential: failed to do something")
	s.assertStorePassword(c, "current-user@local", "old-password")
	s.assertStoreMacaroon(c, "current-user@local", nil)
}

type mockLoginAPI struct {
	mockChangePasswordAPI
}
