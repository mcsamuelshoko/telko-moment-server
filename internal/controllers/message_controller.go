package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
	"github.com/mcsamuelshoko/telko-moment-server/internal/services"
	"github.com/mcsamuelshoko/telko-moment-server/pkg/utils"
	"github.com/rs/zerolog"
	"time"
)

type IMessageController interface {
	// CreateMessage Create a new message
	// (POST /messages)
	CreateMessage(c *fiber.Ctx) error

	// GetAllMessages Get all messages
	// (GET /messages)
	GetAllMessages(c *fiber.Ctx) error

	// GetMessageById Get a message by ID
	// (GET /messages/{messageId})
	GetMessageById(c *fiber.Ctx, messageId string) error

	// GetMessagesByUserId Get a message by ID
	// (GET /messages/user/{userId})
	GetMessagesByUserId(c *fiber.Ctx, userId string) error

	// UpdateMessage Update a message
	// (PUT /messages/{messageId})
	UpdateMessage(c *fiber.Ctx, messageId string) error

	// DeleteMessage Delete a message
	// (DELETE /messages/{messageId})
	DeleteMessage(c *fiber.Ctx, messageId string) error
}

type MessageController struct {
	iName          string
	logger         *zerolog.Logger
	messageService services.IMessageService
}

func NewMessageController(log *zerolog.Logger, messageSvc services.IMessageService) IMessageController {
	return &MessageController{
		iName:          "MessageController",
		logger:         log,
		messageService: messageSvc,
	}
}

func (m MessageController) CreateMessage(c *fiber.Ctx) error {
	const kName = "CreateMessage"
	message := new(models.Message)
	err := c.BodyParser(message)
	if err != nil {
		m.logger.Error().Interface(kName, m.iName).Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}
	//add properties to message
	message.CreatedAt = time.Now()

	//create message via service
	createdMsg, err := m.messageService.Create(c.Context(), message)
	if err != nil {
		m.logger.Error().Interface(kName, m.iName).Err(err).Msg("Failed to create message")
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to create message"))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.SuccessResponse(createdMsg, "Created message"))
}

func (m MessageController) GetAllMessages(c *fiber.Ctx) error {
	const kName = "GetAllMessages"

	return nil
}

func (m MessageController) GetMessageById(c *fiber.Ctx, messageId string) error {
	const kName = "GetMessageById"

	message, err := m.messageService.GetById(c.Context(), messageId)
	if err != nil {
		m.logger.Error().Interface(kName, m.iName).Err(err).Msg("Failed to get message")
		return c.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse("Failed to get message"))
	}
	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(message, "Message Found"))
}

func (m MessageController) GetMessagesByUserId(c *fiber.Ctx, userId string) error {
	const kName = "GetMessagesByUserId"

	senderMsgs, err := m.messageService.GetBySenderId(c.Context(), userId)
	if err != nil {
		m.logger.Error().Interface(kName, m.iName).Err(err).Msg("Failed to get sender id")
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to get sender id"))
	}
	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(senderMsgs, "Messages Found"))
}

func (m MessageController) UpdateMessage(c *fiber.Ctx, messageId string) error {
	const kName = "UpdateMessage"

	message := new(models.Message)
	err := c.BodyParser(message)
	if err != nil {
		m.logger.Error().Interface(kName, m.iName).Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse("Invalid request body"))
	}

	message.UpdatedAt = time.Now()
	err = m.messageService.Update(c.Context(), message)
	if err != nil {
		m.logger.Error().Interface(kName, m.iName).Err(err).Msg("Failed to update message")
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to update message"))
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(message, "Updated message"))
}

func (m MessageController) DeleteMessage(c *fiber.Ctx, messageId string) error {
	const kName = "DeleteMessage"

	err := m.messageService.Delete(c.Context(), messageId)
	if err != nil {
		m.logger.Error().Interface(kName, m.iName).Err(err).Msg("Failed to delete message")
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse("Failed to delete message"))
	}
	return c.Status(fiber.StatusAccepted).JSON(utils.SuccessResponse(nil, "Message Deleted"))
}
