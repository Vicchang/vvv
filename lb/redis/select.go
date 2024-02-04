package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vvv/lb"
	"github.com/vvv/lb/wrr"
)

var _ lb.SelectService = (*SelectService)(nil)

type SelectService struct {
	client redis.UniversalClient
	wrr    wrr.WRR
}

func NewSelectService(client redis.UniversalClient, wrr wrr.WRR) *SelectService {
	return &SelectService{
		client: client,
		wrr:    wrr,
	}
}

func (srv *SelectService) ServerURI() (string, error) {
	uri := srv.wrr.Next()
	if uri == nil {
		return "", fmt.Errorf("no available server")
	}

	return uri.(string), nil
}

func (srv *SelectService) BackgroundUpdate(ctx context.Context) {
	srv.bgUpdate(ctx)

	go func() {
		ticker := time.NewTicker(time.Second * 60)
		defer ticker.Stop()
		for range ticker.C {
			srv.bgUpdate(ctx)
		}
	}()
}

func (srv *SelectService) bgUpdate(ctx context.Context) {
	resp, err := srv.client.HGetAll(ctx, "server_status").Result()
	if err != nil {
		fmt.Printf("get server status fail, %s", err)
		return
	}

	// TODO: change value to be struct in order to store more data
	for k, conns := range resp {
		c, err := strconv.Atoi(conns)
		if err != nil {
			fmt.Println("invalid data in server_status")
			continue
		}

		srv.wrr.Update(k, uint32(c))
	}
}
