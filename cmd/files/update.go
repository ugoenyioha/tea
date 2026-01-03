// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package files

import (
	stdctx "context"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v3"
)

// CmdFilesUpdate updates an existing file in the repository
var CmdFilesUpdate = cli.Command{
	Name:        "update",
	Aliases:     []string{"edit", "modify"},
	Usage:       "Update an existing file in the repository",
	Description: "Update an existing file in the repository via API (commits directly)",
	ArgsUsage:   "<path>",
	Action:      runFilesUpdate,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:    "message",
			Aliases: []string{"m"},
			Usage:   "Commit message",
		},
		&cli.StringFlag{
			Name:    "branch",
			Aliases: []string{"b"},
			Usage:   "Branch to update file in",
		},
		&cli.StringFlag{
			Name:  "new-branch",
			Usage: "Create a new branch for the commit",
		},
		&cli.StringFlag{
			Name:    "content",
			Aliases: []string{"c"},
			Usage:   "New file content (use - for stdin)",
		},
		&cli.StringFlag{
			Name:    "from-file",
			Aliases: []string{"f"},
			Usage:   "Read content from local file",
		},
		&cli.StringFlag{
			Name:  "sha",
			Usage: "SHA of the file to update (auto-detected if not provided)",
		},
	}, flags.AllDefaultFlags...),
}

func runFilesUpdate(_ stdctx.Context, cmd *cli.Command) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	if !cmd.Args().Present() {
		return fmt.Errorf("must specify a file path")
	}

	filePath := cmd.Args().First()
	message := cmd.String("message")
	branch := cmd.String("branch")
	newBranch := cmd.String("new-branch")
	content := cmd.String("content")
	fromFile := cmd.String("from-file")
	sha := cmd.String("sha")

	if message == "" {
		message = fmt.Sprintf("Update %s", filePath)
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

	// Get content
	var fileContent []byte
	var err error

	if fromFile != "" {
		fileContent, err = os.ReadFile(fromFile)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
	} else if content == "-" {
		fileContent, err = io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}
	} else if content != "" {
		fileContent = []byte(content)
	} else {
		return fmt.Errorf("must specify --content or --from-file")
	}

	opts := gitea.UpdateFileOptions{
		FileOptions: gitea.FileOptions{
			Message:       message,
			BranchName:    branch,
			NewBranchName: newBranch,
		},
		SHA:     sha,
		Content: base64.StdEncoding.EncodeToString(fileContent),
	}

	resp, _, err := client.UpdateFile(ctx.Owner, ctx.Repo, filePath, opts)
	if err != nil {
		return fmt.Errorf("failed to update file: %w", err)
	}

	fmt.Printf("Updated %s\n", filePath)
	if resp.Commit != nil {
		fmt.Printf("Commit: %s\n", resp.Commit.SHA)
	}

	return nil
}
