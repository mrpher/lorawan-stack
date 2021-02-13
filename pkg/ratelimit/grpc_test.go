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
	"testing"

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

const (
	maxRate      = 10
	overrideRate = 15
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
	}

	withRateLimitedGRPCServer(rpcserver.RateLimitingConfig{}, func(client ttnpb.GtwGsClient) {
		t.Run("NoRateLimit", func(t *testing.T) {
			t.Run("Invoke", func(t *testing.T) {
				a := assertions.New(t)
				for i := 0; i < 2*maxRate; i++ {
					headers := metadata.MD{}
					_, err := client.GetConcentratorConfig(test.Context(), &pbtypes.Empty{}, grpc.Header(&headers))
					a.So(err, should.BeNil)
					a.So(headers.Get("x-rate-limit-limit"), should.Resemble, []string{"0"})
					a.So(headers.Get("x-rate-limit-reset"), should.Resemble, []string{"0"})
					a.So(headers.Get("x-rate-limit-available"), should.Resemble, []string{"0"})
				}
			})

			t.Run("Stream", func(t *testing.T) {
				// TODO
			})
		})
	})

	withRateLimitedGRPCServer(conf, func(client ttnpb.GtwGsClient) {
		t.Run("RateLimit", func(t *testing.T) {
			t.Run("Invoke", func(t *testing.T) {
				a := assertions.New(t)
				for i := 0; i < 2*maxRate; i++ {
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
				// TODO
			})
		})
	})
}
