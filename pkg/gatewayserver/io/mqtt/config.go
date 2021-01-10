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

package mqtt

import (
	"go.thethings.network/lorawan-stack/v3/pkg/config"
	"go.thethings.network/lorawan-stack/v3/pkg/ratelimit"
)

// RateLimitingConfig represents rate limiting configuration for the MQTT frontend.
type RateLimitingConfig struct {
	Connections ratelimit.Config `name:"connections" description:"Rate limit new connections"`
	Traffic     ratelimit.Config `name:"traffic" description:"Rate limit gateway traffic"`
}

// Config represents configuration for the MQTT frontend.
type Config struct {
	config.MQTT  `name:",squash"`
	RateLimiting RateLimitingConfig `name:"rate-limiting"`
}
