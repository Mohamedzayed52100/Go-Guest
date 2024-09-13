package listeners

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/goplaceapp/goplace-settings/config"
	"github.com/lib/pq"
	"sync"
	"time"
)

var mu sync.Mutex

func GetListener(dbName string, channelName string) (*pq.Listener, func()) {
	mu.Lock()
	defer mu.Unlock()

	cfg := &config.Config{}
	if err := env.Parse(cfg); err != nil {
		panic(fmt.Errorf("failed to read service config, %w", err))
	}

	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require",
		cfg.DbPostgresHost, cfg.DbPostgresUser, cfg.DbPostgresPassword, dbName)

	listener := pq.NewListener(connectionString, 10*time.Second, time.Minute, nil)
	if listener == nil {
		return nil, nil
	}

	err := listener.Listen(channelName)
	if err != nil {
		return nil, nil
	}

	cleanup := func() {
		listener.Unlisten(channelName)
		listener.Close()
	}

	return listener, cleanup
}
