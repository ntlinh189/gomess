package authprovider

type UserInfo struct {
	ID     string
	Email  string
	Name   string
	Avatar string
}

type ProviderInterface interface {
	Verify(token string) (*UserInfo, error)
}
