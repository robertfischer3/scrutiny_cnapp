package database

import (
	"fmt"
	"sync"
)

// Registry manages database providers
type Registry struct {
	mu        sync.RWMutex
	providers map[string]Provider
}

// NewRegistry creates a new registry
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

// Register adds a provider to the registry
func (r *Registry) Register(provider Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	name := provider.Name()
	if _, exists := r.providers[name]; exists {
		// Just log and override instead of failing
		// This allows for provider replacements in tests
		fmt.Printf("Warning: Overriding existing database provider: %s\n", name)
	}
	
	r.providers[name] = provider
}

// Get retrieves a provider by name
func (r *Registry) Get(name string) (Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	provider, exists := r.providers[name]
	if !exists {
		return nil, fmt.Errorf("database provider not found: %s", name)
	}
	
	return provider, nil
}

// List returns all registered provider names
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	
	return names
}

// DefaultRegistry is the global provider registry
var DefaultRegistry = NewRegistry()

// Register adds a provider to the default registry
func Register(provider Provider) {
	DefaultRegistry.Register(provider)
}

// Get retrieves a provider from the default registry
func Get(name string) (Provider, error) {
	return DefaultRegistry.Get(name)
}

// List returns all registered provider names from the default registry
func List() []string {
	return DefaultRegistry.List()
}