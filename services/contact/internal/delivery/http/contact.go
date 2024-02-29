package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"arch/pkg/tools/converter"
	"arch/pkg/type/context"
	"arch/pkg/type/logger"
	"arch/pkg/type/pagination"
	"arch/pkg/type/phoneNumber"
	"arch/pkg/type/query"
	"arch/pkg/type/queryParameter"
	jsonContact "arch/services/contact/internal/delivery/http/contact"
	domainContact "arch/services/contact/internal/domain/contact"
	"arch/services/contact/internal/domain/contact/age"
	"arch/services/contact/internal/domain/contact/name"
	"arch/services/contact/internal/domain/contact/patronymic"
	"arch/services/contact/internal/domain/contact/surname"
	"arch/services/contact/internal/useCase"
)

var mappingSortsContact = query.SortsOptions{
	"name":        {},
	"surname":     {},
	"patronymic":  {},
	"phoneNumber": {},
	"email":       {},
	"gender":      {},
	"age":         {},
}

func (d *Delivery) CreateContact(c *gin.Context) {

	var ctx = context.New(c)

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

	response, err := d.ucContact.Create(ctx, dContact)
	if err != nil {

		SetError(c, http.StatusInternalServerError, err)
		return
	}
	logger.InfoWithContext(ctx, "test log", zap.Any("Test", "Test"))
	if len(response) > 0 {
		c.JSON(http.StatusCreated, jsonContact.ToContactResponse(response[0]))
	} else {
		c.Status(http.StatusOK)
	}
}

func (d *Delivery) UpdateContact(c *gin.Context) {

	var ctx = context.New(c)

	var id jsonContact.ID
	if err := c.ShouldBindUri(&id); err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	contact := jsonContact.ShortContact{}
	if err := c.ShouldBindJSON(&contact); err != nil {
		SetError(c, http.StatusInternalServerError, err)
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

	var dContact, _ = domainContact.NewWithID(
		converter.StringToUUID(id.Value),
		time.Now().UTC(),
		time.Now().UTC(),
		*phoneNumber.New(contact.PhoneNumber),
		contact.Email,
		*contactName,
		*contactSurname,
		*contactPatronymic,
		*contactAge,
		contact.Gender,
	)

	response, err := d.ucContact.Update(ctx, *dContact)
	if err != nil {
		if errors.Is(err, useCase.ErrContactNotFound) {
			SetError(c, http.StatusNotFound, err)
			return
		}

		SetError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, jsonContact.ToContactResponse(response))

}

func (d *Delivery) DeleteContact(c *gin.Context) {

	var ctx = context.New(c)

	var id jsonContact.ID
	if err := c.ShouldBindUri(&id); err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	if err := d.ucContact.Delete(ctx, converter.StringToUUID(id.Value)); err != nil {
		SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (d *Delivery) ListContact(c *gin.Context) {

	var ctx = context.New(c)
	params, err := query.ParseQuery(c, query.Options{
		Sorts: mappingSortsContact,
	})

	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	contacts, err := d.ucContact.List(ctx, queryParameter.QueryParameter{
		Sorts: params.Sorts,
		Pagination: pagination.Pagination{
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	})
	if err != nil {
		SetError(c, http.StatusInternalServerError, err)
		return
	}

	count, err := d.ucContact.Count(ctx)
	if err != nil {
		SetError(c, http.StatusInternalServerError, err)
		return
	}

	var result = jsonContact.ListContact{
		Total:  count,
		Limit:  params.Limit,
		Offset: params.Offset,
		List:   []*jsonContact.ContactResponse{},
	}
	for _, value := range contacts {
		result.List = append(result.List, jsonContact.ToContactResponse(value))
	}

	c.JSON(http.StatusOK, result)
}

func (d *Delivery) ReadContactByID(c *gin.Context) {

	var ctx = context.New(c)

	var id jsonContact.ID
	if err := c.ShouldBindUri(&id); err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	response, err := d.ucContact.ReadByID(ctx, converter.StringToUUID(id.Value))
	if err != nil {
		if errors.Is(err, useCase.ErrContactNotFound) {
			SetError(c, http.StatusNotFound, err)
			return
		}

		SetError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, jsonContact.ToContactResponse(response))

}
