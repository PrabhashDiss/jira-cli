package comment

import (
	"github.com/spf13/cobra"

	"github.com/PrabhashDiss/jira-cli/internal/cmd/issue/comment/add"
	"github.com/PrabhashDiss/jira-cli/internal/cmd/issue/comment/delete"
	"github.com/PrabhashDiss/jira-cli/internal/cmd/issue/comment/edit"
)

const helpText = `Comment command helps you manage issue comments. See available commands below.`

// NewCmdComment is a comment command.
func NewCmdComment() *cobra.Command {
	cmd := cobra.Command{
		Use:     "comment",
		Short:   "Manage issue comments",
		Long:    helpText,
		Aliases: []string{"comments"},
		RunE:    comment,
	}

	cmd.AddCommand(add.NewCmdCommentAdd())
	cmd.AddCommand(edit.NewCmdCommentEdit())
	cmd.AddCommand(delete.NewCmdCommentDelete())

	return &cmd
}

func comment(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}
