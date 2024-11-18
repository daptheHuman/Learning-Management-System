// services/redis/redis.go
package redis

import (
	"github.com/redis/go-redis/v9"
)

// NewClient create a new redis client.
// - The host argument is in the form `host:port`.
// - The db argument is the identifier for the database to be selected after connecting to the server.
// - `MaxRetries` is set to -1 so it is disabled.
func NewClient(host, password string, db int) redis.UniversalClient {
	opts := &redis.Options{
		Addr:       host,
		Password:   password,
		DB:         db,
		MaxRetries: -1, // Disable retry
	}
	return redis.NewClient(opts)
}
