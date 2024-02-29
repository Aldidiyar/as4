package group

import "time"

type GroupResponse struct {
	ID         string    `json:"id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	CreatedAt  time.Time `json:"createdAt"  binding:"required"`
	ModifiedAt time.Time `json:"modifiedAt"  binding:"required"`
	Group
}

type Group struct {
	ShortGroup
	ContactsAmount uint64 `json:"contactsAmount" default:"10" binding:"min=0" minimum:"0"`
}

type ShortGroup struct {
	Name        string `json:"name" binding:"required,max=100" example:"Название группы" maxLength:"100"`
	Description string `json:"description" example:"Описание группы" binding:"max=1000" maxLength:"1000"`
}

type GroupList struct {
	Total  uint64           `json:"total" example:"10" default:"0" binding:"min=0" minimum:"0"`
	Limit  uint64           `json:"limit"  example:"10" default:"10" binding:"min=0" minimum:"0"`
	Offset uint64           `json:"offset" example:"20" default:"0" binding:"min=0" minimum:"0"`
	List   []*GroupResponse `json:"list" binding:"min=0" minimum:"0"`
}

type ID struct {
	Value string `json:"id" uri:"id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
}

type ContactID struct {
	Value string `json:"id" uri:"contactId" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
}
