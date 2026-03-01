package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-git/internal/git"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// GitCheckoutSchema returns the JSON Schema for the git_checkout tool.
func GitCheckoutSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"ref": map[string]any{
				"type":        "string",
				"description": "Branch, tag, or commit to checkout",
			},
			"path": map[string]any{
				"type":        "string",
				"description": "Working directory path (defaults to current directory)",
			},
			"files": map[string]any{
				"type":        "array",
				"description": "Specific files to checkout from the ref",
				"items": map[string]any{
					"type": "string",
				},
			},
		},
		"required": []any{"ref"},
	})
	return s
}

// GitCheckout returns a tool handler that checks out files or refs.
func GitCheckout() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "ref"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		ref := helpers.GetString(req.Arguments, "ref")
		path := helpers.GetString(req.Arguments, "path")
		files := helpers.GetStringSlice(req.Arguments, "files")

		args := []string{"checkout", ref}

		// If specific files are given, use -- to separate.
		if len(files) > 0 {
			args = append(args, "--")
			args = append(args, files...)
		}

		output, err := git.Run(ctx, path, args...)
		if err != nil {
			return helpers.ErrorResult("git_error", err.Error()), nil
		}
		if output == "" {
			output = "Checkout successful"
		}
		return helpers.TextResult(output), nil
	}
}
