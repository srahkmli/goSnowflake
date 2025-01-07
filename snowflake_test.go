package snowflake

import (
	"testing"
	"time"
)

func TestGenerateID(t *testing.T) {
	cfg := Config{
		Epoch:        time.Now().Add(-time.Hour).UnixMilli(), // Set a custom epoch an hour ago.
		NodeID:       1,
		NodeBits:     10,
		SequenceBits: 12,
	}

	ss, err := NewSnowFlake(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize SnowFlake: %v", err)
	}

	id, err := ss.GenerateID()
	if err != nil {
		t.Errorf("GenerateID failed: %v", err)
	}

	if id <= 0 {
		t.Errorf("Generated ID should be greater than 0, got %d", id)
	}
}

func TestDecomposeID(t *testing.T) {
	cfg := Config{
		Epoch:        time.Now().UnixMilli(),
		NodeID:       1,
		NodeBits:     10,
		SequenceBits: 12,
	}

	ss, err := NewSnowFlake(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize SnowFlake: %v", err)
	}

	id, err := ss.GenerateID()
	if err != nil {
		t.Fatalf("GenerateID failed: %v", err)
	}

	timestamp, nodeID, sequence := ss.DecomposeID(id)

	if nodeID != cfg.NodeID {
		t.Errorf("Expected nodeID %d, got %d", cfg.NodeID, nodeID)
	}

	if timestamp <= 0 || sequence < 0 {
		t.Errorf("Decomposed values are invalid: timestamp=%d, sequence=%d", timestamp, sequence)
	}
}

func TestGenerateCustomID(t *testing.T) {
	cfg := Config{
		Epoch:        time.Now().UnixMilli(),
		NodeID:       1,
		NodeBits:     10,
		SequenceBits: 12,
	}

	ss, err := NewSnowFlake(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize SnowFlake: %v", err)
	}

	length := 16
	customID, err := ss.GenerateCustomID(length)
	if err != nil {
		t.Fatalf("GenerateCustomID failed: %v", err)
	}

	if len(customID) != length {
		t.Errorf("Expected custom ID length %d, got %d", length, len(customID))
	}
}

func TestValidateNodeID(t *testing.T) {
	cfg := Config{
		Epoch:        time.Now().UnixMilli(),
		NodeID:       1,
		NodeBits:     10,
		SequenceBits: 12,
	}

	ss, err := NewSnowFlake(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize SnowFlake: %v", err)
	}

	if err := ss.ValidateNodeID(5); err != nil {
		t.Errorf("ValidateNodeID failed for valid node ID: %v", err)
	}

	invalidNodeID := int64(1 << 11) // Out of range for NodeBits = 10.
	if err := ss.ValidateNodeID(invalidNodeID); err == nil {
		t.Errorf("ValidateNodeID did not fail for invalid node ID: %d", invalidNodeID)
	}
}
