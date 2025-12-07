-- Prepended SQL commands --
CREATE EXTENSION IF NOT EXISTS pgcrypto;-- ddl-end ---- ** Database generated with pgModeler (PostgreSQL Database Modeler).
-- ** pgModeler version: 1.2.2
-- ** PostgreSQL version: 18.0
-- ** Project Site: pgmodeler.io
-- ** Model Author: ---

-- ** Database creation must be performed outside a multi lined SQL file. 
-- ** These commands were put in this file only as a convenience.

-- object: actajus | type: DATABASE --
-- DROP DATABASE IF EXISTS actajus;
CREATE DATABASE actajus;
-- ddl-end --


SET search_path TO pg_catalog,public;
-- ddl-end --

-- object: public.people | type: TABLE --
-- DROP TABLE IF EXISTS public.people CASCADE;
CREATE TABLE public.people (
	idpeople bigint NOT NULL GENERATED ALWAYS AS IDENTITY ,
	first_name text NOT NULL,
	last_name text NOT NULL,
	birthday date,
	id_mother bigint,
	id_father bigint,
	id_gender smallint NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz,
	CONSTRAINT people_pk PRIMARY KEY (idpeople)
);
-- ddl-end --
ALTER TABLE public.people OWNER TO postgres;
-- ddl-end --

-- object: public.gender | type: TABLE --
-- DROP TABLE IF EXISTS public.gender CASCADE;
CREATE TABLE public.gender (
	idgender smallint NOT NULL GENERATED ALWAYS AS IDENTITY ,
	name text NOT NULL,
	CONSTRAINT gender_pkey PRIMARY KEY (idgender)
);
-- ddl-end --
ALTER TABLE public.gender OWNER TO postgres;
-- ddl-end --

-- object: public.documents | type: TABLE --
-- DROP TABLE IF EXISTS public.documents CASCADE;
CREATE TABLE public.documents (
	iddocuments bigint NOT NULL GENERATED ALWAYS AS IDENTITY ,
	cpf text,
	id_people bigint NOT NULL,
	CONSTRAINT documents_pk PRIMARY KEY (iddocuments),
	CONSTRAINT cpf_uq UNIQUE (cpf)
);
-- ddl-end --
ALTER TABLE public.documents OWNER TO postgres;
-- ddl-end --

-- object: public.accounts | type: TABLE --
-- DROP TABLE IF EXISTS public.accounts CASCADE;
CREATE TABLE public.accounts (
	idaccounts bigint NOT NULL GENERATED ALWAYS AS IDENTITY ,
	id_people bigint NOT NULL,
	password text NOT NULL,
	avatar text,
	last_login_at timestamptz,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz,
	CONSTRAINT accounts_pk PRIMARY KEY (idaccounts)
);
-- ddl-end --
ALTER TABLE public.accounts OWNER TO postgres;
-- ddl-end --

-- object: gender_fk | type: CONSTRAINT --
-- ALTER TABLE public.people DROP CONSTRAINT IF EXISTS gender_fk CASCADE;
ALTER TABLE public.people ADD CONSTRAINT gender_fk FOREIGN KEY (id_gender)
REFERENCES public.gender (idgender) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: people_fk | type: CONSTRAINT --
-- ALTER TABLE public.accounts DROP CONSTRAINT IF EXISTS people_fk CASCADE;
ALTER TABLE public.accounts ADD CONSTRAINT people_fk FOREIGN KEY (id_people)
REFERENCES public.people (idpeople) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: accounts_uq | type: CONSTRAINT --
-- ALTER TABLE public.accounts DROP CONSTRAINT IF EXISTS accounts_uq CASCADE;
ALTER TABLE public.accounts ADD CONSTRAINT accounts_uq UNIQUE (id_people);
-- ddl-end --

-- object: people_fk | type: CONSTRAINT --
-- ALTER TABLE public.documents DROP CONSTRAINT IF EXISTS people_fk CASCADE;
ALTER TABLE public.documents ADD CONSTRAINT people_fk FOREIGN KEY (id_people)
REFERENCES public.people (idpeople) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: documents_uq | type: CONSTRAINT --
-- ALTER TABLE public.documents DROP CONSTRAINT IF EXISTS documents_uq CASCADE;
ALTER TABLE public.documents ADD CONSTRAINT documents_uq UNIQUE (id_people);
-- ddl-end --

-- object: public.roles | type: TABLE --
-- DROP TABLE IF EXISTS public.roles CASCADE;
CREATE TABLE public.roles (
	idroles smallint NOT NULL GENERATED ALWAYS AS IDENTITY ,
	name text NOT NULL,
	description text NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	CONSTRAINT roles_pk PRIMARY KEY (idroles)
);
-- ddl-end --
ALTER TABLE public.roles OWNER TO postgres;
-- ddl-end --

-- object: public.permissions | type: TABLE --
-- DROP TABLE IF EXISTS public.permissions CASCADE;
CREATE TABLE public.permissions (
	idpermissions smallint NOT NULL GENERATED ALWAYS AS IDENTITY ,
	resource text NOT NULL,
	action text NOT NULL,
	description text NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	CONSTRAINT permissions_pk PRIMARY KEY (idpermissions),
	CONSTRAINT permissions_uq UNIQUE (resource,action)
);
-- ddl-end --
ALTER TABLE public.permissions OWNER TO postgres;
-- ddl-end --

