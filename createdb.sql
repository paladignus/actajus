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
-- CREATE DATABASE actajus;
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
	id_people bigint NOT NULL,
	cpf text,
	CONSTRAINT documents_pk PRIMARY KEY (iddocuments),
	CONSTRAINT cpf_uq UNIQUE (cpf)
);
-- ddl-end --
ALTER TABLE public.documents OWNER TO postgres;
-- ddl-end --

-- object: public.users | type: TABLE --
-- DROP TABLE IF EXISTS public.users CASCADE;
CREATE TABLE public.users (
	idusers bigint NOT NULL,
	password text NOT NULL,
	avatar text,
	last_login_at timestamptz,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz,
	CONSTRAINT users_pk PRIMARY KEY (idusers)
);
-- ddl-end --
ALTER TABLE public.users OWNER TO postgres;
-- ddl-end --

-- object: gender_fk | type: CONSTRAINT --
-- ALTER TABLE public.people DROP CONSTRAINT IF EXISTS gender_fk CASCADE;
ALTER TABLE public.people ADD CONSTRAINT gender_fk FOREIGN KEY (id_gender)
REFERENCES public.gender (idgender) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
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

-- object: public.permission_role | type: TABLE --
-- DROP TABLE IF EXISTS public.permission_role CASCADE;
CREATE TABLE public.permission_role (
	id_roles smallint NOT NULL,
	id_permissions smallint NOT NULL,
	CONSTRAINT permission_role_pk PRIMARY KEY (id_roles,id_permissions)
);
-- ddl-end --
ALTER TABLE public.permission_role OWNER TO postgres;
-- ddl-end --

-- object: public.role_user | type: TABLE --
-- DROP TABLE IF EXISTS public.role_user CASCADE;
CREATE TABLE public.role_user (
	id_users bigint NOT NULL,
	id_roles smallint NOT NULL,
	assigned_by bigint NOT NULL,
	created_at timestamptz NOT NULL DEFAULT NOW(),
	CONSTRAINT role_user_pk PRIMARY KEY (id_roles,id_users)
);
-- ddl-end --
ALTER TABLE public.role_user OWNER TO postgres;
-- ddl-end --

