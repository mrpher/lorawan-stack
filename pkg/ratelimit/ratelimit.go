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

package ratelimit

import (
	"sync"
	"time"

	"github.com/juju/ratelimit"
)

// RateLimiter can be used to rate limit access to a resource (using an id string).
type RateLimiter interface {
	Wait(id string) *Metadata
	WaitMaxDuration(id string, maxDuration time.Duration) (*Metadata, bool)
	WaitMaxDurationWithRate(id string, maxDuration time.Duration, f func() uint64) (*Metadata, bool)
}

// noopRateLimiter does not enforce any rate limits.
type noopRateLimiter struct{}

// Wait implements the RateLimiter interface.
func (*noopRateLimiter) Wait(string) *Metadata {
	return &Metadata{}
}

// WaitMaxDuration implements the RateLimiter interface.
func (*noopRateLimiter) WaitMaxDuration(string, time.Duration) (*Metadata, bool) {
	return &Metadata{}, true
}

// WaitMaxDurationWithRate implements the RateLimiter interface.
func (*noopRateLimiter) WaitMaxDurationWithRate(string, time.Duration, func() uint64) (*Metadata, bool) {
	return &Metadata{}, true
}

// Registry for rate limiting.
type Registry struct {
	mu           sync.RWMutex
	rate         uint64
	per          time.Duration
	resetSeconds int64
	entities     map[string]*ratelimit.Bucket
}

// New returns a new RateLimiter from configuration.
func (c Config) New() RateLimiter {
	if !c.Enable {
		return &noopRateLimiter{}
	}
	return &Registry{
		rate:         uint64(c.Rate),
		per:          time.Second,
		resetSeconds: 1,
		entities:     make(map[string]*ratelimit.Bucket),
	}
}

func (r *Registry) getOrCreate(id string, createFunc func(uint64) *ratelimit.Bucket, rateFunc func() uint64) *ratelimit.Bucket {
	r.mu.RLock()
	limiter, ok := r.entities[id]
	r.mu.RUnlock()
	if ok {
		return limiter
	}
	limiter = createFunc(rateFunc())
	// TODO: this may lead to a race condition and leak buckets
	// TODO: garbage collection. maybe delete old buckets from the map instead of resetting.
	r.mu.Lock()
	r.entities[id] = limiter
	r.mu.Unlock()
	return limiter
}

func (r *Registry) createFunc(rate uint64) *ratelimit.Bucket {
	return ratelimit.NewBucketWithQuantum(r.per, int64(rate), int64(rate))
}

// Wait returns the time to wait until available
func (r *Registry) Wait(id string) *Metadata {
	b := r.getOrCreate(id, r.createFunc, func() uint64 { return r.rate })
	t := b.Take(1)

	return &Metadata{
		Wait:         t,
		ResetSeconds: r.resetSeconds,
		Available:    b.Available(),
		Limit:        int64(b.Rate()),
	}
}

// WaitMaxDurationWithRate returns the time to wait until available, but with a maximum duration to wait.
func (r *Registry) WaitMaxDurationWithRate(id string, max time.Duration, rate func() uint64) (*Metadata, bool) {
	b := r.getOrCreate(id, r.createFunc, rate)
	t, ok := b.TakeMaxDuration(1, max)
	if !ok {
		return &Metadata{
			ResetSeconds: r.resetSeconds,
			Limit:        int64(b.Rate()),
		}, false
	}
	return &Metadata{
		Wait:         t,
		ResetSeconds: r.resetSeconds,
		Available:    b.Available(),
		Limit:        int64(b.Rate()),
	}, true
}

// WaitMaxDuration returns the time to wait until available, but with a maximum duration to wait. The default rate limit is enforced.
func (r *Registry) WaitMaxDuration(id string, max time.Duration) (*Metadata, bool) {
	return r.WaitMaxDurationWithRate(id, max, func() uint64 { return r.rate })
}
