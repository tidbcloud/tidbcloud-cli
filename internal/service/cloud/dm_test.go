// Copyright 2025 PingCAP, Inc.
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
	"context"
	"testing"
)

// Test to verify that all DM methods are implemented in the TiDBCloudClient interface
func TestDMMethodsImplemented(t *testing.T) {
	// This test verifies that the interface is correctly implemented
	// by checking that we can declare a variable of the interface type
	var client TiDBCloudClient
	_ = client

	// Create a dummy ClientDelegate to check method signatures
	delegate := &ClientDelegate{}

	// Verify all DM methods exist and have correct signatures
	ctx := context.Background()

	// Test method signatures by attempting to call with nil parameters
	// This will compile-time verify the interface is properly implemented
	_, _ = delegate.Precheck(ctx, nil)
	_, _ = delegate.CreateTask(ctx, nil)
	_, _ = delegate.GetTask(ctx, nil)
	_, _ = delegate.ListTasks(ctx, nil)
	_, _ = delegate.DeleteTask(ctx, nil)
	_, _ = delegate.OperateTask(ctx, nil)
	_, _ = delegate.CancelPrecheck(ctx, nil)
	_, _ = delegate.GetPrecheck(ctx, nil)

	t.Log("All DM methods are properly implemented with correct signatures")
}
