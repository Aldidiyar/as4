package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"arch/pkg/tools/converter"
	"arch/pkg/type/context"
	"arch/pkg/type/phoneNumber"
	jsonContact "arch/services/contact/internal/delivery/http/contact"
	jsonGroup "arch/services/contact/internal/delivery/http/group"
	domainContact "arch/services/contact/internal/domain/contact"
	"arch/services/contact/internal/domain/contact/age"
	"arch/services/contact/internal/domain/contact/name"
	"arch/services/contact/internal/domain/contact/patronymic"
	"arch/services/contact/internal/domain/contact/surname"
)

func (d *Delivery) CreateContactIntoGroup(c *gin.Context) {

	var ctx = context.New(c)

	var id jsonGroup.ID
	if err := c.ShouldBindUri(&id); err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	contact := jsonContact.ShortContact{}
	if err := c.ShouldBindJSON(&contact); err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	contactAge, err := age.New(contact.Age)
	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	contactName, err := name.New(contact.Name)
	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	contactSurname, err := surname.New(contact.Surname)
	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	contactPatronymic, err := patronymic.New(contact.Patronymic)
	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	dContact, err := domainContact.New(
		*phoneNumber.New(contact.PhoneNumber),
		contact.Email,
		*contactName,
		*contactSurname,
		*contactPatronymic,
		*contactAge,
		contact.Gender,
	)
	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	contacts, err := d.ucGroup.CreateContactIntoGroup(ctx, converter.StringToUUID(id.Value), dContact)
	if err != nil {
		SetError(c, http.StatusInternalServerError, err)
		return
	}

	var list = []*jsonContact.ContactResponse{}
	for _, value := range contacts {
		list = append(list, jsonContact.ToContactResponse(value))
	}

	c.JSON(http.StatusOK, list)

}

func (d *Delivery) AddContactToGroup(c *gin.Context) {

	var ctx = context.New(c)

	var id jsonGroup.ID
	if err := c.ShouldBindUri(&id); err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	var contactID jsonGroup.ContactID
	if err := c.ShouldBindUri(&contactID); err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	if err := d.ucGroup.AddContactToGroup(ctx, converter.StringToUUID(id.Value), converter.StringToUUID(contactID.Value)); err != nil {
		SetError(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)

}

func (d *Delivery) DeleteContactFromGroup(c *gin.Context) {

	var ctx = context.New(c)

	var id jsonGroup.ID
	if err := c.ShouldBindUri(&id); err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	var contactID jsonGroup.ContactID
	if err := c.ShouldBindUri(&contactID); err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	if err := d.ucGroup.DeleteContactFromGroup(ctx, converter.StringToUUID(id.Value), converter.StringToUUID(contactID.Value)); err != nil {
		SetError(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
