package models

type User struct {
	ID         int64
	Provider   string
	ProviderID string
	Account    string
	Name       string
	Avatar     string
}
