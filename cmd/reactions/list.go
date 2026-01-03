// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package reactions

import (
	stdctx "context"
	"fmt"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/utils"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v3"
)

// CmdReactionList lists reactions on an issue or comment
var CmdReactionList = cli.Command{
	Name:        "list",
	Aliases:     []string{"ls"},
	Usage:       "List reactions on an issue or comment",
	Description: "List all emoji reactions on an issue, PR, or comment.",
	Action:      runReactionList,
	Flags: append([]cli.Flag{
		&cli.IntFlag{
			Name:    "issue",
			Aliases: []string{"i"},
			Usage:   "Issue or PR index",
		},
		&cli.IntFlag{
			Name:    "comment",
			Aliases: []string{"c"},
			Usage:   "Comment ID (if listing reactions on a comment)",
		},
	}, flags.AllDefaultFlags...),
}

func runReactionList(_ stdctx.Context, cmd *cli.Command) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	issueIndex := cmd.Int("issue")
	commentID := cmd.Int("comment")

	if issueIndex == 0 && commentID == 0 {
		return fmt.Errorf("must specify --issue or --comment")
	}

	client := ctx.Login.Client()

	var reactions []*gitea.Reaction
	var err error

	if commentID > 0 {
		reactions, _, err = client.GetIssueCommentReactions(ctx.Owner, ctx.Repo, int64(commentID))
		if err != nil {
			return fmt.Errorf("failed to get reactions for comment: %w", err)
		}
		fmt.Printf("Reactions on comment %d:\n", commentID)
	} else {
		index, err2 := utils.ArgToIndex(fmt.Sprintf("%d", issueIndex))
		if err2 != nil {
			return err2
		}
		reactions, _, err = client.GetIssueReactions(ctx.Owner, ctx.Repo, index)
		if err != nil {
			return fmt.Errorf("failed to get reactions for issue: %w", err)
		}
		fmt.Printf("Reactions on issue #%d:\n", issueIndex)
	}

	if len(reactions) == 0 {
		fmt.Println("  No reactions")
		return nil
	}

	// Group reactions by type
	reactionCounts := make(map[string][]string)
	for _, r := range reactions {
		reactionCounts[r.Reaction] = append(reactionCounts[r.Reaction], r.User.UserName)
	}

	for reaction, users := range reactionCounts {
		fmt.Printf("  %s (%d): %v\n", reactionToEmoji(reaction), len(users), users)
	}

	return nil
}

// reactionToEmoji converts reaction names to emoji
func reactionToEmoji(reaction string) string {
	switch reaction {
	case "+1":
		return "+1 (thumbs up)"
	case "-1":
		return "-1 (thumbs down)"
	case "laugh":
		return "laugh"
	case "confused":
		return "confused"
	case "heart":
		return "heart"
	case "hooray":
		return "hooray (tada)"
	case "rocket":
		return "rocket"
	case "eyes":
		return "eyes"
	default:
		return reaction
	}
}
