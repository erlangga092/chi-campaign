package user

import (
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type Repository interface {
	Save(user User) (User, error)
	FindByID(ID int) (User, error)
	FindByEmail(email string) (User, error)
	Update(userID int, user User) (User, error)
}
type repository struct {
	DB *sql.DB
}

func NewUserRepository(DB *sql.DB) Repository {
	return &repository{DB}
}

// https://kodingin.com/golang-insert-data-mysql/
// reference layoutDateTime for UTC
const (
	layoutDateTime string = "2006-01-02 15:04:05"
)

func (r *repository) Save(user User) (User, error) {
	sqlQuery := sq.Insert("users").Columns("name", "occupation", "email", "password_hash", "avatar_file_name", "role", "created_at", "updated_at").Values(user.Name, user.Occupation, user.Email, user.PasswordHash, user.AvatarFileName, user.Role, time.Now().Format(layoutDateTime), time.Now().Format(layoutDateTime)).RunWith(r.DB)

	result, err := sqlQuery.Exec()
	if err != nil {
		return user, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return user, err
	}

	newUser, err := r.FindByID(int(userID))
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (r *repository) FindByID(ID int) (User, error) {
	user := User{}

	sqlQuery := sq.Select("id", "name", "occupation", "email", "password_hash", "avatar_file_name", "role", "created_at", "updated_at").From("users").Where(sq.Eq{"id": ID})

	rows, err := sqlQuery.RunWith(r.DB).Query()
	if err != nil {
		return user, err
	}

	defer rows.Close()

	if rows.Next() {
		rows.Scan(&user.ID, &user.Name, &user.Occupation, &user.Email, &user.PasswordHash, &user.AvatarFileName, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	}

	return user, nil
}

func (r *repository) FindByEmail(email string) (User, error) {
	user := User{}

	sqlQuery := sq.Select("id", "name", "occupation", "email", "password_hash", "avatar_file_name", "role", "created_at", "updated_at").From("users").Where(sq.Eq{"email": email})

	rows, err := sqlQuery.RunWith(r.DB).Query()
	if err != nil {
		return user, err
	}

	defer rows.Close()

	if rows.Next() {
		rows.Scan(&user.ID, &user.Name, &user.Occupation, &user.Email, &user.PasswordHash, &user.AvatarFileName, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	}

	return user, nil
}

func (r *repository) Update(userID int, user User) (User, error) {
	sqlQuery := sq.Update("users").Set("avatar_file_name", user.AvatarFileName).Where(sq.Eq{"id": userID}).RunWith(r.DB)

	result, err := sqlQuery.Exec()
	if err != nil {
		return user, err
	}

	updatedID, err := result.LastInsertId()
	if err != nil {
		return user, err
	}

	updatedUser, err := r.FindByID(int(updatedID))
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}
