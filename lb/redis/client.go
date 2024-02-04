package redis

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func NewClient(ctx context.Context, conf Config) (redis.UniversalClient, error) {
	client := redis.UniversalClient(nil)
	nodes := strings.Split(conf.Nodes, ",")
	switch strings.ToLower(conf.Type) {
	case "sentinel":
		if conf.Master == "" {
			return nil, errors.New("Sentinel mode needs the Master name")
		}
		client = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    conf.Master,
			SentinelAddrs: nodes,
		})

	case "standalone":
		if len(nodes) != 1 {
			return nil, errors.New("standalone redis should have only one node")
		}

		client = redis.NewClient(&redis.Options{
			Addr: nodes[0],
		})
	default:
		return nil, errors.Errorf("Redis type %v is not supported", conf.Type)
	}

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return client, nil
}
