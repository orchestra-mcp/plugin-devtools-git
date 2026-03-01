package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-git/internal/git"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// GHIssueListSchema returns the JSON Schema for the gh_issue_list tool.
func GHIssueListSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Working directory path (defaults to current directory)",
			},
			"state": map[string]any{
				"type":        "string",
				"description": "Filter by state: open, closed, or all (default open)",
				"enum":        []any{"open", "closed", "all"},
			},
			"label": map[string]any{
				"type":        "string",
				"description": "Filter by label",
			},
			"assignee": map[string]any{
				"type":        "string",
				"description": "Filter by assignee username",
			},
		},
	})
	return s
}

// GHIssueList returns a tool handler that lists GitHub issues.
func GHIssueList() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		path := helpers.GetString(req.Arguments, "path")
		state := helpers.GetStringOr(req.Arguments, "state", "open")
		label := helpers.GetString(req.Arguments, "label")
		assignee := helpers.GetString(req.Arguments, "assignee")

		args := []string{"issue", "list",
			"--json", "number,title,state,author,url,labels,assignees,createdAt",
			"--state", state,
		}
		if label != "" {
			args = append(args, "--label", label)
		}
		if assignee != "" {
			args = append(args, "--assignee", assignee)
		}

		output, err := git.GH(ctx, path, args...)
		if err != nil {
			return helpers.ErrorResult("gh_error", err.Error()), nil
		}
		if output == "" || output == "[]" {
			return helpers.TextResult("No issues found"), nil
		}
		return helpers.TextResult(output), nil
	}
}
