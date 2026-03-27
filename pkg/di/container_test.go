package di_test

import (
	"testing"

	"github.com/paladignus/actajus/pkg/di"
)

func TestContainerCreation(t *testing.T) {
	t.Skip("Skipping test - requires database and redis connection")

	container, err := di.NewContainer()
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}
	defer container.Close()

	if container.DB == nil {
		t.Error("Expected DB to be initialized")
	}

	if container.Logger == nil {
		t.Error("Expected Logger to be initialized")
	}

	if container.UoW == nil {
		t.Error("Expected UoW to be initialized")
	}
}

func TestContainerClose(t *testing.T) {
	t.Skip("Skipping test - requires database connection")

	container, err := di.NewContainer()
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}

	// Close should not panic
	container.Close()
}

func TestContainerTypeCompiles(t *testing.T) {
	var c di.Container
	if c.Modules != (di.Modules{}) {
		t.Fatal("expected zero-value modules")
	}
}
