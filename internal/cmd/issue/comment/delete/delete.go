package delete

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/PrabhashDiss/jira-cli/api"
	"github.com/PrabhashDiss/jira-cli/internal/cmdutil"
)

const (
	helpText = `Delete removes a comment from an issue. It asks for confirmation unless '--force' is passed.`
	examples = `$ jira issue comment delete ISSUE-1 10001

# Skip the confirmation prompt
$ jira issue comment delete ISSUE-1 10001 --force`
)

// NewCmdCommentDelete is a comment delete command.
func NewCmdCommentDelete() *cobra.Command {
	cmd := cobra.Command{
		Use:     "delete ISSUE-KEY COMMENT-ID",
		Short:   "Delete a comment from an issue",
		Long:    helpText,
		Example: examples,
		Aliases: []string{"remove", "rm", "del"},
		Annotations: map[string]string{
			"help:args": "ISSUE-KEY\tIssue key of the source issue, eg: ISSUE-1\n" +
				"COMMENT-ID\tID of the comment to delete, eg: 10001",
		},
		Args: cobra.ExactArgs(2),
		Run:  del,
	}

	cmd.Flags().Bool("force", false, "Delete the comment without asking for confirmation")

	return &cmd
}

func del(cmd *cobra.Command, args []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	force, err := cmd.Flags().GetBool("force")
	cmdutil.ExitIfError(err)

	issueKey := cmdutil.GetJiraIssueKey(viper.GetString("project.key"), args[0])
	commentID := args[1]

	if !force {
		var ans bool
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Delete comment %q on issue %q?", commentID, issueKey),
		}
		cmdutil.ExitIfError(survey.AskOne(prompt, &ans))
		if !ans {
			cmdutil.Failed("Action aborted")
		}
	}

	client := api.DefaultClient(debug)

	err = func() error {
		s := cmdutil.Info("Deleting comment")
		defer s.Stop()

		return client.DeleteIssueComment(issueKey, commentID)
	}()
	cmdutil.ExitIfError(err)

	cmdutil.Success("Comment %q deleted from issue %q", commentID, issueKey)
}
