// Copyright 2022 PingCAP, Inc.
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

package prop

const (
	PublicKey  string = "public_key"
	PrivateKey string = "private_key"
	CurProfile string = "current_profile"
)

func GlobalProperties() []string {
	return []string{CurProfile}
}

func ProfileProperties() []string {
	return []string{PublicKey, PrivateKey}
}

func Properties() []string {
	return []string{PublicKey, PrivateKey, CurProfile}
}
