CREATE TABLE user_roles
(
    role_id INT NOT NULL,
    user_id INT NOT NULL,
    PRIMARY KEY (role_id, user_id),
    CONSTRAINT fk_user_roles_roles
        FOREIGN KEY (role_id)
            REFERENCES roles(id),
    CONSTRAINT fk_user_roles_users
        FOREIGN KEY (user_id)
            REFERENCES users(id)
);