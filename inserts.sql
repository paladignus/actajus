INSERT INTO gender (name) VALUES ('Feminino'), ('Masculino');

INSERT INTO people (first_name, last_name, birthday, id_gender)
VALUES
('Janet', 'Jackson', '1980-01-01', 1), ('John', 'Jackson', '1980-01-01', 2) RETURNING *;

WITH parents AS (
    SELECT
        (SELECT idpeople FROM people WHERE first_name = 'Janet' AND last_name = 'Jackson' LIMIT 1) AS id_mother,
        (SELECT idpeople FROM people WHERE first_name = 'John' AND last_name = 'Jackson' LIMIT 1) AS id_father
)
INSERT INTO people (first_name, last_name, birthday, id_mother, id_father, id_gender)
SELECT
    'Sofia',
    'Jackson',
    '1980-01-01',
    p.id_mother,
    p.id_father,
    1
FROM parents p RETURNING *;

-- 4. Criar usuários de exemplo (senha: "password123")
INSERT INTO people (first_name, last_name, birthday, id_gender) VALUES 
('Admin', 'System', '1976-02-26', 2),
('John', 'Manager', '1965-04-13', 2),
('Jane', 'Editor', '1990-09-17', 1),
('Bob', 'Customer', '1993-12-31', 2) RETURNING *;


WITH people AS (
    SELECT
        (SELECT idpeople FROM people WHERE first_name = 'Sofia' AND last_name = 'Jackson' LIMIT 1) AS id_root,
        (SELECT idpeople FROM people WHERE first_name = 'Admin' AND last_name = 'System' LIMIT 1) AS id_admin,
        (SELECT idpeople FROM people WHERE first_name = 'John' AND last_name = 'Manager' LIMIT 1) AS id_manager,
        (SELECT idpeople FROM people WHERE first_name = 'Jane' AND last_name = 'Editor' LIMIT 1) AS id_editor,
        (SELECT idpeople FROM people WHERE first_name = 'Bob' AND last_name = 'Customer' LIMIT 1) AS id_customer
)
INSERT INTO documents (cpf, id_people)
SELECT v.cpf, 
       CASE v.role
            WHEN 'root' THEN p.id_root
            WHEN 'admin' THEN p.id_admin
            WHEN 'manager' THEN p.id_manager
            WHEN 'editor' THEN p.id_editor
            WHEN 'customer' THEN p.id_customer
       END
FROM people p
CROSS JOIN (
    VALUES
        ('11144477735', 'root'),
        ('12345678909', 'admin'),
        ('98765432100', 'manager'),
        ('01234567890', 'editor'),
        ('71460238001', 'customer')
) AS v(cpf, role) RETURNING *;

WITH people AS (
    SELECT
        (SELECT idpeople FROM people WHERE first_name = 'Sofia' AND last_name = 'Jackson' LIMIT 1) AS id_root,
        (SELECT idpeople FROM people WHERE first_name = 'Admin' AND last_name = 'System' LIMIT 1) AS id_admin,
        (SELECT idpeople FROM people WHERE first_name = 'John' AND last_name = 'Manager' LIMIT 1) AS id_manager,
        (SELECT idpeople FROM people WHERE first_name = 'Jane' AND last_name = 'Editor' LIMIT 1) AS id_editor,
        (SELECT idpeople FROM people WHERE first_name = 'Bob' AND last_name = 'Customer' LIMIT 1) AS id_customer
)
INSERT INTO users (idusers, password, avatar)
SELECT
    CASE v.role
        WHEN 'root' THEN p.id_root
        WHEN 'admin' THEN p.id_admin
        WHEN 'manager' THEN p.id_manager
        WHEN 'editor' THEN p.id_editor
        WHEN 'customer' THEN p.id_customer
    END AS id_people,
    v.password,
    v.avatar
