package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-git/internal/git"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// GitStatusSchema returns the JSON Schema for the git_status tool.
func GitStatusSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Working directory path (defaults to current directory)",
			},
		},
	})
	return s
}

// GitStatus returns a tool handler that shows working tree status.
func GitStatus() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		path := helpers.GetString(req.Arguments, "path")

		output, err := git.Run(ctx, path, "status", "--porcelain", "-b")
		if err != nil {
			return helpers.ErrorResult("git_error", err.Error()), nil
		}
		if output == "" {
			output = "Working tree clean"
		}
		return helpers.TextResult(output), nil
	}
}
