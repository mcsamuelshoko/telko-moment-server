package services

import (
	"context"
	"github.com/mcsamuelshoko/telko-moment-server/internal/models"
	"github.com/mcsamuelshoko/telko-moment-server/internal/repository"
)

type chatGroupService interface {
	Create(ctx context.Context, chatGroup *models.ChatGroup) error
	GetById(ctx context.Context, chatGroupId string) (*models.ChatGroup, error)
	List(ctx context.Context, page, limit int) ([]models.ChatGroup, error)
	Update(ctx context.Context, chatGroup *models.ChatGroup) error
	Delete(ctx context.Context, chatGroupId string) error
}

type ChatGroupService struct {
	repo repository.ChatGroupRepository
}

func NewChatGroupService(repo repository.ChatGroupRepository) *ChatGroupService {
	return &ChatGroupService{
		repo: repo,
	}
}

func (c ChatGroupService) Create(ctx context.Context, chatGroup *models.ChatGroup) error {
	return c.repo.Create(ctx, chatGroup)
}

func (c ChatGroupService) GetById(ctx context.Context, chatGroupId string) (*models.ChatGroup, error) {
	return c.repo.GetByID(ctx, chatGroupId)
}

func (c ChatGroupService) List(ctx context.Context, page, limit int) ([]models.ChatGroup, error) {
	return c.repo.List(ctx, page, limit)
}

func (c ChatGroupService) Update(ctx context.Context, chatGroup *models.ChatGroup) error {
	return c.repo.Update(ctx, chatGroup)
}

func (c ChatGroupService) UpdateGroupName(ctx context.Context, chatGroupId string, newGroupName string) error {
	updateData := map[string]interface{}{
		"groupName": newGroupName,
	}
	return c.repo.UpdateWithFilter(ctx, chatGroupId, updateData)
}

//func (c ChatGroupService) UpdateMembers(ctx context.Context, chatGroupId string, newMembers []primitive.ObjectID) error {
//	updateData := map[string]interface{}{
//		"members": newMembers,
//	}
//	return c.repo.UpdateWithFilter(ctx, chatGroupId, updateData)
//}
//
//func (c ChatGroupService) UpdateAdmins(ctx context.Context, chatGroupId string, newAdmins []primitive.ObjectID) error {
//	updateData := map[string]interface{}{
//		"adminIds": newAdmins,
//	}
//	return c.repo.UpdateWithFilter(ctx, chatGroupId, updateData)
//}

func (c ChatGroupService) Delete(ctx context.Context, chatGroupId string) error {
	return c.repo.Delete(ctx, chatGroupId)
}
