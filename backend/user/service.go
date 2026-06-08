package user

import (
	// Internal
	"context"
	"errors"

    // Internal
	"ft_transcendence/backend/auth"
	"ft_transcendence/backend/errs"
)

type UserService interface {
	GetUserByID(ctx context.Context, id uint) (*UserSummaryDTO, error)
	GetUserByName(ctx context.Context, name string) (*UserSummaryDTO, error)
	GetUsers(ctx context.Context, page, pageSize int) ([]UserSummaryDTO, error)
	PatchUser(ctx context.Context, id uint, in PatchUserDTO) (*UserSummaryDTO, error)
	PatchPassword(ctx context.Context, id uint, in PatchPasswordDTO) (*UserSummaryDTO, error)
	DeleteUser(ctx context.Context, id uint) error
}

type UserServiceImpl struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) GetUserByID(ctx context.Context, id uint) (*UserSummaryDTO, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return u.ToSummaryDTO(), nil
}

func (s *UserServiceImpl) GetUserByName(ctx context.Context, name string) (*UserSummaryDTO, error) {
	u, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return u.ToSummaryDTO(), nil
}

func (s *UserServiceImpl) GetUsers(ctx context.Context, page, pageSize int) ([]UserSummaryDTO, error) {
	filter := UserFilter{Page: page, PageSize: pageSize}
	us, err := s.repo.List(ctx, filter)
	if err != nil {
        return nil, err
	}

    userList := make([]UserSummaryDTO, 0, len(us))
    for _, u := range us {
        userList = append(userList, *u.ToSummaryDTO())
    }

	return userList, nil
}

func populateUpdates(updates *map[string]any, in *PatchUserDTO) error {
	if in.Name != nil {
		if err := auth.ValidUserName(*in.Name); err != nil {
			return err
		}
		(*updates)["name"] = *in.Name
	}
	if in.Email != nil {
		if err := auth.ValidUserEmail(*in.Email); err != nil {
			return err
		}
		(*updates)["email"] = *in.Email
	}
	return nil
}


func (s *UserServiceImpl) PatchUser(ctx context.Context, id uint, in PatchUserDTO) (*UserSummaryDTO, error) {
	updates := map[string]any{}
 	if err := populateUpdates(&updates, &in); err != nil {
		return nil, errs.ErrInvalidInput
	}
	u, err := s.repo.UpdateFieldsByID(ctx, id, updates)
	if err != nil {
		return nil, err
	}
	return u.ToSummaryDTO(), nil
}

func (s *UserServiceImpl) PatchPassword(ctx context.Context, id uint, in PatchPasswordDTO) (*UserSummaryDTO, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	match, err := auth.MatchPassword(in.CurrentPassword, u.PasswordHash)
	if err != nil {
		return nil, errs.ErrInternal
	}
	if !match {
		return nil, errs.ErrNotFound
	}

	if in.NewPassword != in.ConfirmPassword {
		return nil, errors.New("new passwords do not match")
	}

	if err := auth.ValidUserPassword(in.NewPassword); err != nil {
		return nil, errs.ErrInvalidInput
	}

	hash, err := auth.CreateHash(in.NewPassword)
	if err != nil {
		return nil, errs.ErrInternal
	}

	u, err = s.repo.UpdateFieldsByID(ctx, id, map[string]any{"password_hash": hash})
	if err != nil {
		return nil, err
	}

	return u.ToSummaryDTO(), nil
}

func (s *UserServiceImpl) DeleteUser(ctx context.Context, id uint) error {
	return s.repo.DeleteByID(ctx, id)
}
