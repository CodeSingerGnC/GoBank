package gapi

import (
	"context"
	// "log"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgenHeader 	= "grpcgateway-user-agent"
	userAgentHeader 			= "user-agent"
	xForwardedForHeader 		= "x-forward-for"
)

type Metadata struct {
	UserAgent string
	ClientAPI string
}

func (server *Server) extraMetadata(ctx context.Context) *Metadata {
	meta := &Metadata{}

	if mt, ok := metadata.FromIncomingContext(ctx); ok {
		// -TODO: remove
		// log.Printf("mt %+v\n", mt)
		if userAgents := mt.Get(grpcGatewayUserAgenHeader); len(userAgents) > 0 {
			meta.UserAgent = userAgents[0]
		}

		if userAgents := mt.Get(userAgentHeader); len(userAgents) > 0 {
			meta.UserAgent = userAgents[0]
		}

		if clientIPs := mt.Get(xForwardedForHeader); len(clientIPs) > 0 {
			meta.ClientAPI = clientIPs[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		meta.ClientAPI = p.Addr.String()
	}

	return meta
}