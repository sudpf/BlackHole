package service

import (
	"BlackHole/internal/voidengine/contract"
	"BlackHole/internal/voidengine/errorcode"
	"BlackHole/internal/voidengine/model"
	"BlackHole/pkg/apperror"
	"context"
	"testing"
)

type userDataAccessStub struct {
	listUsers      []model.User
	user           *model.User
	deleteCalled   bool
	deleteUserName string
}

func (s userDataAccessStub) List(context.Context, model.UserQuery) ([]model.User, error) {
	return s.listUsers, nil
}

func (userDataAccessStub) Create(context.Context, *model.User) error {
	return nil
}

func (s *userDataAccessStub) FindByName(context.Context, string) (*model.User, error) {
	return s.user, nil
}

func (*userDataAccessStub) Update(context.Context, string, *model.User) error {
	return nil
}

func (s *userDataAccessStub) DeleteByName(_ context.Context, username string) error {
	s.deleteCalled = true
	s.deleteUserName = username
	return nil
}

func TestListReturnsContractUsers(t *testing.T) {
	service := NewUserService(&userDataAccessStub{
		listUsers: []model.User{{
			Name:     "alice",
			Password: "secret",
			Email:    "alice@example.com",
			Phone:    "123",
		}},
	})

	users, err := service.List(context.Background(), contract.ListUserRequest{})
	if err != nil {
		t.Fatalf("List error = %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("len(users) = %d, want 1", len(users))
	}
	if users[0].Username != "alice" {
		t.Fatalf("username = %q, want alice", users[0].Username)
	}
	if users[0].Email != "alice@example.com" {
		t.Fatalf("email = %q, want alice@example.com", users[0].Email)
	}
	if users[0].Phone != "123" {
		t.Fatalf("phone = %q, want 123", users[0].Phone)
	}
}

func TestModifyReturnsUserNotFoundCode(t *testing.T) {
	service := NewUserService(&userDataAccessStub{})
	err := service.Modify(context.Background(), contract.ModifyUserRequest{Username: "missing"})

	appErr, ok := apperror.As(err)
	if !ok {
		t.Fatalf("Modify error type = %T, want application error", err)
	}
	if appErr.Code() != errorcode.UserNotFound {
		t.Fatalf("code = %d, want %d", appErr.Code(), errorcode.UserNotFound)
	}
	if appErr.Params()["username"] != "missing" {
		t.Fatalf("username param = %v, want missing", appErr.Params()["username"])
	}
}

func TestDeleteReturnsUserNotFoundWithUsername(t *testing.T) {
	stub := &userDataAccessStub{}
	service := NewUserService(stub)
	err := service.Delete(context.Background(), contract.DeleteUserRequest{Username: "missing"})

	appErr, ok := apperror.As(err)
	if !ok {
		t.Fatalf("Delete error type = %T, want application error", err)
	}
	if appErr.Code() != errorcode.UserNotFound {
		t.Fatalf("code = %d, want %d", appErr.Code(), errorcode.UserNotFound)
	}
	if appErr.Params()["username"] != "missing" {
		t.Fatalf("username param = %v, want missing", appErr.Params()["username"])
	}
	if stub.deleteCalled {
		t.Fatal("DeleteByName was called for a missing user")
	}
}

func TestDeleteChecksUserBeforeDeleting(t *testing.T) {
	stub := &userDataAccessStub{user: &model.User{Name: "alice"}}
	service := NewUserService(stub)

	if err := service.Delete(context.Background(), contract.DeleteUserRequest{Username: "alice"}); err != nil {
		t.Fatalf("Delete error = %v", err)
	}
	if !stub.deleteCalled {
		t.Fatal("DeleteByName was not called")
	}
	if stub.deleteUserName != "alice" {
		t.Fatalf("deleted username = %q, want alice", stub.deleteUserName)
	}
}
