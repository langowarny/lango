package app

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// registerP2PRoutes mounts P2P status endpoints on the gateway router.
// Endpoints are public (no auth) since they expose only node metadata.
func registerP2PRoutes(r chi.Router, p2pc *p2pComponents) {
	r.Route("/api/p2p", func(r chi.Router) {
		r.Get("/status", p2pStatusHandler(p2pc))
		r.Get("/peers", p2pPeersHandler(p2pc))
		r.Get("/identity", p2pIdentityHandler(p2pc))
	})
}

func p2pStatusHandler(p2pc *p2pComponents) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		node := p2pc.node

		addrs := make([]string, 0, len(node.Multiaddrs()))
		for _, a := range node.Multiaddrs() {
			addrs = append(addrs, a.String())
		}

		resp := map[string]interface{}{
			"peerId":         node.PeerID().String(),
			"listenAddrs":    addrs,
			"connectedPeers": len(node.ConnectedPeers()),
			"mdnsEnabled":    p2pc.node.Host().Addrs() != nil,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func p2pPeersHandler(p2pc *p2pComponents) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		node := p2pc.node
		connected := node.ConnectedPeers()

		type peerInfo struct {
			PeerID string   `json:"peerId"`
			Addrs  []string `json:"addrs"`
		}

		peers := make([]peerInfo, 0, len(connected))
		for _, pid := range connected {
			conns := node.Host().Network().ConnsToPeer(pid)
			var addrs []string
			for _, c := range conns {
				addrs = append(addrs, c.RemoteMultiaddr().String())
			}
			peers = append(peers, peerInfo{
				PeerID: pid.String(),
				Addrs:  addrs,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"peers": peers,
			"count": len(peers),
		})
	}
}

func p2pIdentityHandler(p2pc *p2pComponents) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if p2pc.identity == nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"did":    nil,
				"peerId": p2pc.node.PeerID().String(),
			})
			return
		}

		ctx := context.Background()
		did, err := p2pc.identity.DID(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"did":    did.ID,
			"peerId": p2pc.node.PeerID().String(),
		})
	}
}
