// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package version

import (
	"fmt"
	"runtime"

	jc "github.com/juju/testing/checkers"
	semversion "github.com/juju/version"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/testing"
)

type suite struct {
	testing.BaseSuite
}

var _ = gc.Suite(&suite{})

var isDevTests = []struct {
	num semversion.Number
	dev bool
}{{
	num: semversion.Number{},
}, {
	num: semversion.Number{Major: 0, Minor: 0, Patch: 1},
}, {
	num: semversion.Number{Major: 0, Minor: 0, Patch: 2},
}, {
	num: semversion.Number{Major: 0, Minor: 1, Patch: 0},
	dev: true,
}, {
	num: semversion.Number{Major: 0, Minor: 2, Patch: 3},
}, {
	num: semversion.Number{Major: 1, Minor: 0, Patch: 0},
}, {
	num: semversion.Number{Major: 10, Minor: 234, Patch: 3456},
}, {
	num: semversion.Number{Major: 10, Minor: 234, Patch: 3456, Build: 1},
	dev: true,
}, {
	num: semversion.Number{Major: 10, Minor: 234, Patch: 3456, Build: 64},
	dev: true,
}, {
	num: semversion.Number{Major: 10, Minor: 235, Patch: 3456},
}, {
	num: semversion.Number{Major: 1, Minor: 21, Patch: 1, Tag: "alpha"},
	dev: true,
}, {
	num: semversion.Number{Major: 1, Minor: 21, Patch: 1, Tag: "alpha", Build: 1},
	dev: true,
}, {
	num: semversion.Number{Major: 1, Minor: 21},
}}

func (*suite) TestIsDev(c *gc.C) {
	for i, test := range isDevTests {
		c.Logf("test %d", i)
		c.Check(IsDev(test.num), gc.Equals, test.dev)
	}
}

func (s *suite) TestCompiler(c *gc.C) {
	c.Assert(Compiler, gc.Equals, runtime.Compiler)
}

func (s *suite) TestCheckJujuMinVersion(c *gc.C) {
	for _, test := range []struct {
		toCheck     semversion.Number
		jujuVersion semversion.Number
		error       bool
	}{
		{
			toCheck:     semversion.Zero,
			jujuVersion: semversion.MustParse("2.8.0"),
			error:       false,
		}, {
			toCheck:     semversion.MustParse("2.8.0"),
			jujuVersion: semversion.MustParse("2.8.0"),
			error:       false,
		}, {
			toCheck:     semversion.MustParse("2.8.0"),
			jujuVersion: semversion.MustParse("2.8.1"),
			error:       false,
		}, {
			toCheck:     semversion.MustParse("2.8.0"),
			jujuVersion: semversion.MustParse("2.9.0"),
			error:       false,
		}, {
			toCheck:     semversion.MustParse("2.8.0"),
			jujuVersion: semversion.MustParse("3.0.0"),
			error:       false,
		}, {
			toCheck:     semversion.MustParse("2.8.0"),
			jujuVersion: semversion.MustParse("2.8-beta1"),
			error:       false,
		}, {
			toCheck:     semversion.MustParse("2.8.0"),
			jujuVersion: semversion.MustParse("2.7.0"),
			error:       true,
		},
	} {
		err := CheckJujuMinVersion(test.toCheck, test.jujuVersion)
		if test.error {
			c.Assert(err, jc.Satisfies, IsMinVersionError)
			c.Assert(err.Error(), gc.Equals,
				fmt.Sprintf("charm's min version (%s) is higher than this juju model's version (%s)",
					test.toCheck, test.jujuVersion))
		} else {
			c.Assert(err, jc.ErrorIsNil)
		}
	}
}
