package git_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/orchestra-mcp/plugin-devtools-git/internal/git"
)

// initRepo creates a temporary git repository and returns its path.
func initRepo(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	ctx := context.Background()
	cmds := [][]string{
		{"init"},
		{"config", "user.email", "test@example.com"},
		{"config", "user.name", "Test User"},
	}
	for _, args := range cmds {
		if _, err := git.Run(ctx, dir, args...); err != nil {
			t.Fatalf("git %v: %v", args, err)
		}
	}
	return dir
}

// writeAndCommit writes a file and commits it in the repo.
func writeAndCommit(t *testing.T, dir, filename, content, msg string) {
	t.Helper()
	ctx := context.Background()
	if err := os.WriteFile(filepath.Join(dir, filename), []byte(content), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	cmds := [][]string{
		{"add", filename},
		{"commit", "-m", msg},
	}
	for _, args := range cmds {
		if _, err := git.Run(ctx, dir, args...); err != nil {
			t.Fatalf("git %v: %v", args, err)
		}
	}
}

func TestRun_GitInit(t *testing.T) {
	dir, err := os.MkdirTemp("", "git-init-*")
	if err != nil {
		t.Fatalf("temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	out, err := git.Run(context.Background(), dir, "init")
	if err != nil {
		t.Fatalf("git init: %v", err)
	}
	if !strings.Contains(out, "Initialized") && !strings.Contains(out, "Reinitialized") {
		t.Fatalf("unexpected git init output: %s", out)
	}
}

func TestRun_GitStatus_Clean(t *testing.T) {
	dir := initRepo(t)
	writeAndCommit(t, dir, "README.md", "# Repo\n", "init")

	out, err := git.Run(context.Background(), dir, "status", "--porcelain")
	if err != nil {
		t.Fatalf("git status: %v", err)
	}
	// Clean repo = empty porcelain output
	if strings.TrimSpace(out) != "" {
		t.Fatalf("expected clean status, got: %s", out)
	}
}

func TestRun_GitLog(t *testing.T) {
	dir := initRepo(t)
	writeAndCommit(t, dir, "file.txt", "hello\n", "first commit")

	out, err := git.Run(context.Background(), dir, "log", "--oneline")
	if err != nil {
		t.Fatalf("git log: %v", err)
	}
	if !strings.Contains(out, "first commit") {
		t.Fatalf("expected commit in log, got: %s", out)
	}
}

func TestRun_GitDiff(t *testing.T) {
	dir := initRepo(t)
	writeAndCommit(t, dir, "file.txt", "original\n", "init")

	// Modify the file
	if err := os.WriteFile(filepath.Join(dir, "file.txt"), []byte("modified\n"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	out, err := git.Run(context.Background(), dir, "diff")
	if err != nil {
		t.Fatalf("git diff: %v", err)
	}
	if !strings.Contains(out, "modified") {
		t.Fatalf("expected diff to contain modified, got: %s", out)
	}
}

func TestRun_GitBranch_Create(t *testing.T) {
	dir := initRepo(t)
	writeAndCommit(t, dir, "file.txt", "hello\n", "init")

	_, err := git.Run(context.Background(), dir, "branch", "feature-branch")
	if err != nil {
		t.Fatalf("git branch create: %v", err)
	}

	out, err := git.Run(context.Background(), dir, "branch", "--list")
	if err != nil {
		t.Fatalf("git branch list: %v", err)
	}
	if !strings.Contains(out, "feature-branch") {
		t.Fatalf("expected feature-branch in branch list, got: %s", out)
	}
}

func TestRun_GitTag(t *testing.T) {
	dir := initRepo(t)
	writeAndCommit(t, dir, "file.txt", "hello\n", "init")

	_, err := git.Run(context.Background(), dir, "tag", "v1.0.0")
	if err != nil {
		t.Fatalf("git tag: %v", err)
	}

	out, err := git.Run(context.Background(), dir, "tag", "--list")
	if err != nil {
		t.Fatalf("git tag list: %v", err)
	}
	if !strings.Contains(out, "v1.0.0") {
		t.Fatalf("expected v1.0.0 in tags, got: %s", out)
	}
}

func TestRun_GitBlame(t *testing.T) {
	dir := initRepo(t)
	writeAndCommit(t, dir, "blame.txt", "line one\nline two\n", "initial blame file")

	out, err := git.Run(context.Background(), dir, "blame", "--porcelain", "blame.txt")
	if err != nil {
		t.Fatalf("git blame: %v", err)
	}
	if !strings.Contains(out, "blame.txt") {
		t.Fatalf("expected filename in blame output, got: %s", out)
	}
}

func TestRun_GitStash(t *testing.T) {
	dir := initRepo(t)
	writeAndCommit(t, dir, "file.txt", "original\n", "init")

	// Modify file so there's something to stash
	if err := os.WriteFile(filepath.Join(dir, "file.txt"), []byte("modified\n"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	out, err := git.Run(context.Background(), dir, "stash", "push", "-m", "my stash")
	if err != nil {
		t.Fatalf("git stash: %v", err)
	}
	if !strings.Contains(out, "Saved") {
		t.Fatalf("expected stash saved message, got: %s", out)
	}
}

func TestRun_InvalidCommand_ReturnsError(t *testing.T) {
	dir := initRepo(t)
	_, err := git.Run(context.Background(), dir, "nonexistent-subcommand")
	if err == nil {
		t.Fatal("expected error for invalid git command")
	}
}
