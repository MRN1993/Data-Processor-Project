package repository

import (
    "database/sql"
    "data-processor-project/internal/domain/models"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (repo *UserRepository) AddUser(user models.User) error {
    _, err := repo.db.Exec("INSERT INTO users (id, quota) VALUES (?, ?)", user.ID, user.Quota)
    return err
}

func (repo *UserRepository) GetUser(userID string) (models.User, error) {
    var user models.User
    row := repo.db.QueryRow("SELECT id, quota FROM users WHERE id = ?", userID)
    err := row.Scan(&user.ID, &user.Quota)
    return user, err
}
