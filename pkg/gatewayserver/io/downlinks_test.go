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

package io_test

import (
	"testing"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/gatewayserver/io"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

func TestDownlinks(t *testing.T) {
	d := io.NewMemoryDownlinksRegistry()

	msgs := []*ttnpb.DownlinkMessage{
		{
			CorrelationIDs: []string{"corr1", "corr2"},
			RawPayload:     []byte{0, 1},
		},
		{
			CorrelationIDs: []string{"corr1", "corr3"},
			RawPayload:     []byte{0, 1, 2},
		},
	}
	for _, m := range msgs {
		d.Push(m)
	}

	a := assertions.New(t)
	msg, ok := d.Pop([]string{"unknown"})
	a.So(ok, should.BeFalse)
	a.So(msg, should.BeNil)

	msg, ok = d.Pop([]string{"corr1", "corr2"})
	a.So(ok, should.BeTrue)
	a.So(msg, should.Resemble, msgs[0])

	msg, ok = d.Pop([]string{"corr1", "corr3"})
	a.So(ok, should.BeTrue)
	a.So(msg, should.Resemble, msgs[1])

	msg, ok = d.Pop([]string{"corr1", "corr3"})
	a.So(ok, should.BeFalse)
	a.So(msg, should.BeNil)

	msg, ok = d.Pop(nil)
	a.So(ok, should.BeFalse)
	a.So(msg, should.BeNil)
}
