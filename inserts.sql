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
INSERT INTO permission_role (id_roles, id_permissions)
SELECT 1 AS id_roles, idpermissions AS id_permissions
FROM permissions RETURNING *;

-- ADMIN: Quase todas, exceto algumas críticas
INSERT INTO permission_role (id_roles, id_permissions)
SELECT 2, idpermissions FROM permissions
WHERE resource != 'settings' OR action = 'read' RETURNING *;

-- MANAGER: Gerenciar pedidos e visualizar relatórios
INSERT INTO permission_role (id_roles, id_permissions)
SELECT 3, idpermissions FROM permissions
WHERE
    resource = 'orders' AND action != 'create' AND action != 'cancel'
    OR resource = 'reports' AND action = 'read' OR action = 'export' RETURNING *;

-- EDITOR: Gerenciar posts
INSERT INTO permission_role (id_roles, id_permissions)
SELECT 4, idpermissions FROM permissions
WHERE
    resource = 'posts' AND action != 'publish' RETURNING *;

-- CUSTOMER: Criar e ver seus próprios pedidos
INSERT INTO permission_role (id_roles, id_permissions)
SELECT 5, idpermissions FROM permissions
WHERE
    resource = 'orders' AND (action = 'create' OR action = 'read')
	OR resource = 'posts' AND action = 'read' RETURNING *;

