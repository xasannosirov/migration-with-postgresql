package storage

import (
	"database/sql"
	"migration/models"

	_ "github.com/lib/pq"
)

func connect() (*sql.DB, error) {
	dsn := "user=postgres password=1234 dbname=backend sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return &sql.DB{}, err
	}
	return db, nil
}

func CreateUser(user *models.User) (*models.User, error) {
	db, err := connect()
	if err != nil {
		return &models.User{}, err
	}

	defer db.Close()

	query := `
	INSERT INTO users (
		first_name, 
		last_name, 
		gender, 
		email,
		password
	) 
	VALUES($1, $2, $3, $4, $5) 
  	RETURNING 
		id, 
		first_name,
		last_name,
		gender,
		email,
		password`

	var respUser models.User
	if err = db.QueryRow(
		query,
		user.FirstName,
		user.LastName,
		user.Gender,
		user.Email,
		user.Password,
	).Scan(
		&respUser.Id,
		&respUser.FirstName,
		&respUser.LastName,
		&respUser.Gender,
		&respUser.Email,
		&respUser.Password,
	); err != nil {
		return &models.User{}, err
	}

	queryRole := `
	INSERT INTO roles (
		role_name, 
		user_id
	)
	VALUES ($1, $2) 
	RETURNING role_name`

	if err := db.QueryRow(
		queryRole,
		user.Role,
		respUser.Id,
	).Scan(&respUser.Role); err != nil {
		return &models.User{}, err
	}

	return &respUser, nil
}

func UpdateUser(userId string, user *models.User) (*models.User, error) {
	db, err := connect()
	if err != nil {
		return &models.User{}, err
	}
	defer db.Close()

	query := `
  	UPDATE 
    	users 
  	SET 
    	first_name = $1, 
    	last_name = $2
  	WHERE 
    	id = $3
  	RETURNING 
    	id, 
    	first_name, 
    	last_name,
		gender,
		email,
		password`

	var respUser models.User
	if err := db.QueryRow(
		query,
		user.FirstName,
		user.LastName,
		userId,
	).Scan(
		&respUser.Id,
		&respUser.FirstName,
		&respUser.LastName,
		&respUser.Gender,
		&respUser.Email,
		&respUser.Password,
	); err != nil {
		return &models.User{}, err
	}

	queryRole := `
	SELECT 
		role_name 
	FROM 
		roles 
	WHERE 
		user_id = $1`

	if err := db.QueryRow(
		queryRole, 
		userId,
	).Scan(&respUser.Role); err != nil {
		return &models.User{}, err
	}

	return &respUser, nil
}

func DeleteUser(userId string) error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()

	queryRole := `DELETE FROM roles WHERE user_id = $1`
	_, err = db.Exec(queryRole, userId)
	if err != nil {
		return err
	}

	query := `DELETE FROM users WHERE id = $1`
	_, err = db.Exec(query, userId)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(userId string) (*models.User, error) {
	db, err := connect()
	if err != nil {
		return &models.User{}, err
	}
	defer db.Close()

	query := `
	SELECT 
		id, 
		first_name, 
		last_name, 
		gender, 
		email,
		password
	FROM 
		users 
	WHERE 
		id = $1`

	var respUser models.User
	if err = db.QueryRow(query, userId).Scan(
		&respUser.Id,
		&respUser.FirstName,
		&respUser.LastName,
		&respUser.Gender,
		&respUser.Email,
		&respUser.Password,
	); err != nil {
		return &models.User{}, err
	}

	queryRole := `
	SELECT 
		role_name 
	FROM 
		roles 
	WHERE 
		user_id = $1`

	if err := db.QueryRow(
		queryRole, 
		userId,
	).Scan(&respUser.Role); err != nil {
		return &models.User{}, err
	}

	return &respUser, nil
}

func GetAllUsers(page, limit int) ([]*models.User, error) {
	db, err := connect()
	if err != nil {
		return []*models.User{}, err
	}
	defer db.Close()

	var users []*models.User
	offset := limit * (page - 1)
	query := `
	SELECT 
		id, 
		first_name, 
		last_name,
		gender,
		email,
		password
	FROM 
		users 
	LIMIT $1 
	OFFSET $2`

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return []*models.User{}, err
	}

	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.Gender,
			&user.Email,
			&user.Password,
		); err != nil {
			return []*models.User{}, err
		}
		queryRole := `
		SELECT 
			role_name 
		FROM 
			roles 
		WHERE 
			user_id = $1`

		if err := db.QueryRow(
			queryRole, 
			user.Id,
		).Scan(&user.Role); err != nil {
			return []*models.User{}, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func GetUsersByRole(role string, page, limit int) ([]*models.User, error) {
	db, err := connect()
	if err != nil {
		return []*models.User{}, err
	}
	defer db.Close()

	queryRole := `
	SELECT 
		role_name, 
		user_id 
	FROM 
		roles 
	WHERE 
		lower(role_name) = $1 
	LIMIT $2 
	OFFSET $3`

	offset := limit * (page - 1)
	rows, err := db.Query(
		queryRole, 
		role, 
		limit, 
		offset,
	)
	if err != nil {
		return []*models.User{}, err
	}

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.Role, 
			&user.Id,
		); err != nil {
			return []*models.User{}, err
		}
		query := `
		SELECT 
			first_name, 
			last_name,
			gender,
			email,
			password
		FROM 
			users
		WHERE 
			id = $1`

		if err := db.QueryRow(
			query, 
			user.Id,
		).Scan(
			&user.FirstName,
			&user.LastName,
			&user.Gender,
			&user.Email,
			&user.Password,
		); err != nil {
			return []*models.User{}, err
		}
		users = append(users, &user)
	}

	return users, nil
}
