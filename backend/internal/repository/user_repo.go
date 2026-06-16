package repository

import (
	"context"
	"time"

	"oa-nsdiy/backend/ent"
	"oa-nsdiy/backend/ent/user"
	"oa-nsdiy/backend/internal/domain"
)

type UserRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{client: client}
}

// Type alias for backward compatibility
type User = domain.User

func toUser(e *ent.User) *User {
	if e == nil {
		return nil
	}
	return &User{
		ID:                        e.ID,
		Username:                  e.Username,
		Email:                     e.Email,
		Nickname:                  strPtr(e.Nickname),
		Salt:                      e.Salt,
		HashedPassword:            e.HashedPassword,
		RoleID:                    intPtr(e.RoleID),
		UserType:                  e.UserType,
		IsActive:                  e.IsActive,
		TokenVersion:              e.TokenVersion,
		VerificationCode:          strPtr(e.VerificationCode),
		VerificationCodeExpiresAt: timePtr(e.VerificationCodeExpiresAt),
		CreatedAt:                 e.CreatedAt,
		UpdatedAt:                 e.UpdatedAt,
	}
}

func toUsers(es []*ent.User) []*User {
	out := make([]*User, 0, len(es))
	for _, e := range es {
		out = append(out, toUser(e))
	}
	return out
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*User, error) {
	e, err := r.client.User.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toUser(e), nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	e, err := r.client.User.Query().Where(user.UsernameEQ(username)).First(ctx)
	if err != nil {
		return nil, err
	}
	return toUser(e), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	e, err := r.client.User.Query().Where(user.EmailEQ(email)).First(ctx)
	if err != nil {
		return nil, err
	}
	return toUser(e), nil
}

func (r *UserRepository) GetActiveByID(ctx context.Context, id int) (*User, error) {
	e, err := r.client.User.Query().
		Where(user.IDEQ(id), user.IsActive(true)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	return toUser(e), nil
}

func (r *UserRepository) Create(ctx context.Context, u *User) error {
	b := r.client.User.Create().
		SetUsername(u.Username).
		SetEmail(u.Email).
		SetSalt(u.Salt).
		SetHashedPassword(u.HashedPassword).
		SetUserType(u.UserType).
		SetIsActive(u.IsActive).
		SetTokenVersion(u.TokenVersion)

	if u.Nickname != nil {
		b.SetNickname(*u.Nickname)
	}
	if u.RoleID != nil {
		b.SetRoleID(*u.RoleID)
	}

	e, err := b.Save(ctx)
	if err != nil {
		return err
	}
	u.ID = e.ID
	return nil
}

func (r *UserRepository) Update(ctx context.Context, u *User) error {
	b := r.client.User.UpdateOneID(u.ID).
		SetUsername(u.Username).
		SetEmail(u.Email).
		SetUserType(u.UserType).
		SetIsActive(u.IsActive).
		SetTokenVersion(u.TokenVersion)

	if u.Nickname != nil {
		b.SetNickname(*u.Nickname)
	}
	if u.RoleID != nil {
		b.SetRoleID(*u.RoleID)
	}

	_, err := b.Save(ctx)
	return err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id int, salt, hashedPassword string) error {
	_, err := r.client.User.UpdateOneID(id).
		SetSalt(salt).
		SetHashedPassword(hashedPassword).
		AddTokenVersion(1).
		Save(ctx)
	return err
}

func (r *UserRepository) IncrementTokenVersion(ctx context.Context, id int) error {
	_, err := r.client.User.UpdateOneID(id).AddTokenVersion(1).Save(ctx)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	return r.client.User.DeleteOneID(id).Exec(ctx)
}

func (r *UserRepository) SetVerificationCode(ctx context.Context, id int, code string, expiresAt time.Time) error {
	_, err := r.client.User.UpdateOneID(id).
		SetVerificationCode(code).
		SetVerificationCodeExpiresAt(expiresAt).
		Save(ctx)
	return err
}

func (r *UserRepository) ClearVerificationCode(ctx context.Context, id int) error {
	_, err := r.client.User.UpdateOneID(id).
		ClearVerificationCode().
		ClearVerificationCodeExpiresAt().
		Save(ctx)
	return err
}

func (r *UserRepository) List(ctx context.Context, page, pageSize int) ([]*User, int64, error) {
	q := r.client.User.Query()
	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	es, err := q.Order(ent.Desc(user.FieldID)).
		Limit(pageSize).
		Offset(offset).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return toUsers(es), int64(total), nil
}

func (r *UserRepository) Count(ctx context.Context) (int, error) {
	return r.client.User.Query().Count(ctx)
}
