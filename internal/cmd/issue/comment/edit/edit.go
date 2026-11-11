package edit

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ankitpokhrel/jira-cli/api"
	"github.com/ankitpokhrel/jira-cli/internal/cmdutil"
)

const (
	helpText = `Edit replaces the body of an existing comment.`
	examples = `$ jira issue comment edit ISSUE-1 10001 "Updated comment"

# Load comment body from a template file
$ jira issue comment edit ISSUE-1 10001 --template /path/to/template.tmpl

# Get comment body from standard input
$ jira issue comment edit ISSUE-1 10001 --template -

# Or, use pipe to read input directly from standard input
$ echo "Comment from stdin" | jira issue comment edit ISSUE-1 10001

# Positional argument takes precedence over the template flag
$ jira issue comment edit ISSUE-1 10001 "comment from arg" --template /path/to/template.tmpl`
)

// NewCmdCommentEdit is a comment edit command.
func NewCmdCommentEdit() *cobra.Command {
	cmd := cobra.Command{
		Use:     "edit ISSUE-KEY COMMENT-ID [COMMENT_BODY]",
		Short:   "Edit a comment on an issue",
		Long:    helpText,
		Example: examples,
		Annotations: map[string]string{
			"help:args": "ISSUE-KEY\tIssue key of the source issue, eg: ISSUE-1\n" +
				"COMMENT-ID\tID of the comment to edit, eg: 10001\n" +
				"COMMENT_BODY\tNew body of the comment",
		},
		Args: cobra.RangeArgs(2, 3),
		Run:  edit,
	}

	cmd.Flags().StringP("template", "T", "", "Path to a file to read comment body from")

	return &cmd
}

func edit(cmd *cobra.Command, args []string) {
	debug, err := cmd.Flags().GetBool("debug")
	cmdutil.ExitIfError(err)

	template, err := cmd.Flags().GetString("template")
	cmdutil.ExitIfError(err)

	issueKey := cmdutil.GetJiraIssueKey(viper.GetString("project.key"), args[0])
	commentID := args[1]

	var body string
	if len(args) == 3 {
		body = args[2]
	} else {
		b, err := cmdutil.ReadFile(template)
		cmdutil.ExitIfError(err)
		body = string(b)
	}
	if body == "" {
		cmdutil.Failed("Comment body is required: pass it as an argument, with `--template`, or on standard input")
	}

	client := api.DefaultClient(debug)

	err = func() error {
		s := cmdutil.Info("Updating comment")
		defer s.Stop()

		return client.UpdateIssueComment(issueKey, commentID, body)
	}()
	cmdutil.ExitIfError(err)

	cmdutil.Success("Comment %q updated on issue %q", commentID, issueKey)
	fmt.Printf("%s\n", cmdutil.GenerateServerBrowseURL(viper.GetString("server"), issueKey))
}
