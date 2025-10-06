INSERT INTO genders (gender) VALUES ('Feminino'), ('Masculino');

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
    '0199ba1f-505d-71ed-bacd-ae7cdd67c100',
    '0199ba1f-505d-7378-9424-2fd68424db6a',
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
('11122233344455', '0199ba23-bbb9-7c5c-bbaa-be860b76cce3');

INSERT INTO accounts (id_people, password, avatar)
VALUES
(
    '0199ba23-bbb9-7c5c-bbaa-be860b76cce3',
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
('0199ba2a-294f-73ce-8314-404be25debea', 1, '0199ba23-bbb9-7c5c-bbaa-be860b76cce3'),
-- Admin tem role de admin
('0199bac4-d4c6-78b6-9de0-5fcbe0da0f0d', 2, '0199ba23-bbb9-7c5c-bbaa-be860b76cce3'),

-- Manager tem roles de manager E customer (pode fazer pedidos)
('0199bac4-d4cb-7c71-a049-386abd8da042', 3, '0199ba23-bbb9-7c5c-bbaa-be860b76cce3'),
('0199bac4-d4cb-7c71-a049-386abd8da042', 5, '0199ba23-bbb9-7c5c-bbaa-be860b76cce3'),

-- Editor tem role de editor
('0199bac4-d4d0-735e-8331-04354cd7f9fb', 4, '0199ba23-bbb9-7c5c-bbaa-be860b76cce3'),

-- Customer tem role de customer
('0199bac4-d4d4-7306-92f1-c64b90b5a018', 5, '0199ba23-bbb9-7c5c-bbaa-be860b76cce3');
