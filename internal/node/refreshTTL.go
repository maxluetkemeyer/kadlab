package node

import (
	"context"
	"d7024e_group04/env"
	"log"
	"time"
)

func (n *Node) Forget(hash string) {
	n.Store.RemoveRefreshContact(hash)
}

func (n *Node) TTLRefresher(ctx context.Context) {
	for {
		select {
		case newRefreshReq := <-n.RefreshChan:
			go n.RefreshTTL(ctx, newRefreshReq.Key, newRefreshReq.TTL)

		case <-ctx.Done():
			return
		}
	}
}

func (n *Node) RefreshTTL(ctx context.Context, key string, ttl time.Duration) {
	ticker := time.NewTicker(ttl)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			n.RWMutex.RLock()
			contacts := n.Store.GetStoreLocations(key)
			n.RWMutex.RUnlock()

			for _, contact := range contacts {
				ctx, cancel := context.WithTimeout(ctx, env.RPCTimeout)
				defer cancel()

				err := n.Client.SendRefreshTTL(ctx, key, contact)
				if err != nil {
					log.Printf("refreshTTL err: %v", err)
				}
			}

		case <-ctx.Done():
			return
		}
	}

}
