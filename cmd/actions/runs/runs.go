// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runs

import (
	"code.gitea.io/tea/cmd/flags"

	"github.com/urfave/cli/v3"
)

// CmdActionsRuns represents the runs subcommand
var CmdActionsRuns = cli.Command{
	Name:        "runs",
	Aliases:     []string{"run", "workflow"},
	Usage:       "Manage workflow runs",
	Description: "List, view, and manage GitHub Actions-style workflow runs",
	Flags:       flags.AllDefaultFlags,
	Commands: []*cli.Command{
		&CmdRunsList,
		&CmdRunsGet,
		&CmdRunsJobs,
	},
}
