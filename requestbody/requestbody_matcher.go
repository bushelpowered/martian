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
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const X_HASH_HEADER = "X-Request-Hash"

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

	req_hash := req.Header.Get(X_HASH_HEADER)
	if req_hash == "" {
		req_hash = hashBody(string(body))
		req.Header.Set(X_HASH_HEADER, req_hash)
	}

	if m.val_cache == nil {
		m.val_cache = map[string]bool{}
	}

	if _, ok := m.val_cache[req_hash]; ok && os.Getenv("ONLY_UNQIUE_BODY") == "true" {
		m.val_cache[req_hash] = false
		return m.val_cache[req_hash]
	}

	m.val_cache[req_hash] = strings.Contains(string(body), m.value)
	return m.val_cache[req_hash]
}

func hashBody(body string) string {
	s := sha256.New()
	return base64.URLEncoding.EncodeToString(s.Sum([]byte(body)))
}

// MatchResponse evaluates a response and returns whether or not
// the request that resulted in that response has a body which contains the value
func (m *Matcher) MatchResponse(res *http.Response) bool {
	req_hash := res.Request.Header.Get(X_HASH_HEADER)
	if req_hash == "" {
		return false
	}

	if ret, ok := m.val_cache[req_hash]; ok {
		return ret
	}
	return false
}
