package gatherer

// Gatherer is an struct to gather Description
type Gatherer struct {
	descriptors map[string]func(options map[string]string) ([]RawObjectInterface, error)
}

// Gather returns Description for given name and options
func (g *Gatherer) Gather(key string, options map[string]string) ([]RawObject, error) {
	r, err := g.descriptors[key](options)
	if err != nil {
		return nil, err // TODO: wrap error
	}

	result := make([]RawObject, len(r))
	for i, obj := range r {
		result[i] = obj.Copy()
	}
	return result, nil
}

// New creates a new Gatherer
func New(descriptors map[string]func(options map[string]string) ([]RawObjectInterface, error)) *Gatherer {
	g := &Gatherer{descriptors: descriptors}
	return g
}
