package configstore

type inMemoryStore[T DynamicConfig] struct {
	cfg T
}

// Close implements Getter
func (s *inMemoryStore[T]) Close() {
	// no-op
}

// Get implements Getter
func (s *inMemoryStore[T]) Get() (T, error) {
	return s.cfg, nil
}

func NewInMemoryStore[T DynamicConfig](cfg T) (Getter[T], error) {
	err := cfg.Update()
	if err != nil {
		return nil, err
	}
	return &inMemoryStore[T]{
		cfg: cfg,
	}, nil
}
