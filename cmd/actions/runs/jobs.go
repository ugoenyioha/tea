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

// CmdRunsJobs lists jobs for a workflow run
var CmdRunsJobs = cli.Command{
	Name:        "jobs",
	Aliases:     []string{"job"},
	Usage:       "List jobs for a workflow run",
	Description: "View jobs and their steps for a specific workflow run",
	ArgsUsage:   "<run-id>",
	Action:      runRunsJobs,
	Flags:       flags.AllDefaultFlags,
}

func runRunsJobs(_ stdctx.Context, cmd *cli.Command) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	if !cmd.Args().Present() {
		return fmt.Errorf("must specify a run ID")
	}

	runID, err := strconv.ParseInt(cmd.Args().First(), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid run ID: %w", err)
	}

	jobs, err := getWorkflowRunJobs(ctx.Login, ctx.Owner, ctx.Repo, runID)
	if err != nil {
		return fmt.Errorf("failed to get jobs: %w", err)
	}

	if len(jobs.Jobs) == 0 {
		fmt.Println("No jobs found for this run")
		return nil
	}

	if ctx.Output == "json" {
		data, _ := json.MarshalIndent(jobs.Jobs, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	// Print jobs with their steps
	for _, job := range jobs.Jobs {
		statusIcon := getStatusIcon(job.Status, job.Conclusion)
		fmt.Printf("%s Job: %s (%s)\n", statusIcon, job.Name, job.Status)

		if len(job.Steps) > 0 {
			for _, step := range job.Steps {
				stepIcon := getStatusIcon(step.Status, step.Conclusion)
				duration := ""
				if !step.StartedAt.IsZero() && !step.CompletedAt.IsZero() {
					duration = fmt.Sprintf(" (%s)", step.CompletedAt.Sub(step.StartedAt).Round(1e9))
				}
				fmt.Printf("  %s Step %d: %s%s\n", stepIcon, step.Number, step.Name, duration)
			}
		}

		if !job.StartedAt.IsZero() {
			fmt.Printf("  Started:   %s\n", print.FormatTime(job.StartedAt, false))
		}
		if !job.CompletedAt.IsZero() {
			fmt.Printf("  Completed: %s\n", print.FormatTime(job.CompletedAt, false))
		}
		fmt.Println()
	}

	return nil
}

func getStatusIcon(status, conclusion string) string {
	if status == "completed" {
		switch conclusion {
		case "success":
			return "[ok]"
		case "failure":
			return "[FAIL]"
		case "cancelled":
			return "[cancelled]"
		case "skipped":
			return "[skip]"
		default:
			return "[?]"
		}
	}
	switch status {
	case "queued":
		return "[queue]"
	case "in_progress":
		return "[running]"
	case "waiting":
		return "[wait]"
	default:
		return "[" + status + "]"
	}
}
