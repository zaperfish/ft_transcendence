package user

import (
	// Internal
	"context"
	"errors"
	"net/http"

    // Internal
	"ft_transcendence/backend/auth"
	"ft_transcendence/backend/errs"
)

type UserService interface {
	CreateUser(ctx context.Context, in CreateUserDTO) (*UserSummaryDTO, error)
	GetUserByID(ctx context.Context, id uint) (*UserSummaryDTO, error)
	GetUserByName(ctx context.Context, name string) (*UserSummaryDTO, error)
	GetUsers(ctx context.Context, page, pageSize int) ([]UserSummaryDTO, error)
	PatchUser(ctx context.Context, id uint, in PatchUserDTO) (*UserSummaryDTO, error)
	PatchPassword(ctx context.Context, id uint, in PatchPasswordDTO) (*UserSummaryDTO, error)
	DeleteUser(ctx context.Context, id uint) error
	LoginUser(ctx context.Context, name, password string) (*UserSummaryDTO, http.Cookie, error)
}

type UserServiceImpl struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, in CreateUserDTO) (*UserSummaryDTO, error) {
	if err := validateParameters(&in); err != nil {
		return nil, errs.ErrInvalidInput
	}

	hash, err := auth.CreateHash(in.Password)
	if err != nil {
		return nil, errs.ErrInternal
	}

    u := User {
        Name:       	in.Name,
        Email:      	in.Email,
        PasswordHash:   hash,
    }

    if err = s.repo.Create(ctx, &u); err != nil {
		return nil, err
	}

    return u.ToSummaryDTO(), nil
}

func validateParameters(u *CreateUserDTO) error {
	if err := ValidUserName(u.Name); err != nil {
		return err
	}
	if err := ValidUserEmail(u.Email); err != nil {
		return err
	}
	if err := ValidUserPassword(u.Password); err != nil {
		return err
	}
	if u.Password != u.PasswordConfirm {
		return errors.New("passwords do not  match")
	}
	return nil
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
		if err := ValidUserName(*in.Name); err != nil {
			return err
		}
		(*updates)["name"] = *in.Name
	}
	if in.Email != nil {
		if err := ValidUserEmail(*in.Email); err != nil {
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

	if err := ValidUserPassword(in.NewPassword); err != nil {
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

func (s *UserServiceImpl) LoginUser(ctx context.Context, name, password string) (*UserSummaryDTO, http.Cookie, error) {
    u, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, http.Cookie{}, err
	}
	match, err := auth.MatchPassword(password, u.PasswordHash)
	if err != nil {
        return nil, http.Cookie{}, errs.ErrInternal
    }
	if !match {
        return nil, http.Cookie{}, errs.ErrNotFound
	}
	cookie, err := auth.MakeJWTCookieFromID(u.ID)
	if err != nil {
        return nil, http.Cookie{}, errs.ErrInternal
    }
	return u.ToSummaryDTO(), cookie, nil
}
