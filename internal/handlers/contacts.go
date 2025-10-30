package handlers

import (
	"errors"
	"net/http"

	"contactsAI/contacts/internal/config"
	"contactsAI/contacts/internal/db"
	"contactsAI/contacts/internal/routeutils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func RegisterContactsRoutes(router *gin.RouterGroup, env *config.Env) {
	apiGroup := router.Group("/contacts")
	apiGroup.GET("/", func(c *gin.Context) { GetContacts(c, env) })
	apiGroup.GET("/:id", func(c *gin.Context) { GetContactByID(c, env) })
	apiGroup.POST("/", func(c *gin.Context) { CreateContact(c, env) })
	apiGroup.PUT("/:id", func(c *gin.Context) { UpdateContact(c, env) })
	apiGroup.DELETE("/:id", func(c *gin.Context) { DeleteContact(c, env) })
}

type CreateContactBody struct {
	Name  string `json:"name"  binding:"required"`
	Phone string `json:"phone" binding:"required,phonenumber"`
}

// CreateContact godoc
//
//	@Summary		Create new contact
//	@Description	Create a new contact in the system
//	@Tags			contacts
//	@Accept			json
//	@Produce		json
//	@Param			contact	body		CreateContactBody	true	"Contact details"
//	@Success		201		{object}	ContactResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/contacts [post]
func CreateContact(c *gin.Context, env *config.Env) {
	var json CreateContactBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(err.Error()))
		return
	}

	contact := db.CreateContactParams{
		Name:  json.Name,
		Phone: json.Phone,
	}
	createdContact, err := env.CreateContact(c, contact)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse("Failed to create contact"))
		return
	}
	dto := toContactResponse(createdContact)
	c.JSON(http.StatusCreated, dto)
}

// GetContacts godoc
//
//	@Summary		Get all contacts
//	@Description	Retrieve all contacts from the database
//	@Tags			contacts
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]ContactResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/contacts [get]
func GetContacts(c *gin.Context, env *config.Env) {
	contacts, err := env.Queries.GetContacts(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(err.Error()))
		return
	}
	if contacts == nil {
		contacts = []db.Contact{}
	}

	dtos := make([]ContactResponse, len(contacts))
	for i, v := range contacts {
		dtos[i] = toContactResponse(v)
	}
	c.JSON(http.StatusOK, dtos)
}

// GetContactByID godoc
//
//	@Summary		Get contact by ID
//	@Description	Retrieve a single contact by its ID
//	@Tags			contacts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Contact ID"
//	@Success		200	{object}	ContactResponse
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/contacts/{id} [get]
func GetContactByID(c *gin.Context, env *config.Env) {
	contactID, err := routeutils.GetInt32FromPath(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid contact ID"))
		return
	}

	contact, err := env.GetContactByID(c, contactID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, NewErrorResponse("Contact not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(err.Error()))
		return
	}

	dto := toContactResponse(contact)

	c.JSON(http.StatusOK, dto)
}

type UpdateContactBody struct {
	Name  string `json:"name"  binding:"required"`
	Phone string `json:"phone" binding:"required,phonenumber"`
}

// UpdateContact godoc
//
//	@Summary		Update contact
//	@Description	Update an existing contact by ID
//	@Tags			contacts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Contact ID"
//	@Param			contact	body		UpdateContactBody	true	"Updated contact details"
//	@Success		200		{object}	ContactResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/contacts/{id} [put]
func UpdateContact(c *gin.Context, env *config.Env) {
	contactID, parseErr := routeutils.GetInt32FromPath(c, "id")
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid contact ID"))
		return
	}

	var json UpdateContactBody
	if bindErr := c.BindJSON(&json); bindErr != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(bindErr.Error()))
		return
	}
	contactParams := db.UpdateContactParams{
		ID:    contactID,
		Name:  json.Name,
		Phone: json.Phone,
	}

	contact, updateErr := env.UpdateContact(c, contactParams)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse("Updating contact failed."))
		return
	}
	dto := toContactResponse(contact)
	c.JSON(http.StatusOK, dto)
}

// DeleteContact godoc
//
//	@Summary		Delete contact
//	@Description	Delete a contact by ID
//	@Tags			contacts
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Contact ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/contacts/{id} [delete]
func DeleteContact(c *gin.Context, env *config.Env) {
	id, err := routeutils.GetInt32FromPath(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid contact ID"))
		return
	}

	if err = env.DeleteContact(c, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, NewErrorResponse("Contact not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