FROM people p
CROSS JOIN (
    VALUES
        ('root',     crypt('admin123', gen_salt('bf')), 'default.png'),
        ('admin',    crypt('admin123', gen_salt('bf')), 'default.png'),
        ('manager',  crypt('admin123', gen_salt('bf')), 'default.png'),
        ('editor',   crypt('admin123', gen_salt('bf')), 'default.png'),
        ('customer', crypt('admin123', gen_salt('bf')), 'default.png')
) AS v(role, password, avatar) RETURNING *;


INSERT INTO roles (name, description)
VALUES
('root', 'Administrador com acesso total'),
('admin', 'Administrador do sistema'),
('manager', 'Gerente com permissões elevadas'),
('editor', 'Editor de conteúdo'),
('customer', 'Cliente do sistema'),
('viewer', 'Visualizador apenas') RETURNING *;

INSERT INTO permissions (resource, action, description)
VALUES
-- Permissões de Usuários
('users', 'create', 'Criar novos usuários'),
('users', 'read', 'Visualizar usuários'),
('users', 'update', 'Atualizar usuários'),
('users', 'delete', 'Deletar usuários'),

-- Permissões de Posts
('posts', 'create', 'Criar posts'),
('posts', 'read', 'Visualizar posts'),
('posts', 'update', 'Atualizar posts'),
('posts', 'delete', 'Deletar posts'),
('posts', 'publish', 'Publicar posts'),

-- Permissões de Pedidos
('orders', 'create', 'Criar pedidos'),
('orders', 'read', 'Visualizar pedidos'),
('orders', 'update', 'Atualizar pedidos'),
('orders', 'cancel', 'Cancelar pedidos'),
('orders', 'approve', 'Aprovar pedidos'),

-- Permissões de Configurações
('settings', 'read', 'Visualizar configurações'),
('settings', 'update', 'Atualizar configurações'),

-- Permissões de Relatórios
('reports', 'read', 'Visualizar relatórios'),
('reports', 'export', 'Exportar relatórios') RETURNING *;

-- ROOT: Todas as permissões
INSERT INTO role_permission (id_roles, id_permissions)
SELECT 1 AS id_roles, idpermissions AS id_permissions
FROM permissions RETURNING *;

-- ADMIN: Quase todas, exceto algumas críticas
INSERT INTO role_permission (id_roles, id_permissions)
SELECT 2, idpermissions FROM permissions
WHERE resource != 'settings' OR action = 'read' RETURNING *;

-- MANAGER: Gerenciar pedidos e visualizar relatórios
INSERT INTO role_permission (id_roles, id_permissions)
SELECT 3, idpermissions FROM permissions
WHERE
    resource = 'orders' AND action != 'create' AND action != 'cancel'
    OR resource = 'reports' AND action = 'read' OR action = 'export' RETURNING *;

-- EDITOR: Gerenciar posts
INSERT INTO role_permission (id_roles, id_permissions)
SELECT 4, idpermissions FROM permissions
WHERE
    resource = 'posts' AND action != 'publish' RETURNING *;

-- CUSTOMER: Criar e ver seus próprios pedidos
INSERT INTO role_permission (id_roles, id_permissions)
SELECT 5, idpermissions FROM permissions
WHERE
    resource = 'orders' AND (action = 'create' OR action = 'read')
	OR resource = 'posts' AND action = 'read' RETURNING *;

-- VIEWER: Apenas leitura
INSERT INTO role_permission (id_roles, id_permissions)
SELECT 6, idpermissions FROM permissions
WHERE
    resource = 'users' AND action = 'read'
	OR resource = 'posts' AND action = 'read'
	OR resource = 'orders' AND action = 'read' RETURNING *;

