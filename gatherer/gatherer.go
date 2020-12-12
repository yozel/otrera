package gatherer

import (
	"time"
)

// Gatherer is an struct to gather Description
type Gatherer struct {
	cachePath   string
	descriptors map[string]func(options map[string]string) ([]RawObjectInterface, error)
}

// Gather returns Description for given name and options with cache
func (g *Gatherer) Gather(key string, options map[string]string, ttl time.Duration) ([]RawObjectInterface, error) {
	r, err := g.descriptors[key](options)
	if err != nil {
		return nil, err // TODO: wrap error
	}
	return r, nil
}

// New creates a new Gatherer
func New(cachePath string, descriptors map[string]func(options map[string]string) ([]RawObjectInterface, error)) *Gatherer {
	g := &Gatherer{
		cachePath:   cachePath,
		descriptors: descriptors,
	}
	return g
}
