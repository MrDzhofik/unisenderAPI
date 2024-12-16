package entities

type Unisender struct {
	UnisenderKey string `json:"unisender_key"`
	AccountID    string `json:"account_id"`
}

func NewUnisender(unisenderKey string, accountID string) *Unisender {
	return &Unisender{
		UnisenderKey: unisenderKey,
		AccountID:    accountID,
	}
}
