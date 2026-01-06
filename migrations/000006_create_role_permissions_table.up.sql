CREATE TABLE role_permissions
(
    role_id       INT NOT NULL,
    permission_id INT NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    CONSTRAINT fk_role_permissions_roles
        FOREIGN KEY (role_id)
            REFERENCES roles (id),
    CONSTRAINT fk_role_permissions_permissions
        FOREIGN KEY (permission_id)
            REFERENCES users (id)
);