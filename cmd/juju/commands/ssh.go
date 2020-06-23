// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package commands

import (
	"fmt"
	"strconv"

	"github.com/juju/cmd"
	"github.com/juju/errors"
	"github.com/juju/gnuflag"
	"github.com/juju/utils/ssh"

	jujucmd "github.com/juju/juju/cmd"
	"github.com/juju/juju/cmd/modelcmd"
	jujussh "github.com/juju/juju/network/ssh"
)

var usageSSHSummary = `
Initiate a secure shell (SSH) session or execute a command on a Juju machine.`[1:]

var usageSSHDetails = `
This command operates in two modes, like its namesake, the ssh command. It will
either connect to a machine and establish an interactive secure shell (SSH) 
session or execute a command on that machine and exit.

Basic usage involves specifying a target to connect to and optionally a command
to execute. 

	juju ssh (<machine>|<unit>)
	juju ssh (<machine>|<unit>) <command>

'juju ssh' is very powerful,  but mistakes can disrupt the ability for Juju to
operate in a fully automated manner in the future. Where possible, prefer pre-
defined actions to avoid this problem. See the Further Reading section for 
information about actions.

When attempting to execute commands within a "unit context", e.g. to use hook 
tools manually, prefer the 'juju exec' command.  


Establishing a secure shell 


Specify a target machine to connect to via its machine ID or its unit name.

	juju ssh <target>

<target> specifies where Juju should establish the connection. It takes 
the form:

    [<user>@](<machine>|<unit>)

<user> is a user account that exists on the host machine. By default, Juju 
uses the "ubuntu" user.

<machine> is a machine ID. Currently accessible machine IDs are available
via 'juju machines'.

<unit> is a unit name. Currently accessible unit names are available via
'juju status'. Unit names have the form:

    <application>/<unit-id>

<application> is the name of an application that has been added to the model.
The list of applications is available via 'juju status'.

<unit-id> is an integer. The list of units is available via 'juju status'.


Executing a command

The optional command is executed on the remote machine,  and any output is sent
back to the user. If no command is specified, then an interactive shell session
will be initiated.

	juju ssh <target> <command> [<command-argument> [...]]
	
<command> is the command to execute at <target>. Valid commands depend on the 
host's operating system. 

<command-argument> is a command-line argument. Valid arguments are command-
specific.


Advanced usage

This section outlines options for changing the command's default behaviour.


Advanced usage: pseudo-terminal behaviour

When 'juju ssh' is executed without a terminal attached, such as when piping the
output of another command into it, then the default behaviour is to not allocate
a pseudo-terminal (PTY) for the SSH session. This behaviour can be overridden by 
explicitly specifying the behaviour with '--pty=true' or '--pty=false' before 
specifying <target>.

    juju ssh --pty=(true|false) <target>


Advanced usage: disable host key verification

The SSH host keys of the target are verified to prevent man-in-the-middle (MITM)
attacks. The '--no-host-key-checks' option can be used to disable these checks.
Using this option is not recommended.

    juju ssh --no-host-key-keys=true <target>


Advanced usage: identity management

By default, Juju will select a private key (also known as an identity file) to 
connect with from the following locations:

  - ~/.ssh/id_rsa 
  - $XDG_DATA_HOME/juju/ssh/juju_id_rsa

To override this default, add the '-i <path-to-private-key>' option *after* the 
target machine.

    juju ssh <target> -i <path-to-private-key>


Advanced usage: providing options directly to ssh

Options can be passed to the local OpenSSH client (ssh) on platforms where it 
is available. Refer to your system's documentation for options supported by 
your environment.

    juju ssh <target> <ssh-option> [<ssh-option> [...]]


Examples:

    # Connect to machine 0:
    juju ssh 0
	
    # Connect to a mysql unit:
	juju ssh mysql/0
	
    # Run command 'uname -a' on machine 1:
    juju ssh 1 uname -a

	# Run command 'top' on the etcd/0 unit:
	juju ssh etcd/0 top
	
	# Establish a secure shell on the machine hosting the jenkins/0
	# unit as the "jenkins" user:
    juju ssh jenkins@jenkins/0

	# Run the 'echo hello' command on the machine hosting the mysql/0 unit
	# using a custom private key:
	juju ssh mysql/0 -i ~/.ssh/id_alternate echo hello
	

See also: 

	exec
	run-action
	scp
`

func newSSHCommand(
	hostChecker jujussh.ReachableChecker,
	isTerminal func(interface{}) bool,
) cmd.Command {
	c := new(sshCommand)
	c.setHostChecker(hostChecker)
	c.isTerminal = isTerminal
	return modelcmd.Wrap(c)
}

// sshCommand is responsible for launching a ssh shell on a given unit or machine.
type sshCommand struct {
	SSHCommon
	isTerminal func(interface{}) bool
	pty        autoBoolValue
}

func (c *sshCommand) SetFlags(f *gnuflag.FlagSet) {
	c.SSHCommon.SetFlags(f)
	f.Var(&c.pty, "pty", "Enable pseudoterminal (pseudo-tty or PTY) allocation")
}

func (c *sshCommand) Info() *cmd.Info {
	return jujucmd.Info(&cmd.Info{
		Name:    "ssh",
		Args:    "[<user>@][<application>/]<machine> [<ssh-option> ...] [<command> [<command-option> ...]]",
		Purpose: usageSSHSummary,
		Doc:     usageSSHDetails,
	})
}

func (c *sshCommand) Init(args []string) error {
	if len(args) == 0 {
		return errors.Errorf("no target name specified")
	}
	c.Target, c.Args = args[0], args[1:]
	return nil
}

// Run resolves c.Target to a machine, to the address of a i
// machine or unit forks ssh passing any arguments provided.
func (c *sshCommand) Run(ctx *cmd.Context) error {
	err := c.initRun()
	if err != nil {
		return errors.Trace(err)
	}
	defer c.cleanupRun()

	target, err := c.resolveTarget(c.Target)
	if err != nil {
		return err
	}

	var pty bool
	if c.pty.b != nil {
		pty = *c.pty.b
	} else {
		// Flag was not specified: create a pty
		// on the remote side iff this process
		// has a terminal.
		isTerminal := isTerminal
		if c.isTerminal != nil {
			isTerminal = c.isTerminal
		}
		pty = isTerminal(ctx.Stdin)
	}

	options, err := c.getSSHOptions(pty, target)
	if err != nil {
		return err
	}

	cmd := ssh.Command(target.userHost(), c.Args, options)
	cmd.Stdin = ctx.Stdin
	cmd.Stdout = ctx.Stdout
	cmd.Stderr = ctx.Stderr
	return cmd.Run()
}

// autoBoolValue is like gnuflag.boolValue, but remembers
// whether or not a value has been set, so its behaviour
// can be determined dynamically, during command execution.
type autoBoolValue struct {
	b *bool
}

func (b *autoBoolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	b.b = &v
	return nil
}

func (b *autoBoolValue) Get() interface{} {
	if b.b != nil {
		return *b.b
	}
	return b.b // nil
}

func (b *autoBoolValue) String() string {
	if b.b != nil {
		return fmt.Sprint(*b.b)
	}
	return "<auto>"
}

func (b *autoBoolValue) IsBoolFlag() bool { return true }
