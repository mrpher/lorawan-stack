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
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// EchoKeyFunc calculates the rate limiting key from an echo request context.
type EchoKeyFunc func(c echo.Context) string

// EchoPathAndIP is an EchoKeyFunc that uses the path of the request URL and the real IP address as key.
func EchoPathAndIP(c echo.Context) string {
	return fmt.Sprintf("path:%s:from:%s", c.Request().URL.Path, c.RealIP())
}

// EchoMiddleware returns an Echo middleware function that rate limits HTTP requests.
func EchoMiddleware(c Config, keyFunc EchoKeyFunc, maxWait time.Duration) echo.MiddlewareFunc {
	r := c.New()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			md, ok := r.WaitMaxDuration(keyFunc(c), maxWait)
			if !ok {
				return errRateLimitExceeded.New()
			}
			c.Response().Header().Add("x-rate-limit-limit", strconv.FormatInt(md.Limit, 10))
			c.Response().Header().Add("x-rate-limit-available", strconv.FormatInt(md.Available, 10))
			c.Response().Header().Add("x-rate-limit-reset", strconv.FormatInt(md.ResetSeconds, 10))
			time.Sleep(md.Wait)
			return next(c)
		}
	}
}
