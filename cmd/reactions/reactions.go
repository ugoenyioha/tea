// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package reactions

import (
	"code.gitea.io/tea/cmd/flags"

	"github.com/urfave/cli/v3"
)

// CmdReactions is the main command for managing reactions
var CmdReactions = cli.Command{
	Name:        "reaction",
	Aliases:     []string{"reactions", "react"},
	Usage:       "Manage reactions on issues and comments",
	Description: "Add, remove, or list reactions (emoji) on issues, PRs, and comments",
	Flags:       flags.AllDefaultFlags,
	Commands: []*cli.Command{
		&CmdReactionAdd,
		&CmdReactionRemove,
		&CmdReactionList,
	},
}

// Valid reactions in Gitea
var ValidReactions = []string{
	"+1", "-1", "laugh", "confused", "heart", "hooray", "rocket", "eyes",
}

// ReactionHelp provides help text for valid reactions
const ReactionHelp = `Valid reactions: +1, -1, laugh, confused, heart, hooray, rocket, eyes
Aliases: thumbsup (+1), thumbsdown (-1), tada (hooray)`

// NormalizeReaction converts common aliases to their canonical form
func NormalizeReaction(reaction string) string {
	switch reaction {
	case "thumbsup", ":+1:", ":thumbsup:":
		return "+1"
	case "thumbsdown", ":-1:", ":thumbsdown:":
		return "-1"
	case "tada", ":tada:", ":hooray:":
		return "hooray"
	case ":laugh:", ":laughing:":
		return "laugh"
	case ":confused:", ":thinking:":
		return "confused"
	case ":heart:", "<3":
		return "heart"
	case ":rocket:":
		return "rocket"
	case ":eyes:":
		return "eyes"
	default:
		return reaction
	}
}
