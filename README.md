# Snowflake

Snowflake is a Go package for generating unique IDs based on the Snowflake algorithm. It provides flexibility to configure custom epochs, node IDs, and sequence bit sizes, making it ideal for distributed systems requiring unique, high-throughput ID generation.

## Features

- **Thread-safe** ID generation.
- Configurable node and sequence bit sizes.
- Custom epoch for timestamp calculations.
- ID decomposition into timestamp, node ID, and sequence components.
- Optional base62-encoded IDs of specific lengths.

## Installation

Install the package using `go get`:

```sh
go get github.com/srahkmli/gosnowflake
```

## Usage

### Basic Example

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/srahkmli/gosnowflake"
)

func main() {
	cfg := snowflake.Config{
		Epoch:        time.Now().Add(-time.Hour).UnixMilli(),
		NodeID:       1,
		NodeBits:     10,
		SequenceBits: 12,
	}

	ss, err := snowflake.NewSnowflake(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Snowflake: %v", err)
	}

	id, err := ss.GenerateID()
	if err != nil {
		log.Fatalf("Failed to generate ID: %v", err)
	}

	fmt.Printf("Generated ID: %d\n", id)
}
```

### Custom ID Length

Generate a base62-encoded ID with a specific length:

```go
customID, err := ss.GenerateCustomID(16)
if err != nil {
	log.Fatalf("Failed to generate custom ID: %v", err)
}
fmt.Printf("Generated Custom ID: %s\n", customID)
```

### Decomposing an ID

Extract the timestamp, node ID, and sequence from a generated ID:

```go
timestamp, nodeID, sequence := ss.DecomposeID(id)
fmt.Printf("Timestamp: %d, Node ID: %d, Sequence: %d\n", timestamp, nodeID, sequence)
```

## Testing

Run the included tests to verify the functionality:

```sh
go test ./...
```

## Configuration

The `Config` struct allows for the following options:

- **Epoch**: Start time for ID generation in milliseconds.
- **NodeID**: Unique identifier for this generator instance.
- **NodeBits**: Number of bits allocated for the node ID.
- **SequenceBits**: Number of bits allocated for the sequence.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

