CREATE TABLE roles (
    id SERIAL PRIMARY KEY NOT NULL,
    role_name VARCHAR(64) NOT NULL,
    user_id UUID,
    FOREIGN KEY(user_id) REFERENCES users(id)
);