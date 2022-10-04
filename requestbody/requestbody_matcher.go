// Copyright 2017 Google Inc. All rights reserved.
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

package body

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	uuid "github.com/google/uuid"
)

// Matcher is a conditonal evalutor of query string parameters
// to be used in structs that take conditions.
type Matcher struct {
	value     string
	val_cache map[string]bool
}

// NewMatcher builds a new body matcher
func NewMatcher(value string) *Matcher {
	return &Matcher{value: value}
}

// MatchRequest evaluates a request and returns whether or not
// the request contains a querystring param that matches the provided name
// and value.
func (m *Matcher) MatchRequest(req *http.Request) bool {
	if req.Body == nil {
		return false
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return false
	}
	req.Body.Close()

	// And now set a new body, which will simulate the same data we read:
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	req_uuid := req.Header.Get("X-Request-UUID")
	if req_uuid == "" {
		new_uuid := uuid.New()
		req_uuid = new_uuid.String()
		req.Header.Set("X-Request-UUID", req_uuid)
	}

	if m.val_cache == nil {
		m.val_cache = map[string]bool{}
	}

	m.val_cache[req_uuid] = strings.Contains(string(body), m.value)
	return m.val_cache[req_uuid]
}

// MatchResponse evaluates a response and returns whether or not
// the request that resulted in that response has a body which contains the value
func (m *Matcher) MatchResponse(res *http.Response) bool {
	req_uuid := res.Request.Header.Get("X-Request-UUID")
	if req_uuid == "" {
		return false
	}

	if ret, ok := m.val_cache[req_uuid]; ok {
		delete(m.val_cache, req_uuid)
		return ret
	}
	return false
}
