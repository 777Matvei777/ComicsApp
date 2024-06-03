package main

import (
	"context"
)

// MockServer is a mock of the Server interface for testing.
type MockServer struct {
	ShutdownInvoked bool
}

func (m *MockServer) RunServer() {
	// Mock implementation
}

func (m *MockServer) Shutdown(ctx context.Context) error {
	m.ShutdownInvoked = true
	return nil
}

// func TestMain(m *testing.M) {
// 	cfg := &config.Config{}
// 	s := &MockServer{}

// 	// Replace the NewServer function with a mock.
// 	server.NewServer = func(cfg *config.Config) *server.Server {
// 		return &MockServer{}
// 	}

// 	// Set up the signal channel.
// 	ch := make(chan os.Signal, 1)
// 	signal.Notify(ch, syscall.SIGINT)

// 	// Run the main function in a goroutine.
// 	go main()

// 	// Send a SIGINT signal.
// 	ch <- syscall.SIGINT

// 	// Give the main function some time to handle the signal.
// 	time.Sleep(100 * time.Millisecond)

// 	// Check if the server's Shutdown method was invoked.
// 	if !s.ShutdownInvoked {
// 		m.Errorf("expected server Shutdown to be invoked")
// 	}
// }
