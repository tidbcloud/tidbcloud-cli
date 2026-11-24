// Copyright 2026 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cloud

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSeenPageTokens(t *testing.T) {
	seen := newSeenPageTokens()

	// A fresh token passes and is remembered.
	first := "page-1"
	require.NoError(t, seen.check(&first, "cluster"))

	// Empty pagination end always passes and is never remembered.
	require.NoError(t, seen.check(nil, "cluster"))
	empty := ""
	require.NoError(t, seen.check(&empty, "cluster"))

	// A distinct token passes.
	second := "page-2"
	require.NoError(t, seen.check(&second, "cluster"))

	// Repeating an already seen token is rejected, so pagination loops
	// terminate instead of fetching the same page forever.
	repeated := "page-1"
	err := seen.check(&repeated, "cluster")
	require.EqualError(t, err, "cluster list returned a repeated next page token")

	// The empty token still passes after repeats.
	require.NoError(t, seen.check(nil, "cluster"))
}