-- object: roles_fk | type: CONSTRAINT --
-- ALTER TABLE public.role_user DROP CONSTRAINT IF EXISTS roles_fk CASCADE;
ALTER TABLE public.role_user ADD CONSTRAINT roles_fk FOREIGN KEY (id_roles)
REFERENCES public.roles (idroles) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: roles_fk | type: CONSTRAINT --
-- ALTER TABLE public.permission_role DROP CONSTRAINT IF EXISTS roles_fk CASCADE;
ALTER TABLE public.permission_role ADD CONSTRAINT roles_fk FOREIGN KEY (id_roles)
REFERENCES public.roles (idroles) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: permissions_fk | type: CONSTRAINT --
-- ALTER TABLE public.permission_role DROP CONSTRAINT IF EXISTS permissions_fk CASCADE;
ALTER TABLE public.permission_role ADD CONSTRAINT permissions_fk FOREIGN KEY (id_permissions)
REFERENCES public.permissions (idpermissions) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: people_fk | type: CONSTRAINT --
-- ALTER TABLE public.role_user DROP CONSTRAINT IF EXISTS people_fk CASCADE;
ALTER TABLE public.role_user ADD CONSTRAINT people_fk FOREIGN KEY (assigned_by)
REFERENCES public.people (idpeople) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: public.emails | type: TABLE --
-- DROP TABLE IF EXISTS public.emails CASCADE;
CREATE TABLE public.emails (
	idemails bigint NOT NULL GENERATED ALWAYS AS IDENTITY ,
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

-- object: public.password_reset | type: TABLE --
-- DROP TABLE IF EXISTS public.password_reset CASCADE;
CREATE TABLE public.password_reset (
	idpassword_reset bigint NOT NULL GENERATED ALWAYS AS IDENTITY ,
	id_users bigint NOT NULL,
	token text NOT NULL,
	expires_at timestamptz NOT NULL,
	used_at timestamptz

);
-- ddl-end --
ALTER TABLE public.password_reset OWNER TO postgres;
-- ddl-end --

-- object: people_fk | type: CONSTRAINT --
-- ALTER TABLE public.users DROP CONSTRAINT IF EXISTS people_fk CASCADE;
ALTER TABLE public.users ADD CONSTRAINT people_fk FOREIGN KEY (idusers)
REFERENCES public.people (idpeople) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: users_fk | type: CONSTRAINT --
-- ALTER TABLE public.role_user DROP CONSTRAINT IF EXISTS users_fk CASCADE;
ALTER TABLE public.role_user ADD CONSTRAINT users_fk FOREIGN KEY (id_users)
REFERENCES public.users (idusers) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: users_fk | type: CONSTRAINT --
-- ALTER TABLE public.password_reset DROP CONSTRAINT IF EXISTS users_fk CASCADE;
ALTER TABLE public.password_reset ADD CONSTRAINT users_fk FOREIGN KEY (id_users)
REFERENCES public.users (idusers) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: public.companies | type: TABLE --
-- DROP TABLE IF EXISTS public.companies CASCADE;
CREATE TABLE public.companies (
	idcompanies bigint NOT NULL GENERATED ALWAYS AS IDENTITY ,
	registered_by bigint NOT NULL,
	name text NOT NULL,
	trade_name text NOT NULL,
	cnpj text NOT NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NOT NULL DEFAULT now(),
	deleted_at timestamp,
	CONSTRAINT companies_pk PRIMARY KEY (idcompanies)
);
-- ddl-end --
ALTER TABLE public.companies OWNER TO postgres;
-- ddl-end --

-- object: public.phones | type: TABLE --
-- DROP TABLE IF EXISTS public.phones CASCADE;
CREATE TABLE public.phones (
	idphones bigint NOT NULL GENERATED ALWAYS AS IDENTITY ,
	number text NOT NULL,
	kind text NOT NULL,
	department text,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz,
	CONSTRAINT phones_pk PRIMARY KEY (idphones)
);
-- ddl-end --
ALTER TABLE public.phones OWNER TO postgres;
-- ddl-end --

-- object: people_fk | type: CONSTRAINT --
-- ALTER TABLE public.companies DROP CONSTRAINT IF EXISTS people_fk CASCADE;
ALTER TABLE public.companies ADD CONSTRAINT people_fk FOREIGN KEY (registered_by)
REFERENCES public.people (idpeople) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: public.addresses | type: TABLE --
-- DROP TABLE IF EXISTS public.addresses CASCADE;
CREATE TABLE public.addresses (
	idaddresses bigint NOT NULL GENERATED ALWAYS AS IDENTITY ,
	zip text,
	title text,
	street text,
	number smallint,
	complement text,
	reference text,
	neighborhood text,
	city text,
	state text,
	country text,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz,
	CONSTRAINT addresses_pk PRIMARY KEY (idaddresses)
);
-- ddl-end --
ALTER TABLE public.addresses OWNER TO postgres;
-- ddl-end --

-- object: public.address_person | type: TABLE --
-- DROP TABLE IF EXISTS public.address_person CASCADE;
CREATE TABLE public.address_person (
	id_addresses bigint NOT NULL,
	id_people bigint NOT NULL

);
-- ddl-end --
ALTER TABLE public.address_person OWNER TO postgres;
-- ddl-end --

-- object: addresses_fk | type: CONSTRAINT --
-- ALTER TABLE public.address_person DROP CONSTRAINT IF EXISTS addresses_fk CASCADE;
ALTER TABLE public.address_person ADD CONSTRAINT addresses_fk FOREIGN KEY (id_addresses)
REFERENCES public.addresses (idaddresses) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: people_fk | type: CONSTRAINT --
-- ALTER TABLE public.address_person DROP CONSTRAINT IF EXISTS people_fk CASCADE;
ALTER TABLE public.address_person ADD CONSTRAINT people_fk FOREIGN KEY (id_people)
REFERENCES public.people (idpeople) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: public.person_phone | type: TABLE --
-- DROP TABLE IF EXISTS public.person_phone CASCADE;
CREATE TABLE public.person_phone (
	id_phones bigint NOT NULL,
	id_people bigint NOT NULL,
	CONSTRAINT person_phone_pk PRIMARY KEY (id_people)
);
-- ddl-end --
ALTER TABLE public.person_phone OWNER TO postgres;
-- ddl-end --

-- object: phones_fk | type: CONSTRAINT --
-- ALTER TABLE public.person_phone DROP CONSTRAINT IF EXISTS phones_fk CASCADE;
ALTER TABLE public.person_phone ADD CONSTRAINT phones_fk FOREIGN KEY (id_phones)
REFERENCES public.phones (idphones) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: people_fk | type: CONSTRAINT --
-- ALTER TABLE public.person_phone DROP CONSTRAINT IF EXISTS people_fk CASCADE;
ALTER TABLE public.person_phone ADD CONSTRAINT people_fk FOREIGN KEY (id_people)
REFERENCES public.people (idpeople) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: public.address_enterprise | type: TABLE --
-- DROP TABLE IF EXISTS public.address_enterprise CASCADE;
CREATE TABLE public.address_enterprise (
	id_addresses bigint NOT NULL,
	id_companies bigint NOT NULL

);
-- ddl-end --
ALTER TABLE public.address_enterprise OWNER TO postgres;
-- ddl-end --

-- object: public.enterprise_phone | type: TABLE --
-- DROP TABLE IF EXISTS public.enterprise_phone CASCADE;
CREATE TABLE public.enterprise_phone (
	id_companies bigint NOT NULL,
	id_phones bigint NOT NULL

);
-- ddl-end --
ALTER TABLE public.enterprise_phone OWNER TO postgres;
-- ddl-end --

-- object: addresses_fk | type: CONSTRAINT --
-- ALTER TABLE public.address_enterprise DROP CONSTRAINT IF EXISTS addresses_fk CASCADE;
ALTER TABLE public.address_enterprise ADD CONSTRAINT addresses_fk FOREIGN KEY (id_addresses)
REFERENCES public.addresses (idaddresses) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: companies_fk | type: CONSTRAINT --
-- ALTER TABLE public.address_enterprise DROP CONSTRAINT IF EXISTS companies_fk CASCADE;
ALTER TABLE public.address_enterprise ADD CONSTRAINT companies_fk FOREIGN KEY (id_companies)
REFERENCES public.companies (idcompanies) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: phones_fk | type: CONSTRAINT --
-- ALTER TABLE public.enterprise_phone DROP CONSTRAINT IF EXISTS phones_fk CASCADE;
ALTER TABLE public.enterprise_phone ADD CONSTRAINT phones_fk FOREIGN KEY (id_phones)
REFERENCES public.phones (idphones) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: companies_fk | type: CONSTRAINT --
-- ALTER TABLE public.enterprise_phone DROP CONSTRAINT IF EXISTS companies_fk CASCADE;
ALTER TABLE public.enterprise_phone ADD CONSTRAINT companies_fk FOREIGN KEY (id_companies)
REFERENCES public.companies (idcompanies) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
-- ddl-end --

-- object: public.social_media | type: TABLE --
-- DROP TABLE IF EXISTS public.social_media CASCADE;
CREATE TABLE public.social_media (
	idsocial_media bigint NOT NULL GENERATED ALWAYS AS IDENTITY ,
	id_companies bigint NOT NULL,
	name text NOT NULL,
	url text NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz,
	CONSTRAINT social_media_url_uq UNIQUE (url)
);
-- ddl-end --
ALTER TABLE public.social_media OWNER TO postgres;
-- ddl-end --

-- object: public.email_enterprise | type: TABLE --
-- DROP TABLE IF EXISTS public.email_enterprise CASCADE;
CREATE TABLE public.email_enterprise (
	id_emails bigint NOT NULL,
	id_companies bigint NOT NULL,
	CONSTRAINT email_enterprise_pk PRIMARY KEY (id_emails,id_companies)
);
-- ddl-end --
ALTER TABLE public.email_enterprise OWNER TO postgres;
-- ddl-end --

-- object: emails_fk | type: CONSTRAINT --
-- ALTER TABLE public.email_enterprise DROP CONSTRAINT IF EXISTS emails_fk CASCADE;
ALTER TABLE public.email_enterprise ADD CONSTRAINT emails_fk FOREIGN KEY (id_emails)
REFERENCES public.emails (idemails) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: companies_fk | type: CONSTRAINT --
-- ALTER TABLE public.email_enterprise DROP CONSTRAINT IF EXISTS companies_fk CASCADE;
ALTER TABLE public.email_enterprise ADD CONSTRAINT companies_fk FOREIGN KEY (id_companies)
REFERENCES public.companies (idcompanies) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: public.email_person | type: TABLE --
-- DROP TABLE IF EXISTS public.email_person CASCADE;
CREATE TABLE public.email_person (
	id_emails bigint NOT NULL,
	id_people bigint NOT NULL,
	CONSTRAINT email_person_pk PRIMARY KEY (id_emails,id_people)
);
-- ddl-end --
ALTER TABLE public.email_person OWNER TO postgres;
-- ddl-end --

-- object: emails_fk | type: CONSTRAINT --
-- ALTER TABLE public.email_person DROP CONSTRAINT IF EXISTS emails_fk CASCADE;
ALTER TABLE public.email_person ADD CONSTRAINT emails_fk FOREIGN KEY (id_emails)
REFERENCES public.emails (idemails) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: people_fk | type: CONSTRAINT --
-- ALTER TABLE public.email_person DROP CONSTRAINT IF EXISTS people_fk CASCADE;
ALTER TABLE public.email_person ADD CONSTRAINT people_fk FOREIGN KEY (id_people)
REFERENCES public.people (idpeople) MATCH FULL
ON DELETE CASCADE ON UPDATE CASCADE;
-- ddl-end --

-- object: companies_fk | type: CONSTRAINT --
-- ALTER TABLE public.social_media DROP CONSTRAINT IF EXISTS companies_fk CASCADE;
ALTER TABLE public.social_media ADD CONSTRAINT companies_fk FOREIGN KEY (id_companies)
REFERENCES public.companies (idcompanies) MATCH FULL
ON DELETE RESTRICT ON UPDATE CASCADE;
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