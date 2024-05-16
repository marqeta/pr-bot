package rate

import "time"

type Limit struct {
	Window    time.Duration `dynamodbav:"-"`
	Value     int64         `dynamodbav:"value"`
	WindowStr string        `dynamodbav:"window"`
}

type LimiterConfig struct {
	Default   Limit            `dynamodbav:"default"`
	Overrides map[string]Limit `dynamodbav:"overrides"`
}

// Update is called everytime the store updates the config
func (cfg *LimiterConfig) Update() error {
	d, err := time.ParseDuration(cfg.Default.WindowStr)
	if err != nil {
		return err
	}
	cfg.Default.Window = d
	for k, limit := range cfg.Overrides {
		d, err := time.ParseDuration(limit.WindowStr)
		if err != nil {
			return err
		}
		limit.Window = d
		cfg.Overrides[k] = limit
	}
	return nil
}

func (cfg *LimiterConfig) Get(key string) Limit {
	o, ok := cfg.Overrides[key]
	if ok {
		return o
	}
	return cfg.Default
}
