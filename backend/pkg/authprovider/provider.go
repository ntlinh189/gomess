package authprovider

type UserInfo struct {
	ID     string
	Email  string
	Name   string
	Avatar string
}

type Provider interface {
	Verify(token string) (*UserInfo, error)
}
