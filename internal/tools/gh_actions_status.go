package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-git/internal/git"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// GHActionsStatusSchema returns the JSON Schema for the gh_actions_status tool.
func GHActionsStatusSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Working directory path (defaults to current directory)",
			},
			"workflow": map[string]any{
				"type":        "string",
				"description": "Filter by workflow name or filename",
			},
			"branch": map[string]any{
				"type":        "string",
				"description": "Filter by branch name",
			},
		},
	})
	return s
}

// GHActionsStatus returns a tool handler that shows CI/CD workflow status.
func GHActionsStatus() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		path := helpers.GetString(req.Arguments, "path")
		workflow := helpers.GetString(req.Arguments, "workflow")
		branch := helpers.GetString(req.Arguments, "branch")

		args := []string{"run", "list",
			"--json", "databaseId,displayTitle,status,conclusion,url,headBranch,createdAt,workflowName",
			"--limit", "10",
		}
		if workflow != "" {
			args = append(args, "--workflow", workflow)
		}
		if branch != "" {
			args = append(args, "--branch", branch)
		}

		output, err := git.GH(ctx, path, args...)
		if err != nil {
			return helpers.ErrorResult("gh_error", err.Error()), nil
		}
		if output == "" || output == "[]" {
			return helpers.TextResult("No workflow runs found"), nil
		}
		return helpers.TextResult(output), nil
	}
}
