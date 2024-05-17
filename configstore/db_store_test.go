package configstore_test

import (
	"errors"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/marqeta/pr-bot/configstore"
	"github.com/marqeta/pr-bot/metrics"
	"github.com/stretchr/testify/assert"
)

type exampleConfig struct {
	str string
	num int
	arr []string
	m   map[string]string
}

func (e *exampleConfig) Update() error {
	return nil
}

func randomCfg(num int) *exampleConfig {
	return &exampleConfig{
		str: "example",
		num: num,
		arr: []string{"abc", "abc"},
		m: map[string]string{
			"a": "a",
			"b": "b",
			"c": "c",
		},
	}
}

func Test_Config_Store(t *testing.T) {
	//nolint:goerr113
	randomErr := errors.New("random error")
	tableName := "tablename"
	name := "example"
	tests := []struct {
		name            string
		duration        time.Duration
		setExpectations func(dao *configstore.MockDao[*exampleConfig])
		verifyCfg       func(t *testing.T, store configstore.Getter[*exampleConfig], clock clockwork.FakeClock)
		wantErrOnNew    error
	}{
		{
			name:     "Should load example config on creation",
			duration: 10 * time.Second,
			setExpectations: func(dao *configstore.MockDao[*exampleConfig]) {
				dao.EXPECT().GetItem(name, tableName).
					Return(randomCfg(10), nil).Once()
			},
			verifyCfg: func(t *testing.T, store configstore.Getter[*exampleConfig], _ clockwork.FakeClock) {
				verify(t, store, randomCfg(10), nil)
			},
			wantErrOnNew: nil,
		},
		{
			name: "Should return error on creation",
			setExpectations: func(dao *configstore.MockDao[*exampleConfig]) {
				dao.EXPECT().GetItem(name, tableName).
					Return(randomCfg(10), randomErr).Once()
			},
			verifyCfg: func(_ *testing.T, _ configstore.Getter[*exampleConfig], _ clockwork.FakeClock) {
			},
			wantErrOnNew: randomErr,
		},
		{
			name:     "Should load config peridically",
			duration: 10 * time.Second,
			setExpectations: func(dao *configstore.MockDao[*exampleConfig]) {
				dao.EXPECT().GetItem(name, tableName).
					Return(randomCfg(0), nil).Once()
				dao.EXPECT().GetItem(name, tableName).
					Return(randomCfg(10), nil).Once()
				dao.EXPECT().GetItem(name, tableName).
					Return(randomCfg(20), nil).Once()
				dao.EXPECT().GetItem(name, tableName).
					Return(randomCfg(30), nil).Once()
			},
			verifyCfg: func(t *testing.T, store configstore.Getter[*exampleConfig], c clockwork.FakeClock) {
				verify(t, store, randomCfg(0), nil)
				advanceTime(c, 10*time.Second)
				verify(t, store, randomCfg(10), nil)
				advanceTime(c, 10*time.Second)
				verify(t, store, randomCfg(20), nil)
				advanceTime(c, 10*time.Second)
				verify(t, store, randomCfg(30), nil)
			},
			wantErrOnNew: nil,
		},
		{
			name:     "Should return old config when load fails",
			duration: 10 * time.Second,
			setExpectations: func(dao *configstore.MockDao[*exampleConfig]) {
				dao.EXPECT().GetItem(name, tableName).
					Return(randomCfg(0), nil).Once()
				dao.EXPECT().GetItem(name, tableName).
					Return(randomCfg(10), nil).Once()
				dao.EXPECT().GetItem(name, tableName).
					Return(randomCfg(-1), randomErr).Once()
				dao.EXPECT().GetItem(name, tableName).
					Return(randomCfg(-1), randomErr).Once()
			},
			verifyCfg: func(t *testing.T, store configstore.Getter[*exampleConfig], c clockwork.FakeClock) {
				verify(t, store, randomCfg(0), nil)
				advanceTime(c, 10*time.Second)
				verify(t, store, randomCfg(10), nil)
				advanceTime(c, 10*time.Second)
				verify(t, store, randomCfg(10), nil)
				advanceTime(c, 10*time.Second)
				verify(t, store, randomCfg(10), nil)
			},
			wantErrOnNew: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := configstore.NewMockDao[*exampleConfig](t)
			clock := clockwork.NewFakeClock()
			ticker := clock.NewTicker(tt.duration)

			tt.setExpectations(dao)
			store, err := configstore.NewDBStore[*exampleConfig](dao, name,
				tableName, ticker, metrics.NewNoopEmitter())
			if !errors.Is(err, tt.wantErrOnNew) {
				t.Errorf("configstore.NewDBStore() error = %v, wantErr %v", err, tt.wantErrOnNew)
				return
			}
			tt.verifyCfg(t, store, clock)
			if store != nil {
				store.Close()
			}
		})
	}
}

func verify(t *testing.T, store configstore.Getter[*exampleConfig],
	wantCfg *exampleConfig, wantErr error) {
	gotCfg, err := store.Get()
	assert.Equal(t, wantCfg, gotCfg)
	if !errors.Is(err, wantErr) {
		t.Errorf("configstore.NewDBStore() error = %v, wantErr %v", err, wantErr)
		return
	}

}

func advanceTime(c clockwork.FakeClock, d time.Duration) {
	c.Advance(d)
	// wait for the lock to be acquired by atomic.Value.Store() before atomic.Value.Load() is called.
	time.Sleep(100 * time.Millisecond)
}
