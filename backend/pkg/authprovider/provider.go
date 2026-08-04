package authprovider

type UserInfo struct {
	ID      string
	Account string
	Name    string
	Avatar  string
}

type ProviderInterface interface {
	Verify(token string) (*UserInfo, error)
}
