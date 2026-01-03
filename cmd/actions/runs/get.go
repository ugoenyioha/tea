// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runs

import (
	stdctx "context"
	"encoding/json"
	"fmt"
	"strconv"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"

	"github.com/urfave/cli/v3"
)

// CmdRunsGet gets details of a workflow run
var CmdRunsGet = cli.Command{
	Name:        "get",
	Aliases:     []string{"view", "show"},
	Usage:       "Get details of a workflow run",
	Description: "View detailed information about a specific workflow run",
	ArgsUsage:   "<run-id>",
	Action:      runRunsGet,
	Flags:       flags.AllDefaultFlags,
}

func runRunsGet(_ stdctx.Context, cmd *cli.Command) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	if !cmd.Args().Present() {
		return fmt.Errorf("must specify a run ID")
	}

	runID, err := strconv.ParseInt(cmd.Args().First(), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid run ID: %w", err)
	}

	run, err := getWorkflowRun(ctx.Login, ctx.Owner, ctx.Repo, runID)
	if err != nil {
		return fmt.Errorf("failed to get workflow run: %w", err)
	}

	if ctx.Output == "json" {
		data, _ := json.MarshalIndent(run, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	// Print run details
	fmt.Printf("Run #%d: %s\n", run.RunNumber, run.DisplayTitle)
	fmt.Printf("  Status:     %s\n", run.Status)
	fmt.Printf("  Conclusion: %s\n", run.Conclusion)
	fmt.Printf("  Event:      %s\n", run.Event)
	fmt.Printf("  Branch:     %s\n", run.HeadBranch)
	fmt.Printf("  Commit:     %s\n", run.HeadSha[:8])
	if run.Actor != nil {
		fmt.Printf("  Actor:      %s\n", run.Actor.UserName)
	}
	if !run.StartedAt.IsZero() {
		fmt.Printf("  Started:    %s\n", print.FormatTime(run.StartedAt, false))
	}
	if !run.CompletedAt.IsZero() {
		fmt.Printf("  Completed:  %s\n", print.FormatTime(run.CompletedAt, false))
	}
	fmt.Printf("  URL:        %s\n", run.HTMLURL)

	return nil
}
