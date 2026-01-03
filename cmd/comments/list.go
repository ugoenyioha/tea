// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package comments

import (
	stdctx "context"
	"fmt"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"
	"code.gitea.io/tea/modules/utils"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v3"
)

// CmdCommentsList lists comments on an issue or PR
var CmdCommentsList = cli.Command{
	Name:        "list",
	Aliases:     []string{"ls"},
	Usage:       "List comments on an issue or pull request",
	Description: "List all comments on an issue or pull request",
	ArgsUsage:   "<issue/pr index>",
	Action:      runCommentsList,
	Flags: append([]cli.Flag{
		&cli.IntFlag{
			Name:    "limit",
			Aliases: []string{"lm"},
			Usage:   "Limit number of comments to return",
			Value:   0,
		},
	}, flags.AllDefaultFlags...),
}

func runCommentsList(_ stdctx.Context, cmd *cli.Command) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	if !cmd.Args().Present() {
		return fmt.Errorf("must specify issue/pr index")
	}

	index, err := utils.ArgToIndex(cmd.Args().First())
	if err != nil {
		return err
	}

	client := ctx.Login.Client()

	opt := gitea.ListIssueCommentOptions{}
	if limit := cmd.Int("limit"); limit > 0 {
		opt.PageSize = int(limit)
	}

	comments, _, err := client.ListIssueComments(ctx.Owner, ctx.Repo, index, opt)
	if err != nil {
		return err
	}

	if len(comments) == 0 {
		fmt.Println("No comments found")
		return nil
	}

	print.Comments(comments)
	return nil
}