WITH mapping AS (
  SELECT * FROM (VALUES
    ('root','Sofia','Jackson'),
    ('admin','Admin','System'),
    ('manager','John','Manager'),
    ('editor','Jane','Editor'),
    ('customer','Bob','Customer'),
    ('viewer', NULL, NULL)
  ) AS t(role, first_name, last_name)
),
people_map AS (
  SELECT m.role, p.idpeople AS id_person
  FROM mapping m
  LEFT JOIN people p
    ON p.first_name = m.first_name
   AND p.last_name  = m.last_name
),
role_ids AS (
  SELECT name AS role, idroles
  FROM roles
),
user_map AS (
  SELECT pm.role, a.idusers
  FROM people_map pm
  LEFT JOIN users a
    ON a.idusers = pm.id_person
),
root_person AS (
  SELECT idpeople AS assigned_by
  FROM people
  WHERE first_name = 'Sofia' AND last_name = 'Jackson'
  LIMIT 1
)
INSERT INTO user_role (id_users, id_roles, assigned_by)
SELECT am.idusers, ri.idroles, rp.assigned_by
FROM user_map am
JOIN role_ids ri ON ri.role = am.role
CROSS JOIN root_person rp
WHERE am.idusers IS NOT NULL
  AND ri.idroles IS NOT NULL
ON CONFLICT (id_users, id_roles) DO NOTHING
RETURNING *;

-- 5. Associar Roles aos Usuários
-- INSERT INTO user_role (id_users, id_roles, assigned_by) VALUES 
-- --Super_Admin tem todas as roles
-- ('0199bc5a-62a1-7ca3-994c-3c6493dd1375', 1, '0199bc58-f61a-76c8-8837-6a273aa9b2d0'),
-- -- Admin tem role de admin
-- ('0199bc60-1e3d-7b0d-a559-00ff768266a6', 2, '0199bc58-f61a-76c8-8837-6a273aa9b2d0'),

-- -- Manager tem roles de manager E customer (pode fazer pedidos)
-- ('0199bc60-1e40-79b8-a07d-8ca6004c12ba', 3, '0199bc58-f61a-76c8-8837-6a273aa9b2d0'),
-- ('0199bc60-1e40-79b8-a07d-8ca6004c12ba', 5, '0199bc58-f61a-76c8-8837-6a273aa9b2d0'),

-- -- Editor tem role de editor
-- ('0199bc60-1e43-75d4-aa01-e507a23116fc', 4, '0199bc58-f61a-76c8-8837-6a273aa9b2d0'),

-- -- Customer tem role de customer
-- ('0199bc60-1e46-70da-a034-33bc4296f326', 5, '0199bc58-f61a-76c8-8837-6a273aa9b2d0');

INSERT INTO user_role (id_users, id_roles, assigned_by) VALUES ('3', 5, '1') RETURNING *;

WITH people AS (
    SELECT
        (SELECT idpeople FROM people WHERE first_name = 'Sofia' AND last_name = 'Jackson' LIMIT 1) AS id_root,
        (SELECT idpeople FROM people WHERE first_name = 'Admin' AND last_name = 'System' LIMIT 1) AS id_admin,
        (SELECT idpeople FROM people WHERE first_name = 'John' AND last_name = 'Manager' LIMIT 1) AS id_manager,
        (SELECT idpeople FROM people WHERE first_name = 'Jane' AND last_name = 'Editor' LIMIT 1) AS id_editor,
        (SELECT idpeople FROM people WHERE first_name = 'Bob' AND last_name = 'Customer' LIMIT 1) AS id_customer
)
INSERT INTO emails (id_people, address)
SELECT
    CASE v.role
        WHEN 'root' THEN p.id_root
        WHEN 'admin' THEN p.id_admin
        WHEN 'manager' THEN p.id_manager
        WHEN 'editor' THEN p.id_editor
        WHEN 'customer' THEN p.id_customer
    END AS id_people,
    v.address
FROM people p
CROSS JOIN (
    VALUES
        ('root',     'root@email.com.br'),
        ('admin',    'admin@email.com.br'),
        ('manager',  'manager@email.com.br'),
        ('editor',   'editor@email.com.br'),
        ('customer', 'customer@email.com.br')
) AS v(role, address) RETURNING *;

