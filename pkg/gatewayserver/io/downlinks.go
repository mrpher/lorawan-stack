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

package io

import (
	"strings"
	"sync"

	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
)

// DownlinksRegistry stores downlink messages, that can later be retrieved using the correlation IDs.
type DownlinksRegistry interface {
	// Push stores a new downlink message
	Push(*ttnpb.DownlinkMessage) bool
	// Pop retrieves the downlink message with the specified correlation IDs, if any.
	Pop(cids []string) (msg *ttnpb.DownlinkMessage, exists bool)
}

// NewMemoryDownlinksRegistry creates a new in-memory registry for downlink messages.
func NewMemoryDownlinksRegistry() DownlinksRegistry {
	return &memoryDownlinks{
		downlinks: make(map[string]*ttnpb.DownlinkMessage),
	}
}

// TODO: this is a temporary implementation, replace with a decent implementation by addressing TODOs below.
type memoryDownlinks struct {
	downlinksMu sync.Mutex
	downlinks   map[string]*ttnpb.DownlinkMessage
}

func concat(cids []string) string {
	// TODO: cids order should not matter
	return strings.Join(cids, "/")
}

func (d *memoryDownlinks) Push(msg *ttnpb.DownlinkMessage) bool {
	d.downlinksMu.Lock()
	defer d.downlinksMu.Unlock()

	// TODO: check existing, garbage collection
	key := concat(msg.GetCorrelationIDs())
	if key == "" {
		return false
	}
	d.downlinks[key] = msg
	return true
}

func (d *memoryDownlinks) Pop(cids []string) (*ttnpb.DownlinkMessage, bool) {
	d.downlinksMu.Lock()
	defer d.downlinksMu.Unlock()

	key := concat(cids)
	msg, ok := d.downlinks[key]
	if ok {
		delete(d.downlinks, key)
	}
	return msg, ok
}
