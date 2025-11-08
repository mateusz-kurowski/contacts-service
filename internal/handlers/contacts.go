package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"contactsAI/contacts/internal/config"
	"contactsAI/contacts/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func RegisterContactsRoutes(router *gin.RouterGroup, env *config.Env) {
	apiGroup := router.Group("/contacts")
	apiGroup.GET("/", func(c *gin.Context) { GetContacts(c, env) })
	apiGroup.GET("/:id", func(c *gin.Context) { GetContactByID(c, env) })
	apiGroup.POST("/", func(c *gin.Context) { CreateContact(c, env) })
	apiGroup.PUT("/:id", func(c *gin.Context) { UpdateContact(c, env) })
	apiGroup.PUT("/:id/avatar", func(c *gin.Context) { UploadContactAvatar(c, env) })
	apiGroup.GET("/:id/avatar", func(c *gin.Context) { DownloadContactAvatar(c, env) })
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
//	@Success		200	{array}		ContactResponse
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
	contactID, err := getIntFromPath(c, "id")
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
	contactID, parseErr := getIntFromPath(c, "id")
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

const (
	BytesPerKB = 1024
	KBPerMB    = 1024
	MaxMBSize  = 10
)

const maxAvatarSize = int64(MaxMBSize * KBPerMB * BytesPerKB) // 10 MiB

// UploadContactAvatar godoc
//
//	@Summary		Upload contact avatar
//	@Description	Upload an avatar image for a contact by ID
//	@Tags			contacts
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id		path		int		true	"Contact ID"
//	@Param			avatar	formData	file	true	"Avatar file"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/contacts/{id}/avatar [put]
func UploadContactAvatar(c *gin.Context, env *config.Env) {
	avatar, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse("Invalid avatar file provided"))
		return
	}
	if avatar.Size == 0 {
		c.JSON(http.StatusBadRequest, NewErrorResponse("Avatar file is empty"))
		return
	}
	if avatar.Size > maxAvatarSize {
		c.JSON(http.StatusBadRequest, NewErrorResponse(fmt.Sprintf("Avatar size cannot exceed %dMB", MaxMBSize)))
		return
	}

	f, err := avatar.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse("Failed to process avatar file"))
		return
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse("Failed to read avatar file"))
		return
	}

	key := "test" // TODO: make unique per contact
	contentType := avatar.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	if err = env.Bucket.Upload(c.Request.Context(), key, data, contentType); err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse("Could not upload avatar"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Avatar uploaded"})
}

// DownloadContactAvatar godoc
//
//	@Summary		Download contact's avatar
//	@Description	Streams a contact's avatar by contact ID
//	@Tags			contacts
//	@Produce		octet-stream
//	@Param			id	path		int		true	"Contact ID"
//	@Success		200	{file}		file	"The avatar file stream"
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/contacts/{id}/avatar [get]
func DownloadContactAvatar(c *gin.Context, env *config.Env) {
	objectKey := c.Param("id")

	s3Object, err := env.Bucket.GetStream(c.Request.Context(), objectKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Avatar not found."})
		return
	}
	defer s3Object.Body.Close()

	if s3Object.ContentType != nil {
		c.Header("Content-Type", *s3Object.ContentType)
	} else {
		c.Header("Content-Type", "application/octet-stream")
	}

	if s3Object.ContentLength != nil {
		c.Header("Content-Length", strconv.FormatInt(*s3Object.ContentLength, 10))
	}

	c.Header("Content-Disposition", "inline")

	_, err = io.Copy(c.Writer, s3Object.Body)
	if err != nil {
		env.Logger.Error("Error streaming file to client", "error", err)
	}
}

// DeleteContact godoc
//
//	@Summary		Delete contact
//	@Description	Delete a contact by ID
//	@Tags			contacts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int		true	"Contact ID"
//	@Success		204	{string}	string	"No Content"
//	@Failure		400	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/contacts/{id} [delete]
func DeleteContact(c *gin.Context, env *config.Env) {
	id, err := getIntFromPath(c, "id")
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
