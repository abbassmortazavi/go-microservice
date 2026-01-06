CREATE TABLE user_roles
(
    role_id INT NOT NULL,
    user_id INT NOT NULL,
    PRIMARY KEY (role_id, user_id),
    CONSTRAINT fk_user_roles_role
        FOREIGN KEY (role_id)
            REFERENCES roles(role_id),
    CONSTRAINT fk_user_roles_user
        FOREIGN KEY (user_id)
            REFERENCES users(user_id)
);