package entities

type ContactTask struct {
	Action       string `json:"action"`
	ContactID    string `json:"contactId"`
	ContactName  string `json:"contactName"`
	ContactEmail string `json:"contactEmail"`
}

func NewContactTask(action, contactID, contactName, contactEmail string) ContactTask {
	return ContactTask{
		Action:       action,
		ContactID:    contactID,
		ContactName:  contactName,
		ContactEmail: contactEmail,
	}
}