-- VIEWER: Apenas leitura
INSERT INTO permission_role (id_roles, id_permissions)
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
role_map AS (
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
INSERT INTO role_user (id_users, id_roles, assigned_by)
SELECT am.idusers, ri.idroles, rp.assigned_by
FROM role_map am
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

INSERT INTO role_user (id_users, id_roles, assigned_by) VALUES ('3', 5, '1') RETURNING *;


-- 6. Associar Emails aos Usuários
-- WITH people AS (
--     SELECT
--         (SELECT idpeople FROM people WHERE first_name = 'Sofia' AND last_name = 'Jackson' LIMIT 1) AS id_root,
--         (SELECT idpeople FROM people WHERE first_name = 'Admin' AND last_name = 'System' LIMIT 1) AS id_admin,
--         (SELECT idpeople FROM people WHERE first_name = 'John' AND last_name = 'Manager' LIMIT 1) AS id_manager,
--         (SELECT idpeople FROM people WHERE first_name = 'Jane' AND last_name = 'Editor' LIMIT 1) AS id_editor,
--         (SELECT idpeople FROM people WHERE first_name = 'Bob' AND last_name = 'Customer' LIMIT 1) AS id_customer
-- )
-- INSERT INTO emails (id_people, address)
-- SELECT
--     CASE v.role
--         WHEN 'root' THEN p.id_root
--         WHEN 'admin' THEN p.id_admin
--         WHEN 'manager' THEN p.id_manager
--         WHEN 'editor' THEN p.id_editor
--         WHEN 'customer' THEN p.id_customer
--     END AS id_people,
--     v.address
-- FROM people p
-- CROSS JOIN (
--     VALUES
--         ('root',     'marcelo@marcelo.eti.br'),
--         ('admin',    'admin@email.com.br'),
--         ('manager',  'wostemberg3@gmail.com'),
--         ('editor',   'editor@email.com.br'),
--         ('customer', 'customer@email.com.br')
-- ) AS v(role, address) RETURNING *;

INSERT INTO emails (address) VALUES 
    ('marcelo@marcelo.eti.br'),
    ('admin@email.com.br'),
    ('wostemberg3@gmail.com'),
    ('editor@email.com.br'),
    ('customer@email.com.br') RETURNING *;

WITH person_map AS (
    SELECT idpeople, first_name, last_name
    FROM people
    WHERE (first_name, last_name) IN (
        ('Sofia', 'Jackson'),
        ('Admin', 'System'),
        ('John', 'Manager'),
        ('Jane', 'Editor'),
        ('Bob', 'Customer')
    )
),
email_map AS (
    SELECT idemails, address
    FROM emails
    WHERE address IN (
        'marcelo@marcelo.eti.br',
        'admin@email.com.br',
        'wostemberg3@gmail.com',
        'editor@email.com.br',
        'customer@email.com.br'
    )
),
mapping AS (
    SELECT 
        em.idemails,
        pm.idpeople
    FROM (VALUES
        ('marcelo@marcelo.eti.br', 'Sofia', 'Jackson'),
        ('admin@email.com.br', 'Admin', 'System'),
        ('wostemberg3@gmail.com', 'John', 'Manager'),
        ('editor@email.com.br', 'Jane', 'Editor'),
        ('customer@email.com.br', 'Bob', 'Customer')
    ) AS v(email_addr, fname, lname)
    JOIN person_map pm ON pm.first_name = v.fname AND pm.last_name = v.lname
    JOIN email_map em ON em.address = v.email_addr
)
INSERT INTO email_person (id_emails, id_people)
SELECT idemails, idpeople
FROM mapping
RETURNING *;

-- =============================================================
-- Função para inserir 100.000 empresas fake no banco Actajus
-- Schema baseado em: createdb.sql
-- 
-- Pré-requisitos:
--   1. Existir ao menos 1 registro em public.gender
--   2. Existir ao menos 1 registro em public.people (para registered_by)
--
-- Execução:
--   SELECT seed_fake_companies(100000);
-- =============================================================

CREATE OR REPLACE FUNCTION seed_fake_companies(p_total INT DEFAULT 100000)
RETURNS VOID AS $$
DECLARE
  -- IDs retornados pelos INSERTs
  v_company_id   BIGINT;
  v_address_id   BIGINT;
  v_phone_id     BIGINT;
  v_email_id     BIGINT;
  v_people_id    BIGINT;
  v_doc_id       BIGINT;

  -- Variáveis de apoio
  v_cnpj         TEXT;
  v_cnpj_raw     TEXT;
  v_email        TEXT;
  v_cpf          TEXT;
  v_cpf_raw      TEXT;
  v_digits       INT[];
  v_sum          INT;
  v_rem          INT;
  d1             INT;
  d2             INT;
  v_platform     TEXT;
  v_gender_id    SMALLINT;
  v_registered_by BIGINT;

  -- Arrays de dados fake
  v_first_names  TEXT[] := ARRAY[
    'Ana','Bruno','Carlos','Daniela','Eduardo','Fernanda','Gabriel','Helena',
    'Igor','Juliana','Kevin','Larissa','Marcos','Natália','Otávio','Paula',
    'Rafael','Sabrina','Thiago','Úrsula','Vitor','Wendy','Xavier','Yasmin','Zeca'
  ];
  v_last_names   TEXT[] := ARRAY[
    'Silva','Santos','Oliveira','Souza','Rodrigues','Ferreira','Alves','Pereira',
    'Lima','Gomes','Costa','Ribeiro','Martins','Carvalho','Almeida','Lopes',
    'Sousa','Fernandes','Vieira','Barbosa','Rocha','Dias','Nascimento','Andrade'
  ];
  v_company_names TEXT[] := ARRAY[
    'Alpha','Beta','Gamma','Delta','Epsilon','Zeta','Iota','Kappa',
    'Lambda','Sigma','Nexus','Apex','Vertex','Orbit','Fusion','Pulse',
    'Vortex','Zenith','Praxis','Synapse','Axiom','Cipher','Forge','Helix'
  ];
  v_trade_suffixes TEXT[] := ARRAY[
    'Sistema Jurídico','Consultoria','Soluções','Tecnologia','Serviços',
    'Assessoria','Advocacia','Gestão','Inovação','Digital'
  ];
  v_streets     TEXT[] := ARRAY[
    'Paulista','Atlântica','Brasil','Independência','Liberdade','Consolação',
    'Augusta','Rebouças','Faria Lima','Berrini','Ipiranga','das Flores',
    'do Comércio','das Nações','da República'
  ];
  v_titles      TEXT[] := ARRAY['Avenida','Rua','Alameda','Travessa','Praça'];
  v_neighborhoods TEXT[] := ARRAY[
    'Centro','Jardins','Pinheiros','Vila Madalena','Moema','Itaim Bibi',
    'Lapa','Santana','Tatuapé','Mooca','Bela Vista','Consolação'
  ];
  v_cities      TEXT[] := ARRAY[
    'São Paulo','Rio de Janeiro','Belo Horizonte','Salvador','Curitiba',
    'Porto Alegre','Fortaleza','Recife','Goiânia','Manaus','Belém','Campinas'
  ];
  v_states      TEXT[] := ARRAY[
    'São Paulo','Rio de Janeiro','Minas Gerais','Bahia','Paraná',
    'Rio Grande do Sul','Ceará','Pernambuco','Goiás','Amazonas','Pará','São Paulo'
  ];
  v_kinds       TEXT[] := ARRAY['Celular','Fixo','VoIP'];
  v_departments TEXT[] := ARRAY['Comercial','Suporte','Financeiro','RH','Jurídico'];
  v_platforms   TEXT[] := ARRAY['facebook','Instagram','Linkedin','X','Site'];

  i          INT;
  arr_len    INT;

BEGIN

  -- Busca um gender_id válido (deve existir ao menos 1)
  SELECT idgender INTO v_gender_id FROM public.gender LIMIT 1;
  IF v_gender_id IS NULL THEN
    RAISE EXCEPTION 'Nenhum registro encontrado em public.gender. Insira ao menos um gênero antes de executar esta função.';
  END IF;

  FOR i IN 1..p_total LOOP

    -- =========================================================
    -- 1. Gerar CNPJ válido
    -- =========================================================
    v_cnpj_raw := lpad((floor(random()*99999999)::BIGINT)::TEXT, 8, '0')
               || lpad((floor(random()*9999)::INT + 1)::TEXT, 4, '0'); -- evita 0000

    v_digits := ARRAY(
      SELECT substring(v_cnpj_raw, s, 1)::INT
      FROM generate_series(1, 12) s
    );

    -- 1º dígito verificador
    v_sum := v_digits[1]*5  + v_digits[2]*4  + v_digits[3]*3  + v_digits[4]*2
           + v_digits[5]*9  + v_digits[6]*8  + v_digits[7]*7  + v_digits[8]*6
           + v_digits[9]*5  + v_digits[10]*4 + v_digits[11]*3 + v_digits[12]*2;
    v_rem := v_sum % 11;
    d1 := CASE WHEN v_rem < 2 THEN 0 ELSE 11 - v_rem END;

    -- 2º dígito verificador
    v_sum := v_digits[1]*6  + v_digits[2]*5  + v_digits[3]*4  + v_digits[4]*3
           + v_digits[5]*2  + v_digits[6]*9  + v_digits[7]*8  + v_digits[8]*7
           + v_digits[9]*6  + v_digits[10]*5 + v_digits[11]*4 + v_digits[12]*3
           + d1*2;
    v_rem := v_sum % 11;
    d2 := CASE WHEN v_rem < 2 THEN 0 ELSE 11 - v_rem END;

    v_cnpj := substring(v_cnpj_raw,1,2) || '.'
           || substring(v_cnpj_raw,3,3) || '.'
           || substring(v_cnpj_raw,6,3) || '/'
           || substring(v_cnpj_raw,9,4) || '-'
           || lpad(d1::TEXT,2,'0');
    -- Corrigindo: os dois dígitos finais são d1 e d2
    v_cnpj := substring(v_cnpj_raw,1,2) || '.'
           || substring(v_cnpj_raw,3,3) || '.'
           || substring(v_cnpj_raw,6,3) || '/'
           || substring(v_cnpj_raw,9,4) || '-'
           || lpad(d1::TEXT,1,'0') || lpad(d2::TEXT,1,'0');

    -- =========================================================
    -- 2. Gerar CPF válido (para a pessoa registrante)
    -- =========================================================
    v_cpf_raw := lpad((floor(random()*999999999)::BIGINT)::TEXT, 9, '0');

    v_digits := ARRAY(
      SELECT substring(v_cpf_raw, s, 1)::INT
      FROM generate_series(1, 9) s
    );

    -- 1º dígito CPF
    v_sum := 0;
    FOR arr_len IN 1..9 LOOP
      v_sum := v_sum + v_digits[arr_len] * (11 - arr_len);
    END LOOP;
    v_rem := v_sum % 11;
    d1 := CASE WHEN v_rem < 2 THEN 0 ELSE 11 - v_rem END;

    -- 2º dígito CPF
    v_sum := 0;
    FOR arr_len IN 1..9 LOOP
      v_sum := v_sum + v_digits[arr_len] * (12 - arr_len);
    END LOOP;
    v_sum := v_sum + d1 * 2;
    v_rem := v_sum % 11;
    d2 := CASE WHEN v_rem < 2 THEN 0 ELSE 11 - v_rem END;

    v_cpf := substring(v_cpf_raw,1,3) || '.'
          || substring(v_cpf_raw,4,3) || '.'
          || substring(v_cpf_raw,7,3) || '-'
          || lpad(d1::TEXT,1,'0') || lpad(d2::TEXT,1,'0');

    -- =========================================================
    -- 3. Inserir People (pessoa que registra a empresa)
    -- =========================================================
    INSERT INTO public.people (
      first_name, last_name, birthday, id_gender,
      created_at, updated_at
    ) VALUES (
      v_first_names[1 + (floor(random() * array_length(v_first_names,1)))::INT % array_length(v_first_names,1)],
      v_last_names [1 + (floor(random() * array_length(v_last_names ,1)))::INT % array_length(v_last_names ,1)],
      ('1960-01-01'::date + (floor(random() * 22000))::INT),  -- datas entre 1960 e ~2020
      v_gender_id,
      NOW(), NOW()
    ) RETURNING idpeople INTO v_people_id;

    -- =========================================================
    -- 4. Inserir Documents (CPF da pessoa)
    -- =========================================================
    INSERT INTO public.documents (id_people, cpf)
    VALUES (v_people_id, v_cpf);

    -- =========================================================
    -- 5. Inserir Address
    -- =========================================================
    INSERT INTO public.addresses (
      zip, title, street, number, complement,
      reference, neighborhood, city, state, country,
      created_at, updated_at
    ) VALUES (
      lpad((floor(random()*99999))::INT::TEXT, 5, '0') || '-'
        || lpad((floor(random()*999))::INT::TEXT, 3, '0'),
      v_titles[1 + (floor(random()*array_length(v_titles,1)))::INT % array_length(v_titles,1)],
      v_streets[1 + (floor(random()*array_length(v_streets,1)))::INT % array_length(v_streets,1)],
      (floor(random()*9999 + 1))::SMALLINT,
      'Sala ' || (floor(random()*999 + 1))::INT::TEXT,
      'Próximo ao ponto ' || i::TEXT,
      v_neighborhoods[1 + (floor(random()*array_length(v_neighborhoods,1)))::INT % array_length(v_neighborhoods,1)],
      v_cities[1 + (floor(random()*array_length(v_cities,1)))::INT % array_length(v_cities,1)],
      v_states[1 + (floor(random()*array_length(v_states,1)))::INT % array_length(v_states,1)],
      'Brasil',
      NOW(), NOW()
    ) RETURNING idaddresses INTO v_address_id;

    -- =========================================================
    -- 6. Inserir Phone
    -- =========================================================
    INSERT INTO public.phones (
      number, kind, department, created_at, updated_at
    ) VALUES (
      '+5511' || lpad((floor(random()*999999999))::BIGINT::TEXT, 9, '0'),
      v_kinds[1 + (floor(random()*array_length(v_kinds,1)))::INT % array_length(v_kinds,1)],
      v_departments[1 + (floor(random()*array_length(v_departments,1)))::INT % array_length(v_departments,1)],
      NOW(), NOW()
    ) RETURNING idphones INTO v_phone_id;

    -- =========================================================
    -- 7. Inserir Email (único, usando i como semente)
    -- =========================================================
    v_email := 'empresa' || i::TEXT
            || '.' || lpad((floor(random()*9999))::INT::TEXT, 4, '0')
            || '@actajusfake.com.br';

    INSERT INTO public.emails (address, created_at, updated_at)
    VALUES (v_email, NOW(), NOW())
    RETURNING idemails INTO v_email_id;

    -- =========================================================
    -- 8. Inserir Company
    -- =========================================================
    INSERT INTO public.companies (
      registered_by, name, trade_name, cnpj,
      created_at, updated_at
    ) VALUES (
      v_people_id,
      v_company_names[1 + (floor(random()*array_length(v_company_names,1)))::INT % array_length(v_company_names,1)]
        || ' ' || lpad(i::TEXT, 6, '0'),
      v_company_names[1 + (floor(random()*array_length(v_company_names,1)))::INT % array_length(v_company_names,1)]
        || ' ' || v_trade_suffixes[1 + (floor(random()*array_length(v_trade_suffixes,1)))::INT % array_length(v_trade_suffixes,1)],
      v_cnpj,
      NOW() + (i || ' seconds')::INTERVAL,
      NOW() + (i || ' seconds')::INTERVAL
    ) RETURNING idcompanies INTO v_company_id;
    -- INSERT INTO public.companies (
    --   registered_by, name, trade_name, cnpj,
    --   created_at, updated_at
    -- ) VALUES (
    --   v_people_id,
    --   v_company_names[1 + (floor(random()*array_length(v_company_names,1)))::INT % array_length(v_company_names,1)]
    --     || ' ' || lpad(i::TEXT, 6, '0'),
    --   v_company_names[1 + (floor(random()*array_length(v_company_names,1)))::INT % array_length(v_company_names,1)]
    --     || ' ' || v_trade_suffixes[1 + (floor(random()*array_length(v_trade_suffixes,1)))::INT % array_length(v_trade_suffixes,1)],
    --   v_cnpj,
    --   NOW(), NOW()
    -- ) RETURNING idcompanies INTO v_company_id;

    -- =========================================================
    -- 9. Relacionamentos (tabelas pivot)
    -- =========================================================

    -- Endereço da empresa
    INSERT INTO public.company_address (id_addresses, id_companies)
    VALUES (v_address_id, v_company_id);

    -- Telefone da empresa
    INSERT INTO public.company_phone (id_companies, id_phones)
    VALUES (v_company_id, v_phone_id);

    -- Email da empresa
    INSERT INTO public.company_email (id_emails, id_companies)
    VALUES (v_email_id, v_company_id);

    -- =========================================================
    -- 10. Social Media (todas as 5 plataformas por empresa)
    -- =========================================================
    FOREACH v_platform IN ARRAY v_platforms LOOP
      INSERT INTO public.social_media (
        id_companies, platform, url, created_at, updated_at
      ) VALUES (
        v_company_id,
        v_platform,
        'https://www.' || lower(v_platform) || '.com/actajus_fake_'
          || lower(v_platform) || '_' || i::TEXT,
        NOW(), NOW()
      );
    END LOOP;

    -- =========================================================
    -- Log de progresso a cada 5.000 registros
    -- =========================================================
    IF i % 5000 = 0 THEN
      RAISE NOTICE '[seed_fake_companies] % / % empresas inseridas...', i, p_total;
    END IF;

  END LOOP;

  RAISE NOTICE '[seed_fake_companies] Concluído! % empresas inseridas com sucesso.', p_total;

END;
$$ LANGUAGE plpgsql;


-- =============================================================
-- EXECUTAR:
-- =============================================================
-- SELECT seed_fake_companies(100000);


-- =============================================================
-- DICA DE PERFORMANCE: desative índices e triggers antes da carga
-- =============================================================
ALTER TABLE public.companies       DISABLE TRIGGER ALL;
ALTER TABLE public.people          DISABLE TRIGGER ALL;
ALTER TABLE public.emails          DISABLE TRIGGER ALL;
ALTER TABLE public.phones          DISABLE TRIGGER ALL;
ALTER TABLE public.addresses       DISABLE TRIGGER ALL;
ALTER TABLE public.social_media    DISABLE TRIGGER ALL;
ALTER TABLE public.company_address DISABLE TRIGGER ALL;
ALTER TABLE public.company_phone   DISABLE TRIGGER ALL;
ALTER TABLE public.company_email   DISABLE TRIGGER ALL;
ALTER TABLE public.documents          DISABLE TRIGGER ALL;
--
SELECT seed_fake_companies(1000);
--
ALTER TABLE public.companies       ENABLE TRIGGER ALL;
ALTER TABLE public.people          ENABLE TRIGGER ALL;
ALTER TABLE public.emails          ENABLE TRIGGER ALL;
ALTER TABLE public.phones          ENABLE TRIGGER ALL;
ALTER TABLE public.addresses       ENABLE TRIGGER ALL;
ALTER TABLE public.social_media    ENABLE TRIGGER ALL;
ALTER TABLE public.company_address ENABLE TRIGGER ALL;
ALTER TABLE public.company_phone   ENABLE TRIGGER ALL;
ALTER TABLE public.company_email   ENABLE TRIGGER ALL;
ALTER TABLE public.documents          ENABLE TRIGGER ALL;