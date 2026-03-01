package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-git/internal/git"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// GHRepoInfoSchema returns the JSON Schema for the gh_repo_info tool.
func GHRepoInfoSchema() *structpb.Struct {
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

// GHRepoInfo returns a tool handler that shows repository metadata.
func GHRepoInfo() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		path := helpers.GetString(req.Arguments, "path")

		args := []string{"repo", "view",
			"--json", "name,owner,description,url,defaultBranchRef,stargazerCount,forkCount,isPrivate,languages,createdAt,pushedAt",
		}

		output, err := git.GH(ctx, path, args...)
		if err != nil {
			return helpers.ErrorResult("gh_error", err.Error()), nil
		}
		return helpers.TextResult(output), nil
	}
}
