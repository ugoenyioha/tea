// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package files

import (
	stdctx "context"
	"fmt"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v3"
)

// CmdFilesDelete deletes a file from the repository
var CmdFilesDelete = cli.Command{
	Name:        "delete",
	Aliases:     []string{"rm", "remove"},
	Usage:       "Delete a file from the repository",
	Description: "Delete a file from the repository via API (commits directly)",
	ArgsUsage:   "<path>",
	Action:      runFilesDelete,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:    "message",
			Aliases: []string{"m"},
			Usage:   "Commit message",
		},
		&cli.StringFlag{
			Name:    "branch",
			Aliases: []string{"b"},
			Usage:   "Branch to delete file from",
		},
		&cli.StringFlag{
			Name:  "new-branch",
			Usage: "Create a new branch for the commit",
		},
		&cli.StringFlag{
			Name:  "sha",
			Usage: "SHA of the file to delete (auto-detected if not provided)",
		},
	}, flags.AllDefaultFlags...),
}

func runFilesDelete(_ stdctx.Context, cmd *cli.Command) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	if !cmd.Args().Present() {
		return fmt.Errorf("must specify a file path")
	}

	filePath := cmd.Args().First()
	message := cmd.String("message")
	branch := cmd.String("branch")
	newBranch := cmd.String("new-branch")
	sha := cmd.String("sha")

	if message == "" {
		message = fmt.Sprintf("Delete %s", filePath)
	}

	client := ctx.Login.Client()

	// Auto-detect SHA if not provided
	if sha == "" {
		existing, _, err := client.GetContents(ctx.Owner, ctx.Repo, branch, filePath)
		if err != nil {
			return fmt.Errorf("failed to get existing file (use --sha to provide manually): %w", err)
		}
		sha = existing.SHA
	}

	opts := gitea.DeleteFileOptions{
		FileOptions: gitea.FileOptions{
			Message:       message,
			BranchName:    branch,
			NewBranchName: newBranch,
		},
		SHA: sha,
	}

	_, err := client.DeleteFile(ctx.Owner, ctx.Repo, filePath, opts)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	fmt.Printf("Deleted %s\n", filePath)

	return nil
}
