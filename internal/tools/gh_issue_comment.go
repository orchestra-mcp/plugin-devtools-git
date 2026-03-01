package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-git/internal/git"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// GHIssueCommentSchema returns the JSON Schema for the gh_issue_comment tool.
func GHIssueCommentSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"issue_number": map[string]any{
				"type":        "number",
				"description": "Issue number to comment on",
			},
			"body": map[string]any{
				"type":        "string",
				"description": "Comment body",
			},
			"path": map[string]any{
				"type":        "string",
				"description": "Working directory path (defaults to current directory)",
			},
		},
		"required": []any{"issue_number", "body"},
	})
	return s
}

// GHIssueComment returns a tool handler that comments on a GitHub issue.
func GHIssueComment() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "body"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		issueNumber := helpers.GetInt(req.Arguments, "issue_number")
		if issueNumber == 0 {
			return helpers.ErrorResult("validation_error", "issue_number is required"), nil
		}

		body := helpers.GetString(req.Arguments, "body")
		path := helpers.GetString(req.Arguments, "path")

		args := []string{"issue", "comment", fmt.Sprintf("%d", issueNumber), "--body", body}

		output, err := git.GH(ctx, path, args...)
		if err != nil {
			return helpers.ErrorResult("gh_error", err.Error()), nil
		}
		if output == "" {
			output = fmt.Sprintf("Comment added to issue #%d", issueNumber)
		}
		return helpers.TextResult(output), nil
	}
}