-- object: public.role_permission | type: TABLE --
-- DROP TABLE IF EXISTS public.role_permission CASCADE;
CREATE TABLE public.role_permission (
	id_roles smallint NOT NULL,
	id_permissions smallint NOT NULL,
	CONSTRAINT role_permission_pk PRIMARY KEY (id_roles,id_permissions)
);
-- ddl-end --
ALTER TABLE public.role_permission OWNER TO postgres;
-- ddl-end --

-- object: public.account_role | type: TABLE --
-- DROP TABLE IF EXISTS public.account_role CASCADE;
CREATE TABLE public.account_role (
	id_accounts bigint NOT NULL,
	id_roles smallint NOT NULL,
	assigned_by bigint NOT NULL,
	created_at timestamptz NOT NULL DEFAULT NOW(),
	CONSTRAINT account_role_pk PRIMARY KEY (id_accounts,id_roles)
);
-- ddl-end --
ALTER TABLE public.account_role OWNER TO postgres;
-- ddl-end --

-- object: accounts_fk | type: CONSTRAINT --
-- ALTER TABLE public.account_role DROP CONSTRAINT IF EXISTS accounts_fk CASCADE;
ALTER TABLE public.account_role ADD CONSTRAINT accounts_fk FOREIGN KEY (id_accounts)
REFERENCES public.accounts (idaccounts) MATCH FULL
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --

-- object: roles_fk | type: CONSTRAINT --
-- ALTER TABLE public.account_role DROP CONSTRAINT IF EXISTS roles_fk CASCADE;
ALTER TABLE public.account_role ADD CONSTRAINT roles_fk FOREIGN KEY (id_roles)
REFERENCES public.roles (idroles) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: roles_fk | type: CONSTRAINT --
-- ALTER TABLE public.role_permission DROP CONSTRAINT IF EXISTS roles_fk CASCADE;
ALTER TABLE public.role_permission ADD CONSTRAINT roles_fk FOREIGN KEY (id_roles)
REFERENCES public.roles (idroles) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: permissions_fk | type: CONSTRAINT --
-- ALTER TABLE public.role_permission DROP CONSTRAINT IF EXISTS permissions_fk CASCADE;
ALTER TABLE public.role_permission ADD CONSTRAINT permissions_fk FOREIGN KEY (id_permissions)
REFERENCES public.permissions (idpermissions) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: people_fk | type: CONSTRAINT --
-- ALTER TABLE public.account_role DROP CONSTRAINT IF EXISTS people_fk CASCADE;
ALTER TABLE public.account_role ADD CONSTRAINT people_fk FOREIGN KEY (assigned_by)
REFERENCES public.people (idpeople) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: public.emails | type: TABLE --
-- DROP TABLE IF EXISTS public.emails CASCADE;
CREATE TABLE public.emails (
	idemails bigint NOT NULL GENERATED ALWAYS AS IDENTITY ,
	id_people bigint NOT NULL,
	address text NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz,
	CONSTRAINT emails_pk PRIMARY KEY (idemails),
	CONSTRAINT emails_name_uq UNIQUE (address)
);
-- ddl-end --
ALTER TABLE public.emails OWNER TO postgres;
-- ddl-end --

-- object: people_fk | type: CONSTRAINT --
-- ALTER TABLE public.emails DROP CONSTRAINT IF EXISTS people_fk CASCADE;
ALTER TABLE public.emails ADD CONSTRAINT people_fk FOREIGN KEY (id_people)
REFERENCES public.people (idpeople) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: public.password_reset | type: TABLE --
-- DROP TABLE IF EXISTS public.password_reset CASCADE;
CREATE TABLE public.password_reset (
	idpassword_reset bigint NOT NULL GENERATED ALWAYS AS IDENTITY ,
	id_people bigint NOT NULL,
	token text NOT NULL,
	expires_at timestamptz NOT NULL,
	used_at timestamptz

);
-- ddl-end --
ALTER TABLE public.password_reset OWNER TO postgres;
-- ddl-end --

-- object: mother_fk | type: CONSTRAINT --
-- ALTER TABLE public.people DROP CONSTRAINT IF EXISTS mother_fk CASCADE;
ALTER TABLE public.people ADD CONSTRAINT mother_fk FOREIGN KEY (id_mother)
REFERENCES public.people (idpeople) MATCH FULL
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --

-- object: father_fk | type: CONSTRAINT --
-- ALTER TABLE public.people DROP CONSTRAINT IF EXISTS father_fk CASCADE;
ALTER TABLE public.people ADD CONSTRAINT father_fk FOREIGN KEY (id_father)
REFERENCES public.people (idpeople) MATCH FULL
ON DELETE SET NULL ON UPDATE CASCADE;
-- ddl-end --


-- Appended SQL commands --
ALTER DATABASE actajus SET datestyle TO "ISO, DMY";
-- ddl-end --