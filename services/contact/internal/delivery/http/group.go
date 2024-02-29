package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"arch/pkg/tools/converter"
	"arch/pkg/type/context"
	"arch/pkg/type/pagination"
	"arch/pkg/type/query"
	"arch/pkg/type/queryParameter"
	jsonGroup "arch/services/contact/internal/delivery/http/group"
	domainGroup "arch/services/contact/internal/domain/group"
	"arch/services/contact/internal/domain/group/description"
	"arch/services/contact/internal/domain/group/name"
	"arch/services/contact/internal/useCase"
)

var mappingSortsGroup = query.SortsOptions{
	"id":           {},
	"name":         {},
	"description":  {},
	"contactCount": {},
}

func (d *Delivery) CreateGroup(c *gin.Context) {

	var ctx = context.New(c)

	var group = &jsonGroup.ShortGroup{}

	if err := c.ShouldBindJSON(&group); err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	groupName, err := name.New(group.Name)
	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}
	groupDescription, err := description.New(group.Description)
	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}
	newGroup, err := d.ucGroup.Create(ctx, domainGroup.New(
		groupName,
		groupDescription,
	))
	if err != nil {
		SetError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, jsonGroup.GroupResponse{
		ID:         newGroup.ID().String(),
		CreatedAt:  newGroup.CreatedAt(),
		ModifiedAt: newGroup.ModifiedAt(),
		Group: jsonGroup.Group{
			ShortGroup: jsonGroup.ShortGroup{
				Name:        newGroup.Name().Value(),
				Description: newGroup.Description().Value(),
			},
			ContactsAmount: newGroup.ContactCount(),
		},
	})
}

func (d *Delivery) UpdateGroup(c *gin.Context) {

	var ctx = context.New(c)

	var id jsonGroup.ID
	if err := c.ShouldBindUri(&id); err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	group := jsonGroup.ShortGroup{}
	if err := c.ShouldBindJSON(&group); err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	groupName, err := name.New(group.Name)
	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}
	groupDescription, err := description.New(group.Description)
	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	response, err := d.ucGroup.Update(ctx, domainGroup.NewWithID(
		converter.StringToUUID(id.Value),
		time.Now().UTC(),
		time.Now().UTC(),
		groupName,
		groupDescription,
		0,
	))
	if err != nil {
		if errors.Is(err, useCase.ErrGroupNotFound) {
			SetError(c, http.StatusNotFound, err)
			return
		}

		SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, jsonGroup.ProtoToGroupResponse(response))
}

func (d *Delivery) DeleteGroup(c *gin.Context) {

	var ctx = context.New(c)

	var id jsonGroup.ID
	if err := c.ShouldBindUri(&id); err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	if err := d.ucGroup.Delete(ctx, converter.StringToUUID(id.Value)); err != nil {
		SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (d *Delivery) ListGroup(c *gin.Context) {

	var ctx = context.New(c)

	params, err := query.ParseQuery(c, query.Options{
		Sorts: mappingSortsGroup,
	})

	if err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	groups, err := d.ucGroup.List(ctx, queryParameter.QueryParameter{
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

	var list = make([]*jsonGroup.GroupResponse, len(groups))

	for i, elem := range groups {
		list[i] = jsonGroup.ProtoToGroupResponse(elem)
	}

	c.JSON(http.StatusOK, jsonGroup.GroupList{
		Total:  count,
		Limit:  params.Limit,
		Offset: params.Offset,
		List:   list,
	})
}

func (d *Delivery) ReadGroupByID(c *gin.Context) {

	var ctx = context.New(c)

	var id jsonGroup.ID
	if err := c.ShouldBindUri(&id); err != nil {
		SetError(c, http.StatusBadRequest, err)
		return
	}

	response, err := d.ucGroup.ReadByID(ctx, converter.StringToUUID(id.Value))
	if err != nil {
		if errors.Is(err, useCase.ErrGroupNotFound) {
			SetError(c, http.StatusNotFound, err)
			return
		}

		SetError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, jsonGroup.ProtoToGroupResponse(response))
}
