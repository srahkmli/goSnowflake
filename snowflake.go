package snowflake

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// SnowFlake is a generator for creating unique IDs based on the Snowflake algorithm.
type SnowFlake struct {
	mu            sync.Mutex // Ensures thread-safe ID generation.
	epoch         int64      // Custom epoch to calculate timestamps.
	nodeID        int64      // Unique identifier for the generator instance.
	sequence      int64      // Tracks the sequence number within the same millisecond.
	lastTimestamp int64      // Keeps the last used timestamp to handle clock adjustments.
	maxSequence   int64      // Maximum value the sequence can take.
	nodeShift     uint       // Bit shift for the node ID.
	sequenceShift uint       // Bit shift for the sequence.
	maxNodeID     int64      // Maximum valid node ID.
}

// Config holds the settings for initializing a SnowFlake generator.
type Config struct {
	Epoch        int64 // Start time for ID generation in milliseconds.
	NodeID       int64 // Unique node ID for this generator.
	NodeBits     int   // Number of bits allocated for the node ID.
	SequenceBits int   // Number of bits allocated for the sequence.
}

// NewSnowFlake creates and configures a new instance of the SnowFlake generator.
func NewSnowFlake(cfg Config) (*SnowFlake, error) {
	if cfg.NodeBits+cfg.SequenceBits >= 63 {
		return nil, errors.New("the sum of NodeBits and SequenceBits must be less than 63")
	}

	maxNodeID := (1 << cfg.NodeBits) - 1
	if cfg.NodeID < 0 || cfg.NodeID > int64(maxNodeID) {
		return nil, fmt.Errorf("nodeID must be between 0 and %d", maxNodeID)
	}

	return &SnowFlake{
		epoch:         cfg.Epoch,
		nodeID:        cfg.NodeID,
		maxSequence:   (1 << cfg.SequenceBits) - 1,
		nodeShift:     uint(cfg.SequenceBits),
		sequenceShift: uint(63 - cfg.NodeBits - cfg.SequenceBits),
		maxNodeID:     int64(maxNodeID),
	}, nil
}

// GenerateID produces a unique ID using the current timestamp, node ID, and sequence number.
func (ss *SnowFlake) GenerateID() (int64, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	now := time.Now().UnixMilli()
	if now < ss.lastTimestamp {
		return 0, fmt.Errorf("system clock moved backward: refusing to generate ID for %d milliseconds", ss.lastTimestamp-now)
	}

	if now == ss.lastTimestamp {
		ss.sequence = (ss.sequence + 1) & ss.maxSequence
		if ss.sequence == 0 {
			// Wait for the next millisecond when the sequence overflows.
			for now <= ss.lastTimestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		ss.sequence = 0 // Reset sequence for a new millisecond.
	}

	ss.lastTimestamp = now

	// Construct the ID by combining timestamp, node ID, and sequence.
	id := ((now - ss.epoch) << ss.sequenceShift) |
		(ss.nodeID << ss.nodeShift) |
		ss.sequence

	return id, nil
}

// DecomposeID breaks down a generated ID into its timestamp, node ID, and sequence components.
func (ss *SnowFlake) DecomposeID(id int64) (timestamp, nodeID, sequence int64) {
	timestamp = (id >> ss.sequenceShift) + ss.epoch
	nodeID = (id >> ss.nodeShift) & ((1 << uint(63-ss.sequenceShift-ss.nodeShift)) - 1)
	sequence = id & ss.maxSequence
	return
}

// GenerateCustomID creates a unique ID and returns it as a base62-encoded string with a specific length.
func (ss *SnowFlake) GenerateCustomID(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("length must be greater than zero")
	}

	id, err := ss.GenerateID()
	if err != nil {
		return "", err
	}

	idStr := encodeToBase62(id)

	// Adjust the ID length by truncating or padding it as needed.
	if len(idStr) > length {
		return idStr[:length], nil
	} else if len(idStr) < length {
		return padLeft(idStr, '0', length), nil
	}

	return idStr, nil
}

// encodeToBase62 converts an integer ID into a base62-encoded string.
func encodeToBase62(id int64) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	if id == 0 {
		return "0"
	}

	result := make([]byte, 0)
	for id > 0 {
		result = append([]byte{charset[id%62]}, result...)
		id /= 62
	}
	return string(result)
}

// padLeft adds padding characters to the left of a string to meet a specific length.
func padLeft(str string, pad rune, length int) string {
	if len(str) >= length {
		return str
	}
	padding := make([]rune, length-len(str))
	for i := range padding {
		padding[i] = pad
	}
	return string(padding) + str
}

// ValidateNodeID checks if a given node ID is within the valid range for this generator.
func (ss *SnowFlake) ValidateNodeID(nodeID int64) error {
	if nodeID < 0 || nodeID > ss.maxNodeID {
		return fmt.Errorf("nodeID must be between 0 and %d", ss.maxNodeID)
	}
	return nil
}

// Version specifies the current version of the SnowFlake package.
const Version = "1.0.0"

// Example Usage:
// cfg := snowstorm.Config{
// 	Epoch:        time.Now().UnixMilli(),
// 	NodeID:       1,
// 	NodeBits:     10,
// 	SequenceBits: 12,
// }
// ss, err := snowstorm.NewSnowFlake(cfg)
// if err != nil {
// 	log.Fatalf("Failed to initialize SnowFlake: %v", err)
// }
// id, err := ss.GenerateID()
// if err != nil {
// 	log.Fatalf("Failed to generate ID: %v", err)
// }
// fmt.Printf("Generated ID: %d\n", id)
