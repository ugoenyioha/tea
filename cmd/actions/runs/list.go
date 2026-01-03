// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package runs

import (
	stdctx "context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"

	"github.com/urfave/cli/v3"
)

// CmdRunsList lists workflow runs
var CmdRunsList = cli.Command{
	Name:        "list",
	Aliases:     []string{"ls"},
	Usage:       "List workflow runs",
	Description: "List workflow runs for a repository",
	Action:      runRunsList,
	Flags: append([]cli.Flag{
		&cli.IntFlag{
			Name:    "limit",
			Aliases: []string{"lm"},
			Usage:   "Limit number of runs to return",
			Value:   10,
		},
		&cli.StringFlag{
			Name:    "status",
			Aliases: []string{"s"},
			Usage:   "Filter by status (queued, in_progress, completed)",
		},
		&cli.StringFlag{
			Name:  "branch",
			Usage: "Filter by branch name",
		},
		&cli.StringFlag{
			Name:  "event",
			Usage: "Filter by event type (push, pull_request, issues, issue_comment, etc)",
		},
	}, flags.AllDefaultFlags...),
}

func runRunsList(_ stdctx.Context, cmd *cli.Command) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	// Build query parameters
	params := url.Values{}
	if limit := cmd.Int("limit"); limit > 0 {
		params.Set("limit", strconv.Itoa(int(limit)))
	}
	if status := cmd.String("status"); status != "" {
		params.Set("status", status)
	}
	if branch := cmd.String("branch"); branch != "" {
		params.Set("branch", branch)
	}
	if event := cmd.String("event"); event != "" {
		params.Set("event", event)
	}

	runList, err := getWorkflowRuns(ctx.Login, ctx.Owner, ctx.Repo, params.Encode())
	if err != nil {
		return fmt.Errorf("failed to get workflow runs: %w", err)
	}

	if len(runList.WorkflowRuns) == 0 {
		fmt.Println("No workflow runs found")
		return nil
	}

	// Print runs
	printRuns(runList.WorkflowRuns, ctx.Output)
	return nil
}

func printRuns(runs []*ActionRun, output string) {
	if output == "json" {
		data, _ := json.MarshalIndent(runs, "", "  ")
		fmt.Println(string(data))
		return
	}

	// Print table header (ID is used for get/jobs commands, # is display number)
	fmt.Printf("%-7s %-4s %-12s %-10s %-15s %-10s %-35s %s\n",
		"ID", "#", "STATUS", "CONCLUSION", "EVENT", "BRANCH", "TITLE", "STARTED")
	fmt.Println(strings.Repeat("-", 115))

	for _, run := range runs {
		started := ""
		if !run.StartedAt.IsZero() {
			started = print.FormatTime(run.StartedAt, false)
		}
		fmt.Printf("%-7d %-4d %-12s %-10s %-15s %-10s %-35s %s\n",
			run.ID,
			run.RunNumber,
			run.Status,
			run.Conclusion,
			run.Event,
			truncateString(run.HeadBranch, 10),
			truncateString(run.DisplayTitle, 35),
			started,
		)
	}
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
