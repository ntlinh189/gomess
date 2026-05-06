package model

type User struct {
	ID         int64
	Provider   string
	ProviderID string
	Email      string
	Name       string
	Avatar     string
}
