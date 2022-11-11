package util

import (
	"fmt"

	"tidbcloud-cli/internal/prop"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CheckAuth(cmd *cobra.Command, args []string) error {
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

func GetAccessKeys() (publicKey string, privateKey string) {
	curP := viper.Get(prop.CurProfile)
	publicKey = viper.Get(fmt.Sprintf("%s.%s", curP, prop.PublicKey)).(string)
	privateKey = viper.Get(fmt.Sprintf("%s.%s", curP, prop.PrivateKey)).(string)
	return
}
