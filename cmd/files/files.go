// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package files

import (
	"code.gitea.io/tea/cmd/flags"

	"github.com/urfave/cli/v3"
)

// CmdFiles is the main command for managing repository files
var CmdFiles = cli.Command{
	Name:        "files",
	Aliases:     []string{"file", "content", "contents"},
	Usage:       "Manage repository files",
	Description: "Get, create, update, or delete files in a repository via API",
	Flags:       flags.AllDefaultFlags,
	Commands: []*cli.Command{
		&CmdFilesGet,
		&CmdFilesCreate,
		&CmdFilesUpdate,
		&CmdFilesDelete,
	},
}
