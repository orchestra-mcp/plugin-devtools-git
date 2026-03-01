package tools_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-git/internal/git"
	"github.com/orchestra-mcp/plugin-devtools-git/internal/tools"
	"google.golang.org/protobuf/types/known/structpb"
)

// ---------- Helpers ----------

func makeArgs(t *testing.T, m map[string]any) *structpb.Struct {
	t.Helper()
	s, err := structpb.NewStruct(m)
	if err != nil {
		t.Fatalf("makeArgs: %v", err)
	}
	return s
}

func initRepo(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "git-tools-test-*")
	if err != nil {
		t.Fatalf("temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	ctx := context.Background()
	for _, args := range [][]string{
		{"init"},
		{"config", "user.email", "test@example.com"},
		{"config", "user.name", "Test User"},
	} {
		if _, err := git.Run(ctx, dir, args...); err != nil {
			t.Fatalf("git %v: %v", args, err)
		}
	}
	return dir
}

func writeAndCommit(t *testing.T, dir, filename, content, msg string) {
	t.Helper()
	ctx := context.Background()
	if err := os.WriteFile(filepath.Join(dir, filename), []byte(content), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	for _, args := range [][]string{{"add", filename}, {"commit", "-m", msg}} {
		if _, err := git.Run(ctx, dir, args...); err != nil {
			t.Fatalf("git %v: %v", args, err)
		}
	}
}

func responseText(t *testing.T, resp *pluginv1.ToolResponse) string {
	t.Helper()
	if resp.Result == nil {
		return ""
	}
	v, ok := resp.Result.Fields["text"]
	if !ok {
		t.Fatalf("result missing 'text' key")
	}
	return v.GetStringValue()
}

// ---------- git_status ----------

func TestGitStatus_CleanRepo(t *testing.T) {
	dir := initRepo(t)
	writeAndCommit(t, dir, "README.md", "# Test\n", "init")

	fn := tools.GitStatus()
	resp, err := fn(context.Background(), &pluginv1.ToolRequest{
		Arguments: makeArgs(t, map[string]any{"path": dir}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Success {
		t.Fatalf("expected success: %s", resp.ErrorMessage)
	}
}

func TestGitStatus_WithModifiedFile(t *testing.T) {
	dir := initRepo(t)
	writeAndCommit(t, dir, "file.txt", "original\n", "init")
	// Modify without staging
	os.WriteFile(filepath.Join(dir, "file.txt"), []byte("changed\n"), 0644)

	fn := tools.GitStatus()
	resp, err := fn(context.Background(), &pluginv1.ToolRequest{
		Arguments: makeArgs(t, map[string]any{"path": dir}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Success {
		t.Fatalf("expected success: %s", resp.ErrorMessage)
	}
	text := responseText(t, resp)
	if !strings.Contains(text, "file.txt") {
		t.Fatalf("expected modified file in status, got: %s", text)
	}
}

// ---------- git_log ----------

func TestGitLog_ShowsCommits(t *testing.T) {
	dir := initRepo(t)
	writeAndCommit(t, dir, "a.txt", "a\n", "first commit")
	writeAndCommit(t, dir, "b.txt", "b\n", "second commit")

	fn := tools.GitLog()
	resp, err := fn(context.Background(), &pluginv1.ToolRequest{
		Arguments: makeArgs(t, map[string]any{"path": dir}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Success {
		t.Fatalf("expected success: %s", resp.ErrorMessage)
	}
	text := responseText(t, resp)
	if !strings.Contains(text, "first commit") || !strings.Contains(text, "second commit") {
		t.Fatalf("expected both commits in log, got: %s", text)
	}
}

// ---------- git_diff ----------

func TestGitDiff_ShowsDiff(t *testing.T) {
	dir := initRepo(t)
	writeAndCommit(t, dir, "file.txt", "original content\n", "init")
	os.WriteFile(filepath.Join(dir, "file.txt"), []byte("modified content\n"), 0644)

	fn := tools.GitDiff()
	resp, err := fn(context.Background(), &pluginv1.ToolRequest{
		Arguments: makeArgs(t, map[string]any{"path": dir}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Success {
		t.Fatalf("expected success: %s", resp.ErrorMessage)
	}
	text := responseText(t, resp)
	if !strings.Contains(text, "modified content") {
		t.Fatalf("expected diff in output, got: %s", text)
	}
}

// ---------- git_commit ----------

func TestGitCommit_MissingMessage(t *testing.T) {
	fn := tools.GitCommit()
	resp, err := fn(context.Background(), &pluginv1.ToolRequest{
		Arguments: makeArgs(t, map[string]any{"path": "/tmp"}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Success {
		t.Fatal("expected error when message missing")
	}
	if resp.ErrorCode != "validation_error" {
		t.Fatalf("expected validation_error, got %s", resp.ErrorCode)
	}
}

func TestGitCommit_WithFiles(t *testing.T) {
	dir := initRepo(t)
	// Create initial commit first so repo is not empty
	writeAndCommit(t, dir, "init.txt", "init\n", "bootstrap")
	// Add new untracked file
	os.WriteFile(filepath.Join(dir, "new.txt"), []byte("new file\n"), 0644)

	fn := tools.GitCommit()
	resp, err := fn(context.Background(), &pluginv1.ToolRequest{
		Arguments: makeArgs(t, map[string]any{
			"path":    dir,
			"message": "add new file",
			"files":   []any{"new.txt"},
		}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Success {
		t.Fatalf("expected success: %s — %s", resp.ErrorCode, resp.ErrorMessage)
	}
	text := responseText(t, resp)
	if !strings.Contains(text, "add new file") {
		t.Fatalf("expected commit message in output, got: %s", text)
	}
}

// ---------- git_branch ----------

func TestGitBranch_ListBranches(t *testing.T) {
	dir := initRepo(t)
	writeAndCommit(t, dir, "file.txt", "content\n", "init")

	fn := tools.GitBranch()
	resp, err := fn(context.Background(), &pluginv1.ToolRequest{
		Arguments: makeArgs(t, map[string]any{
			"path":   dir,
			"action": "list",
		}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Success {
		t.Fatalf("expected success: %s", resp.ErrorMessage)
	}
}

func TestGitBranch_CreateBranch(t *testing.T) {
	dir := initRepo(t)
	writeAndCommit(t, dir, "file.txt", "content\n", "init")

	fn := tools.GitBranch()
	resp, err := fn(context.Background(), &pluginv1.ToolRequest{
		Arguments: makeArgs(t, map[string]any{
			"path":   dir,
			"action": "create",
			"name":   "feature-xyz",
		}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Success {
		t.Fatalf("expected success: %s", resp.ErrorMessage)
	}
}

func TestGitBranch_MissingAction(t *testing.T) {
	fn := tools.GitBranch()
	resp, err := fn(context.Background(), &pluginv1.ToolRequest{
		Arguments: makeArgs(t, map[string]any{"path": "/tmp"}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Success {
		t.Fatal("expected error when action missing")
	}
}

// ---------- git_tag ----------

func TestGitTag_CreateTag(t *testing.T) {
	dir := initRepo(t)
	writeAndCommit(t, dir, "file.txt", "content\n", "init")

	fn := tools.GitTag()
	resp, err := fn(context.Background(), &pluginv1.ToolRequest{
		Arguments: makeArgs(t, map[string]any{
			"path":   dir,
			"action": "create",
			"name":   "v1.0.0",
		}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Success {
		t.Fatalf("expected success: %s", resp.ErrorMessage)
	}
}

func TestGitTag_ListTags(t *testing.T) {
	dir := initRepo(t)
	writeAndCommit(t, dir, "file.txt", "content\n", "init")
	// Create a tag first
	git.Run(context.Background(), dir, "tag", "v0.1.0")

	fn := tools.GitTag()
	resp, err := fn(context.Background(), &pluginv1.ToolRequest{
		Arguments: makeArgs(t, map[string]any{
			"path":   dir,
			"action": "list",
		}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Success {
		t.Fatalf("expected success: %s", resp.ErrorMessage)
	}
	text := responseText(t, resp)
	if !strings.Contains(text, "v0.1.0") {
		t.Fatalf("expected tag in list, got: %s", text)
	}
}

func TestGitTag_MissingAction(t *testing.T) {
	fn := tools.GitTag()
	resp, err := fn(context.Background(), &pluginv1.ToolRequest{
		Arguments: makeArgs(t, map[string]any{"path": "/tmp"}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Success {
		t.Fatal("expected error when action missing")
	}
}
