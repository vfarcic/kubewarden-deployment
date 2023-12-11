package main

import (
	"encoding/json"
	"testing"
)

func TestParsingSettingsWithNoValueProvided(t *testing.T) {
	rawSettings := []byte(`{}`)
	settings := &Settings{}
	if err := json.Unmarshal(rawSettings, settings); err != nil {
		t.Errorf("Unexpected error %+v", err)
	}

	if settings.ReplicasGreaterThan != 0 {
		t.Errorf("Expected ReplicasGreaterThan to be empty")
	}

	valid, err := settings.Valid()
	if !valid {
		t.Errorf("Settings are reported as not valid")
	}
	if err != nil {
		t.Errorf("Unexpected error %+v", err)
	}
}

func TestIsReplicasGreaterThanDenied(t *testing.T) {
	settings := Settings{
		ReplicasGreaterThan: 2,
	}
	if settings.IsReplicasGreaterThanAllowed(0) {
		t.Errorf("replicasGreaterThan should NOT be allowed")
	}
	if settings.IsReplicasGreaterThanAllowed(1) {
		t.Errorf("replicasGreaterThan should be allowed")
	}
	if settings.IsReplicasGreaterThanAllowed(2) {
		t.Errorf("replicasGreaterThan should be allowed")
	}
	if !settings.IsReplicasGreaterThanAllowed(3) {
		t.Errorf("replicasGreaterThan should NOT be allowed")
	}
}
