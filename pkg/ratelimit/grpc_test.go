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
	"fmt"
	"sync"
	"testing"
	"time"

	pbtypes "github.com/gogo/protobuf/types"
	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/ratelimit"
	"go.thethings.network/lorawan-stack/v3/pkg/rpcserver"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	maxRate      uint64 = 10
	overrideRate uint64 = 15

	timeout = (1 << 5) * test.Delay
)

func TestGRPC(t *testing.T) {
	conf := rpcserver.RateLimitingConfig{
		ByRemoteIP: ratelimit.Config{
			Enable: true,
			Rate:   maxRate,
			Overrides: map[string]uint64{
				"/ttn.lorawan.v3.GtwGs/GetMQTTConnectionInfo": overrideRate,
			},
		},
		Stream: ratelimit.Config{
			Enable: true,
			Rate:   maxRate,
		},
	}

	withRateLimitedGRPCServer(rpcserver.RateLimitingConfig{}, func(server *mockServer, client ttnpb.GtwGsClient) {
		t.Run("NoRateLimit", func(t *testing.T) {
			t.Run("Invoke", func(t *testing.T) {
				a := assertions.New(t)
				for i := uint64(0); i < 2*maxRate; i++ {
					headers := metadata.MD{}
					_, err := client.GetConcentratorConfig(test.Context(), &pbtypes.Empty{}, grpc.Header(&headers))
					a.So(err, should.BeNil)
					a.So(headers.Get("x-rate-limit-limit"), should.Resemble, []string{"0"})
					a.So(headers.Get("x-rate-limit-reset"), should.Resemble, []string{"0"})
					a.So(headers.Get("x-rate-limit-available"), should.Resemble, []string{"0"})
				}
			})

			t.Run("Stream", func(t *testing.T) {
				a := assertions.New(t)

				ctx := metadata.NewOutgoingContext(test.Context(), metadata.Pairs("test-stream-id", "1"))
				stream, err := client.LinkGateway(ctx)
				if err != nil {
					panic(err)
				}
				md, err := stream.Header()
				if err != nil {
					panic(err)
				}
				a.So(md.Get("x-rate-limit-limit")[0], should.Equal, "0")
				a.So(md.Get("x-rate-limit-reset")[0], should.Equal, "1")

				for i := uint64(0); i < 2*maxRate; i++ {
					err := stream.Send(&ttnpb.GatewayUp{})
					a.So(err, should.BeNil)
				}

				select {
				case <-server.ch["1"]:
					t.Fatal("Unexpected rate limiting error")
				case <-time.After(timeout):
				}
			})
		})
	})

	withRateLimitedGRPCServer(conf, func(server *mockServer, client ttnpb.GtwGsClient) {
		t.Run("RateLimit", func(t *testing.T) {
			t.Run("Invoke", func(t *testing.T) {
				a := assertions.New(t)
				for i := uint64(0); i < 2*maxRate; i++ {
					headers := metadata.MD{}
					_, err := client.GetConcentratorConfig(test.Context(), &pbtypes.Empty{}, grpc.Header(&headers))
					a.So(headers.Get("x-rate-limit-limit"), should.Resemble, []string{fmt.Sprintf("%d", maxRate)})
					a.So(headers.Get("x-rate-limit-reset"), should.Resemble, []string{fmt.Sprintf("%d", 1)})
					if i < maxRate {
						a.So(headers.Get("x-rate-limit-available"), should.Resemble, []string{fmt.Sprintf("%d", maxRate-(i+1))})
						a.So(err, should.BeNil)
					} else {
						a.So(headers.Get("x-rate-limit-available"), should.Resemble, []string{"0"})
						a.So(errors.IsResourceExhausted(err), should.BeTrue)
					}
				}
			})

			t.Run("OverrideMethod", func(t *testing.T) {
				a := assertions.New(t)
				headers := metadata.MD{}
				_, err := client.GetMQTTConnectionInfo(test.Context(), &ttnpb.GatewayIdentifiers{
					GatewayID: "gtw1",
				}, grpc.Header(&headers))
				a.So(err, should.BeNil)
				a.So(headers.Get("x-rate-limit-limit"), should.Resemble, []string{fmt.Sprintf("%d", overrideRate)})
				a.So(headers.Get("x-rate-limit-reset"), should.Resemble, []string{fmt.Sprintf("%d", 1)})
				a.So(headers.Get("x-rate-limit-available"), should.Resemble, []string{fmt.Sprintf("%d", overrideRate-1)})
			})

			t.Run("Stream", func(t *testing.T) {
				a := assertions.New(t)
				wg := sync.WaitGroup{}
				errCh := make(chan error, 1)
				defer close(errCh)

				for k := 0; k < 5; k++ {
					wg.Add(1)
					go func(k int) {
						ctx := metadata.NewOutgoingContext(test.Context(), metadata.Pairs("test-stream-id", fmt.Sprintf("%d", k)))
						stream, err := client.LinkGateway(ctx)
						if err != nil {
							panic(err)
						}
						md, err := stream.Header()
						if err != nil {
							panic(err)
						}
						a.So(md.Get("x-rate-limit-limit")[0], should.Equal, fmt.Sprintf("%d", maxRate))
						a.So(md.Get("x-rate-limit-reset")[0], should.Equal, "1")

						for i := uint64(0); i < 2*maxRate; i++ {
							err := stream.Send(&ttnpb.GatewayUp{})
							if err != nil && !errors.IsResourceExhausted(err) {
								panic(err)
							}
						}

						select {
						case count := <-server.ch[fmt.Sprintf("%d", k)]:
							a.So(count, should.Equal, maxRate)
						case <-time.After(timeout):
							errCh <- fmt.Errorf("Timed out waiting for rate limit exceeded error")
						}

						wg.Done()
					}(k)
				}

				select {
				case err := <-errCh:
					t.Fatal("Received unexpected error", err)
				case <-time.After(timeout):
				}

				wg.Wait()
			})

		})
	})
}
