// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package comments

import (
	stdctx "context"
	"fmt"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"

	"github.com/urfave/cli/v3"
)

// CmdCommentsDelete deletes a comment
var CmdCommentsDelete = cli.Command{
	Name:        "delete",
	Aliases:     []string{"rm"},
	Usage:       "Delete a comment",
	Description: "Delete an existing comment by ID",
	ArgsUsage:   "<comment id>",
	Action:      runCommentsDelete,
	Flags:       flags.AllDefaultFlags,
}

func runCommentsDelete(_ stdctx.Context, cmd *cli.Command) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	args := ctx.Args()
	if args.Len() == 0 {
		return fmt.Errorf("must specify comment ID")
	}

	commentID, err := parseCommentID(args.First())
	if err != nil {
		return err
	}

	client := ctx.Login.Client()
	_, err = client.DeleteIssueComment(ctx.Owner, ctx.Repo, commentID)
	if err != nil {
		return err
	}

	fmt.Printf("Comment %d deleted\n", commentID)
	return nil
}
