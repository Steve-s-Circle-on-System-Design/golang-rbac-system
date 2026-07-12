package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, u *User) error
	// FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	// FindByGoogleID(ctx context.Context, googleID string) (*User, error)
	// Update(ctx context.Context, u *User) error
	// IncrementFailedLogin(ctx context.Context, id string) (int, error)
	// ResetFailedLogin(ctx context.Context, id string) error
	// SetLock(ctx context.Context, id string, until time.Time) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, u *User) error {
	_, err := r.db.Exec(
		ctx,
		`INSERT INTO users (email, password_hash)
     VALUES ($1, $2)`,
		u.Email,
		u.PasswordHash,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	row := r.db.QueryRow(
        ctx,
        `SELECT email, password_hash
         FROM users
         WHERE email = $1`,
        email,
    )

    var u User

    err := row.Scan(
        &u.Email,
        &u.PasswordHash,
    )
    if err != nil {
        return nil, err
    }

    return &u, nil
}
