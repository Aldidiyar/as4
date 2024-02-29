package contact

import (
	"time"

	"arch/pkg/type/email"
	"arch/pkg/type/gender"
)

type ID struct {
	Value string `json:"id" uri:"id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
}

type ContactResponse struct {
	ID         string    `json:"id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	CreatedAt  time.Time `json:"createdAt"  binding:"required"`
	ModifiedAt time.Time `json:"modifiedAt"  binding:"required"`
	ShortContact
}

type ShortContact struct {
	PhoneNumber string        `json:"phoneNumber" binding:"required,max=50" maxLength:"50" example:"78002002020"`
	Email       email.Email   `json:"email" binding:"omitempty,max=250,email" maxLength:"250" example:"example@gmail.com" format:"email" swaggertype:"string"`
	Gender      gender.Gender `json:"gender" example:"1" enums:"1,2" swaggertype:"integer"`
	Age         uint8         `json:"age" binding:"min=0,max=200" minimum:"0" maximum:"200" default:"0" example:"42"`
	Name        string        `json:"name" binding:"max=50" maxLength:"50" example:"Иван"`
	Surname     string        `json:"surname" binding:"max=100" maxLength:"100" example:"Иванов"`
	Patronymic  string        `json:"patronymic" binding:"max=100" maxLength:"100" example:"Иванович"`
}

type ListContact struct {
	Total  uint64             `json:"total" example:"10" default:"0" binding:"min=0" minimum:"0"`
	Limit  uint64             `json:"limit"  example:"10" default:"10" binding:"min=0" minimum:"0"`
	Offset uint64             `json:"offset" example:"20" default:"0" binding:"min=0" minimum:"0"`
	List   []*ContactResponse `json:"list"`
}
