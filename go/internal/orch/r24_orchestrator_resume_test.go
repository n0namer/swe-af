package orch

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestR24OrchestratorResumeRejectsUnreadableCheckpoint(t *testing.T) {
	scenarios := []struct {
		name       string
		checkpoint []byte
	}{
		{name: "missing"},
		{name: "malformed", checkpoint: []byte("{not-json\n")},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			repository := t.TempDir()
			checkpointDirectory := filepath.Join(repository, ".artifacts", "execution")
			if err := os.MkdirAll(checkpointDirectory, 0700); err != nil {
				t.Fatal(err)
			}
			checkpointPath := filepath.Join(checkpointDirectory, "checkpoint.json")
			if scenario.checkpoint != nil {
				if err := os.WriteFile(checkpointPath, scenario.checkpoint, 0600); err != nil {
					t.Fatal(err)
				}
			}
			_, err := ResumeBuildHandler(context.Background(), &Deps{}, map[string]any{
				"repo_path":     repository,
				"artifacts_dir": ".artifacts",
			})
			if err == nil {
				t.Fatal("unreadable checkpoint must not return successful resume")
			}
			actual, readErr := os.ReadFile(checkpointPath)
			if scenario.checkpoint == nil {
				if !os.IsNotExist(readErr) {
					t.Fatalf("missing checkpoint must remain absent: %v", readErr)
				}
			} else {
				if readErr != nil {
					t.Fatal(readErr)
				}
				if !bytes.Equal(actual, scenario.checkpoint) {
					t.Fatal("rejected malformed checkpoint was modified")
				}
			}
		})
	}
}
