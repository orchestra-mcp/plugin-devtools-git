package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-git/internal/git"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// GHPRListSchema returns the JSON Schema for the gh_pr_list tool.
func GHPRListSchema() *structpb.Struct {
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
			"author": map[string]any{
				"type":        "string",
				"description": "Filter by author username",
			},
			"label": map[string]any{
				"type":        "string",
				"description": "Filter by label",
			},
		},
	})
	return s
}

// GHPRList returns a tool handler that lists GitHub pull requests.
func GHPRList() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		path := helpers.GetString(req.Arguments, "path")
		state := helpers.GetStringOr(req.Arguments, "state", "open")
		author := helpers.GetString(req.Arguments, "author")
		label := helpers.GetString(req.Arguments, "label")

		args := []string{"pr", "list",
			"--json", "number,title,state,author,url,labels,createdAt,headRefName",
			"--state", state,
		}
		if author != "" {
			args = append(args, "--author", author)
		}
		if label != "" {
			args = append(args, "--label", label)
		}

		output, err := git.GH(ctx, path, args...)
		if err != nil {
			return helpers.ErrorResult("gh_error", err.Error()), nil
		}
		if output == "" || output == "[]" {
			return helpers.TextResult("No pull requests found"), nil
		}
		return helpers.TextResult(output), nil
	}
}
