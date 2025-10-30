package handlers

import "contactsAI/contacts/internal/db"

type ContactResponse struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	OwnerID *int32 `json:"owner_id,omitempty"`
}

func toContactResponse(contact db.Contact) ContactResponse {
	var ownerID *int32
	if contact.OwnerID.Valid {
		ownerID = &contact.OwnerID.Int32
	}
	return ContactResponse{
		ID:      contact.ID,
		Name:    contact.Name,
		Phone:   contact.Phone,
		OwnerID: ownerID,
	}
}
