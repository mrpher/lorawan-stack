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
	"testing"
	"time"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/ratelimit"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

func TestRateLimit(t *testing.T) {
	a := assertions.New(t)

	r := ratelimit.Config{
		Enable: true,
		Rate:   3,
	}.New()

	const (
		id1 = "id1"
		id2 = "id2"
		id3 = "id3"
	)

	t.Run("Wait", func(t *testing.T) {
		for i := int64(0); i < 3; i++ {
			for _, id := range []string{id1, id2} {
				a.So(r.Wait(id), should.Resemble, &ratelimit.Metadata{
					ResetSeconds: 1,
					Limit:        3,
					Available:    3 - (i + 1),
				})
			}
		}
	})

	t.Run("WaitMaxDuration", func(t *testing.T) {
		_, ok := r.WaitMaxDuration(id1, 100*time.Microsecond)
		a.So(ok, should.BeFalse)
		md, ok := r.WaitMaxDuration(id1, time.Second)
		a.So(ok, should.BeTrue)
		a.So(md.Wait, should.BeGreaterThan, time.Millisecond)
	})

	t.Run("WaitMaxDurationWithRate", func(t *testing.T) {
		md, ok := r.WaitMaxDurationWithRate(id3, time.Second, func() uint64 { return 100 })
		a.So(ok, should.BeTrue)
		a.So(md.Limit, should.Equal, 100)
	})
}
