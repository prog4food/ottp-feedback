package main

type DummyChat struct {
	ID string `json:"id"`
}

func (u *DummyChat) Recipient() string {
	return u.ID
}
