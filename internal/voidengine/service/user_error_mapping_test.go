package service

import (
	"BlackHole/internal/voidengine/contract"
	"BlackHole/internal/voidengine/errorcode"
	"BlackHole/internal/voidengine/model"
	"BlackHole/pkg/apperror"
	"context"
	"errors"
	"testing"
)

type failingUserDAO struct {
	createErr error
	updateErr error
}

func (f failingUserDAO) List(context.Context, model.UserQuery) ([]model.User, error) {
	return nil, nil
}
func (f failingUserDAO) Create(context.Context, *model.User) error { return f.createErr }
func (f failingUserDAO) FindByName(context.Context, string) (*model.User, error) {
	return &model.User{Name: "alice"}, nil
}
func (f failingUserDAO) Update(context.Context, string, *model.User) error { return f.updateErr }
func (f failingUserDAO) DeleteByName(context.Context, string) error        { return nil }

func TestAddMapsDuplicateNameError(t *testing.T) {
	dao := failingUserDAO{
		createErr: errors.New("Error 1062: Duplicate entry 'alice' for key 'name'"),
	}
	svc := NewUserService(dao)

	err := svc.Add(context.Background(), contract.AddUserRequest{
		Username: "alice",
		Password: "secret123",
	})
	appErr, ok := apperror.As(err)
	if !ok {
		t.Fatalf("expected apperror, got %v", err)
	}
	if appErr.Code() != errorcode.UserAlreadyExists {
		t.Fatalf("code = %d, want %d", appErr.Code(), errorcode.UserAlreadyExists)
	}
}

func TestAddMapsDuplicateEmailError(t *testing.T) {
	dao := failingUserDAO{
		createErr: errors.New("Error 1062: Duplicate entry 'a@example.com' for key 'email'"),
	}
	svc := NewUserService(dao)

	err := svc.Add(context.Background(), contract.AddUserRequest{
		Username: "alice",
		Password: "secret123",
	})
	appErr, ok := apperror.As(err)
	if !ok {
		t.Fatalf("expected apperror, got %v", err)
	}
	if appErr.Code() != errorcode.EmailAlreadyExists {
		t.Fatalf("code = %d, want %d", appErr.Code(), errorcode.EmailAlreadyExists)
	}
}

func TestModifyMapsDuplicateEmailError(t *testing.T) {
	email := "dup@example.com"
	dao := failingUserDAO{
		updateErr: errors.New("Error 1062: Duplicate entry 'dup@example.com' for key 'idx_email'"),
	}
	svc := NewUserService(dao)

	err := svc.Modify(context.Background(), contract.ModifyUserRequest{
		Username: "alice",
		Email:    &email,
	})
	appErr, ok := apperror.As(err)
	if !ok {
		t.Fatalf("expected apperror, got %v", err)
	}
	if appErr.Code() != errorcode.EmailAlreadyExists {
		t.Fatalf("code = %d, want %d", appErr.Code(), errorcode.EmailAlreadyExists)
	}
}

func TestAddWrapsUnknownDBErrorAsSystemError(t *testing.T) {
	dao := failingUserDAO{
		createErr: errors.New("connection refused"),
	}
	svc := NewUserService(dao)

	err := svc.Add(context.Background(), contract.AddUserRequest{
		Username: "alice",
		Password: "secret123",
	})
	appErr, ok := apperror.As(err)
	if !ok {
		t.Fatalf("expected apperror, got %v", err)
	}
	if appErr.Code() != errorcode.SystemError {
		t.Fatalf("code = %d, want %d", appErr.Code(), errorcode.SystemError)
	}
	if !errors.Is(err, dao.createErr) {
		t.Fatal("cause error is not preserved")
	}
}
