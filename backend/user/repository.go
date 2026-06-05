package user

import (
	// Std
	"context"

	// Internal
	"ft_transcendence/backend/errs"

	// External
	"gorm.io/gorm"
)

type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id uint) (*User, error)
	GetByName(ctx context.Context, name string) (*User, error)
    List(ctx context.Context, filter UserFilter) ([]User, error)
	UpdateFieldsByID(ctx context.Context, id uint, fields map[string]any) (*User, error)
    DeleteByID(ctx context.Context, id uint) error
}

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) GetByID(ctx context.Context, id uint) (*User, error) {
	u, err := gorm.G[User](r.db.Debug()).Where("id = ?", id).First(ctx)
	return &u, errs.ErrorDB(err)
}

func (r *userRepositoryImpl) GetByName(ctx context.Context, name string) (*User, error) {
    u, err := gorm.G[User](r.db.Debug()).Where("name = ?", name).First(ctx)
	return &u, errs.ErrorDB(err)
}

func (r *userRepositoryImpl) List(ctx context.Context, filter UserFilter) ([]User, error) {
	offset := (filter.Page - 1) * filter.PageSize

    us, err := gorm.G[User](r.db.Debug()).Limit(filter.PageSize).Offset(offset).Find(ctx)
    if err != nil {
        return nil, errs.ErrorDB(err)
    }

	return us, nil
}

func (r *userRepositoryImpl) Create(ctx context.Context, input *User) error {

	err := gorm.G[User](r.db.Debug()).Create(ctx, input)
	if err != nil {
		return errs.ErrorDB(err)
	}

	return nil
}

func (r *userRepositoryImpl) UpdateFieldsByID(ctx context.Context, id uint, fields map[string]any) (*User, error) {

	_, err := gorm.G[map[string]any](r.db.Debug()).Table("users").Where("id = ?", id).Updates(ctx, fields)
	if err != nil {
		return nil, errs.ErrorDB(err)
	}

	updated, err := gorm.G[User](r.db.Debug()).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, errs.ErrorDB(err)
	}

	return &updated, nil
}

func (r *userRepositoryImpl) DeleteByID(ctx context.Context, id uint) error {
	rows, err := gorm.G[User](r.db.Debug()).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return errs.ErrorDB(err)
	}
	if rows == 0 {
		return errs.ErrNotFound
	}
	return nil
}
