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
