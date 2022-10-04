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

// MatchRequest always returns false because the response body does not exist yet
func (m *Matcher) MatchRequest(req *http.Request) bool {
	return false
}

func (m *Matcher) MatchResponse(res *http.Response) bool {
	if res.Body == nil {
		return false
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return false
	}
	res.Body.Close()

	// And now set a new body, which will simulate the same data we read:
	res.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return strings.Contains(string(body), m.value)
}
