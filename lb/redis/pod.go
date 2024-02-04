package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/redis/go-redis/v9"
)

// race condition?
type PodService struct {
	client redis.UniversalClient
	table  map[string]struct{}
}

func NewPodService(client redis.UniversalClient) *PodService {
	return &PodService{
		client: client,
		table:  map[string]struct{}{},
	}
}

func (srv *PodService) Add(url string) {
	srv.table[url] = struct{}{}
}

// vpc should maintain the virtual ip instead of update the pod status,
// shouldn't this be lb's responsibility
func (srv *PodService) BackgroundUpdate(ctx context.Context) {
	go func() {
		for {
			for u := range srv.table {
				timer := time.NewTimer(time.Second)
				<-timer.C

				hearbeatURL, err := url.JoinPath(u, "heartbeat")
				if err != nil {
					fmt.Printf("join path %s err, %s\n", u, err)
					continue
				}

				resp, err := http.Get(hearbeatURL)
				if err != nil {
					fmt.Printf("get %s err, %s\n", hearbeatURL, err)
					continue
				}

				if resp.StatusCode != http.StatusOK {
					if resp.StatusCode == http.StatusNotFound {
						delete(srv.table, u)
					}
					continue
				}

				var status struct {
					Connections int `json:"connections"`
				}
				bs, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Printf("read body fail, %s\n", err)
					continue
				}

				err = json.Unmarshal(bs, &status)
				if err != nil {
					fmt.Printf("unmarshal body fail, %s\n", err)
					continue
				}

				err = srv.client.HSet(ctx, "server_status", u, status.Connections).Err()
				if err != nil {
					fmt.Printf("set server status fail, %s, %d, %s\n", u, status.Connections, err)
					continue
				}
			}

			timer := time.NewTimer(time.Second)
			<-timer.C
		}
	}()
}
