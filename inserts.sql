INSERT INTO gender (name) VALUES ('Feminino'), ('Masculino');

INSERT INTO people (first_name, last_name, birthday, id_gender)
VALUES
('Janet', 'Jackson', '1980-01-01', 1), ('John', 'Jackson', '1980-01-01', 2);

INSERT INTO
people (first_name, last_name, birthday, id_mother, id_father, id_gender)
VALUES
(
    'Sofia',
    'Jackson',
    '1980-01-01',
    '0199bc57-f525-7eed-9ea2-e0098c19607f',
    '0199bc57-f526-7b75-99db-ae3293049df3',
    1
);

-- 4. Criar usuários de exemplo (senha: "password123")
INSERT INTO people (first_name, last_name, birthday, id_gender) VALUES 
('Admin', 'System', '1976-02-26', 2),
('John', 'Manager', '1965-04-13', 2),
('Jane', 'Editor', '1990-09-17', 1),
('Bob', 'Customer', '1993-12-31', 2);

INSERT INTO documents (cpf, id_people)
VALUES
('11122233344', '0199bc58-f61a-76c8-8837-6a273aa9b2d0');

INSERT INTO accounts (id_people, password, avatar)
VALUES
(
    '0199bc58-f61a-76c8-8837-6a273aa9b2d0',
    crypt('admin123', gen_salt('bf')),
    'default.png'
);

INSERT INTO accounts (id_people, password, avatar)
VALUES
(
    '0199bc59-3d57-76f4-a493-7f7929b60c8a',
    crypt('admin123', gen_salt('bf')),
    'default.png'
),
(
    '0199bc59-3d57-77fb-a7ba-1bf13c1b4723',
    crypt('admin123', gen_salt('bf')),
    'default.png'
),
(
    '0199bc59-3d57-781e-8af0-e9a168f2cca7',
    crypt('admin123', gen_salt('bf')),
    'default.png'
),
(
    '0199bc59-3d57-783b-8ee2-848c75d457cd',
    crypt('admin123', gen_salt('bf')),
    'default.png'
);

INSERT INTO roles (name, description)
VALUES
('super_admin', 'Administrador com acesso total'),
('admin', 'Administrador do sistema'),
('manager', 'Gerente com permissões elevadas'),
('editor', 'Editor de conteúdo'),
('customer', 'Cliente do sistema'),
('viewer', 'Visualizador apenas');

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
('reports', 'export', 'Exportar relatórios');

-- SUPER ADMIN: Todas as permissões
INSERT INTO role_permission (id_roles, id_permissions)
SELECT 1 AS id_roles, idpermissions AS id_permissions
FROM permissions;

-- ADMIN: Quase todas, exceto algumas críticas
INSERT INTO role_permission (id_roles, id_permissions)
SELECT 2, idpermissions FROM permissions
WHERE resource != 'settings' OR action = 'read';

-- MANAGER: Gerenciar pedidos e visualizar relatórios
INSERT INTO role_permission (id_roles, id_permissions)
SELECT 3, idpermissions FROM permissions
WHERE
    resource = 'orders' AND action != 'create' AND action != 'cancel'
    OR resource = 'reports' AND action = 'read' OR action = 'export'

-- EDITOR: Gerenciar posts
INSERT INTO role_permission (id_roles, id_permissions)
SELECT 4, idpermissions FROM permissions
WHERE
    resource = 'posts' AND action != 'publish'

-- CUSTOMER: Criar e ver seus próprios pedidos
INSERT INTO role_permission (id_roles, id_permissions)
SELECT 5, idpermissions FROM permissions
WHERE
    resource = 'orders' AND (action = 'create' OR action = 'read')
	OR resource = 'posts' AND action = 'read'

-- VIEWER: Apenas leitura
INSERT INTO role_permission (id_roles, id_permissions)
SELECT 6, idpermissions FROM permissions
WHERE
    resource = 'users' AND action = 'read'
	OR resource = 'posts' AND action = 'read'
	OR resource = 'orders' AND action = 'read'




-- 4. Criar usuários de exemplo (senha: "password123")
INSERT INTO users (id, email, password, first_name, last_name) VALUES 
('99999999-9999-9999-9999-999999999991', 'admin@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye5JtJXY6XvH9q7CNoXvFUr1T/lXHHJuO', 'Admin', 'System'),
('99999999-9999-9999-9999-999999999992', 'manager@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye5JtJXY6XvH9q7CNoXvFUr1T/lXHHJuO', 'John', 'Manager'),
('99999999-9999-9999-9999-999999999993', 'editor@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye5JtJXY6XvH9q7CNoXvFUr1T/lXHHJuO', 'Jane', 'Editor'),
('99999999-9999-9999-9999-999999999994', 'customer@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye5JtJXY6XvH9q7CNoXvFUr1T/lXHHJuO', 'Bob', 'Customer');

-- 5. Associar Roles aos Usuários
INSERT INTO account_role (id_accounts, id_roles, assigned_by) VALUES 
--Super_Admin tem todas as roles
('0199bc5a-62a1-7ca3-994c-3c6493dd1375', 1, '0199bc58-f61a-76c8-8837-6a273aa9b2d0'),
-- Admin tem role de admin
('0199bc60-1e3d-7b0d-a559-00ff768266a6', 2, '0199bc58-f61a-76c8-8837-6a273aa9b2d0'),

-- Manager tem roles de manager E customer (pode fazer pedidos)
('0199bc60-1e40-79b8-a07d-8ca6004c12ba', 3, '0199bc58-f61a-76c8-8837-6a273aa9b2d0'),
('0199bc60-1e40-79b8-a07d-8ca6004c12ba', 5, '0199bc58-f61a-76c8-8837-6a273aa9b2d0'),

-- Editor tem role de editor
('0199bc60-1e43-75d4-aa01-e507a23116fc', 4, '0199bc58-f61a-76c8-8837-6a273aa9b2d0'),

-- Customer tem role de customer
('0199bc60-1e46-70da-a034-33bc4296f326', 5, '0199bc58-f61a-76c8-8837-6a273aa9b2d0');
