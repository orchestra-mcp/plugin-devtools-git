package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-git/internal/git"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// GHIssueCreateSchema returns the JSON Schema for the gh_issue_create tool.
func GHIssueCreateSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"title": map[string]any{
				"type":        "string",
				"description": "Issue title",
			},
			"body": map[string]any{
				"type":        "string",
				"description": "Issue body/description",
			},
			"path": map[string]any{
				"type":        "string",
				"description": "Working directory path (defaults to current directory)",
			},
			"labels": map[string]any{
				"type":        "array",
				"description": "Labels to apply to the issue",
				"items": map[string]any{
					"type": "string",
				},
			},
			"assignees": map[string]any{
				"type":        "array",
				"description": "Users to assign to the issue",
				"items": map[string]any{
					"type": "string",
				},
			},
		},
		"required": []any{"title", "body"},
	})
	return s
}

// GHIssueCreate returns a tool handler that creates a GitHub issue.
func GHIssueCreate() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "title", "body"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		title := helpers.GetString(req.Arguments, "title")
		body := helpers.GetString(req.Arguments, "body")
		path := helpers.GetString(req.Arguments, "path")
		labels := helpers.GetStringSlice(req.Arguments, "labels")
		assignees := helpers.GetStringSlice(req.Arguments, "assignees")

		args := []string{"issue", "create", "--title", title, "--body", body}

		for _, label := range labels {
			args = append(args, "--label", label)
		}
		for _, assignee := range assignees {
			args = append(args, "--assignee", assignee)
		}

		output, err := git.GH(ctx, path, args...)
		if err != nil {
			return helpers.ErrorResult("gh_error", err.Error()), nil
		}
		return helpers.TextResult(output), nil
	}
}
