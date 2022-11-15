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

package util

import (
	"fmt"

	"tidbcloud-cli/internal/prop"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CheckAuth() func(*cobra.Command, []string) error {
	return func(*cobra.Command, []string) error {
		curP := viper.Get(prop.CurProfile)
		if curP == nil {
			return fmt.Errorf("no active profile for auth, please use `config init` to create one")
		}

		publicKey := viper.Get(fmt.Sprintf("%s.%s", curP, prop.PublicKey))
		if publicKey == nil {
			return fmt.Errorf("no public key configured for auth, please use `set %s <publicKey>` to set one", prop.PublicKey)
		}

		privateKey := viper.Get(fmt.Sprintf("%s.%s", curP, prop.PrivateKey))
		if privateKey == nil {
			return fmt.Errorf("no private key configured for auth, please use `set %s <privateKey>` to set one", prop.PrivateKey)
		}

		return nil
	}
}

func GetAccessKeys() (publicKey string, privateKey string) {
	curP := viper.Get(prop.CurProfile)
	publicKey = viper.Get(fmt.Sprintf("%s.%s", curP, prop.PublicKey)).(string)
	privateKey = viper.Get(fmt.Sprintf("%s.%s", curP, prop.PrivateKey)).(string)
	return
}
