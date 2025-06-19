//go:build mage
// +build mage

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/azrod/kivigo/magefiles/testinfra"
)

// StartBackend starts a specific backend (e.g. mage startBackend redis)
func StartBackend(name string) error {
	backend := testinfra.GetBackend(name)
	if backend == nil {
		return fmt.Errorf("backend %q not found", name)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	fmt.Printf("🚀 Starting backend: %s...\n", name)
	if err := backend.Start(ctx); err != nil {
		return fmt.Errorf("failed to start backend %q: %w", name, err)
	}
	fmt.Printf("✅ Backend %s started successfully.\n", name)
	return nil
}

// StopBackend stops a specific backend (e.g. mage stopBackend redis)
func StopBackend(name string) error {
	backend := testinfra.GetBackend(name)
	if backend == nil {
		return fmt.Errorf("backend %q not found", name)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	fmt.Printf("🛑 Stopping backend: %s...\n", name)
	if err := backend.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop backend %q: %w", name, err)
	}
	fmt.Printf("✅ Backend %s stopped successfully.\n", name)
	return nil
}

// TestBackend runs tests for a specific backend (e.g. mage testBackend redis)
func TestBackend(name string) error {
	backend := testinfra.GetBackend(name)
	if backend == nil {
		return fmt.Errorf("backend %q not found", name)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if err := backend.Start(ctx); err != nil {
		return err
	}
	defer backend.Stop(ctx)

	fmt.Printf("🧪 Running tests for backend: %s...\n", name)
	cmd := exec.Command("go", "test", "-v", "-cover", fmt.Sprintf("github.com/azrod/kivigo/pkg/backend/%s", backend.Name()))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// TestAllBackends runs tests for all registered backends
func TestAllBackends() error {
	for _, name := range testinfra.ListBackends() {
		if err := TestBackend(name); err != nil {
			return err
		}
	}
	return nil
}

// ListOfBackends lists all registered backends
func ListOfBackends() error {
	if len(testinfra.ListBackends()) == 0 {
		fmt.Println("No backends registered.")
		return nil
	}
	fmt.Println("🔍 Listing all registered backends:")
	for _, name := range testinfra.ListBackends() {
		fmt.Println(name)
	}
	return nil
}
