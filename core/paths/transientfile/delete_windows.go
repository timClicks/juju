// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

//+build windows

// Package transientfile provides helpers for creating files that do not
// survive machine reboots.
package transientfile

import (
	"github.com/juju/errors"
	"golang.org/x/sys/windows"
)

// ensureDeleteAfterReboot arranges for the specified file to be removed once
// the machine reboots. It exploits the MoveFileEx API call which allows us to
// defer the deletion of a file by specifying a nil value as the destination
// target in conjunction with the MOVEFILE_DELAY_UNTIL_REBOOT flag.
//
// The file to be deleted is appended to the windows registry and automatically
// cleaned up by windows when the system reboots.
//
// For more info see: https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-movefileexa
func ensureDeleteAfterReboot(file string) error {
	from, err := windows.UTF16PtrFromString(file)
	if err != nil {
		return errors.Trace(err)
	}
	return windows.MoveFileEx(from, nil, windows.MOVEFILE_DELAY_UNTIL_REBOOT)
}
