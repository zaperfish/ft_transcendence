package user

import (
	// Std
	"context"

	// Internal
	"ft_transcendence/backend/errs"

	// External
	"gorm.io/gorm"
)

// type UserRepository interface {
//     create(ctx context.Context, user *User) error
//     getUserByID(ctx context.Context, id uint) (*User, error)
//     list(ctx context.Context, filter UserFilter) ([]User, error)
// 	updateUserFieldsByID(ctx context.Context, id uint, fields map[string]any) (*User, error)
//     delete(ctx context.Context, id uint) error
// }

func (h Handler) getUserByID(ctx context.Context, id uint) (*User, error) {
	u, err := gorm.G[User](h.DB).Where("id = ?", id).First(ctx)
	return &u, errs.ErrorDB(err)
}

func (h Handler) getUserByName(ctx context.Context, name string) (*User, error) {
    u, err := gorm.G[User](h.DB).Where("name = ?", name).First(ctx)
	return &u, errs.ErrorDB(err)
}

func (h Handler) getUsersList(ctx context.Context, filter UserFilter) ([]User, error) {
	offset := (filter.Page - 1) * filter.PageSize

    us, err := gorm.G[User](h.DB).Limit(filter.PageSize).Offset(offset).Find(ctx)
    if err != nil {
        return nil, errs.ErrorDB(err)
    }

	return us, nil
}

func (h Handler) creatUser(ctx context.Context, input *User) error {

	err := gorm.G[User](h.DB).Create(ctx, input)
	if err != nil {
		return errs.ErrorDB(err)
	}

	return nil
}

func (h Handler) updateUserFieldsByID(ctx context.Context, id uint, fields map[string]any) (*User, error) {

	_, err := gorm.G[map[string]any](h.DB.Debug()).Table("users").Where("id = ?", id).Updates(ctx, fields)
	if err != nil {
		return nil, errs.ErrorDB(err)
	}

	updated, err := gorm.G[User](h.DB.Debug()).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, errs.ErrorDB(err)
	}

	return &updated, nil
}

func (h Handler) deleteUserByID(ctx context.Context, id uint) error {
	rows, err := gorm.G[User](h.DB).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return errs.ErrorDB(err)
	}
	if rows == 0 {
		return errs.ErrNotFound
	}
	return nil
}
