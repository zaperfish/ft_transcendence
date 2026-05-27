package user

import (
	// Std
	"context"
	"errors"

	// External
	"gorm.io/gorm"
)

// type UserRepository interface {
//     create(ctx context.Context, user *User) error
//     getByID(ctx context.Context, id uint) (*User, error)
//     list(ctx context.Context, filter UserFilter) ([]User, error)
// 	updateFieldsByID(ctx context.Context, id uint, fields map[string]any) (*User, error)
//     delete(ctx context.Context, id uint) error
// }

func (h handler) getByID(ctx context.Context, id uint) (*User, error) {
	u, err := gorm.G[User](h.db).Where("id = ?", id).First(ctx)
	return &u, err
}

func (h handler) list(ctx context.Context, filter UserFilter) ([]User, error) {
	offset := (filter.Page - 1) * filter.PageSize

    us, err := gorm.G[User](h.db).Limit(filter.PageSize).Offset(offset).Find(ctx)
    if err != nil {
        return nil, err
    }

	return us, nil
}

func (h handler) create(ctx context.Context, input *User) error {

	err := gorm.G[User](h.db).Create(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (h handler) updateFieldsByID(ctx context.Context, id uint, fields map[string]any) (*User, error) {

	_, err := gorm.G[map[string]any](h.db.Debug()).Table("users").Where("id = ?", id).Updates(ctx, fields)
	if err != nil {
		return nil, err
	}

	updated, err := gorm.G[User](h.db.Debug()).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (h handler) deleteByID(ctx context.Context, id uint) error {
	rows, err := gorm.G[User](h.db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("no user deleted")
	}
	return nil
}
