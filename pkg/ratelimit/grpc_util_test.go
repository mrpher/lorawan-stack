// Copyright © 2021 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ratelimit_test

import (
	"context"
	"net"

	pbtypes "github.com/gogo/protobuf/types"
	"go.thethings.network/lorawan-stack/v3/pkg/rpcserver"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"google.golang.org/grpc"
)

// mockServer
type mockServer struct {
}

// Link a gateway to the Gateway Server for streaming upstream messages and downstream messages.
func (s *mockServer) LinkGateway(stream ttnpb.GtwGs_LinkGatewayServer) error {
	return nil
}

// Get configuration for the concentrator.
func (s *mockServer) GetConcentratorConfig(_ context.Context, _ *pbtypes.Empty) (*ttnpb.ConcentratorConfig, error) {
	return &ttnpb.ConcentratorConfig{}, nil
}

// Get connection information to connect an MQTT gateway.
func (s *mockServer) GetMQTTConnectionInfo(_ context.Context, _ *ttnpb.GatewayIdentifiers) (*ttnpb.MQTTConnectionInfo, error) {
	return &ttnpb.MQTTConnectionInfo{}, nil
}

// Get legacy connection information to connect a The Things Network Stack V2 MQTT gateway.
func (s *mockServer) GetMQTTV2ConnectionInfo(_ context.Context, _ *ttnpb.GatewayIdentifiers) (*ttnpb.MQTTConnectionInfo, error) {
	return &ttnpb.MQTTConnectionInfo{}, nil
}

func withRateLimitedGRPCServer(conf rpcserver.RateLimitingConfig, f func(ttnpb.GtwGsClient)) {
	ctx, cancel := context.WithCancel(test.Context())
	s := rpcserver.New(test.Context(), rpcserver.WithRateLimitingConfig(conf))
	ttnpb.RegisterGtwGsServer(s.Server, &mockServer{})
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	go func() {
		s.Serve(l)
	}()
	go func() {
		<-ctx.Done()
		s.Stop()
	}()
	cc, err := grpc.Dial(l.Addr().String(), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	f(ttnpb.NewGtwGsClient(cc))
	cancel()
}
