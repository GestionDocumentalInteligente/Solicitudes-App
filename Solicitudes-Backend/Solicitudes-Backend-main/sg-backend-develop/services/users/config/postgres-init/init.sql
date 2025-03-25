-- Table: public.persons

DROP TABLE IF EXISTS public.persons;

CREATE TABLE IF NOT EXISTS public.persons
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    cuil character varying(20) COLLATE pg_catalog."default" NOT NULL,
    dni character varying(20) COLLATE pg_catalog."default" NOT NULL,
    first_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    last_name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    email character varying(255) COLLATE pg_catalog."default",
    phone character varying(50) COLLATE pg_catalog."default",
    CONSTRAINT persons_pkey PRIMARY KEY (id),
    CONSTRAINT persons_cuil_key UNIQUE (cuil),
    CONSTRAINT persons_email_key UNIQUE (email)
)

TABLESPACE pg_default;

DROP INDEX IF EXISTS public.idx_persons_cuil;

CREATE INDEX IF NOT EXISTS idx_persons_cuil
    ON public.persons USING btree
    (cuil COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;

DROP INDEX IF EXISTS public.idx_persons_email;

CREATE INDEX IF NOT EXISTS idx_persons_email
    ON public.persons USING btree
    (email COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;
	
-- Table: public.users

DROP TABLE IF EXISTS public.users;

CREATE TABLE IF NOT EXISTS public.users
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    email_validated boolean DEFAULT false,
    accepts_notifications boolean DEFAULT false,
    logged_at timestamp with time zone,
    person_id bigint,
    deleted_at timestamp with time zone,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT fk_person FOREIGN KEY (person_id)
        REFERENCES public.persons (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE SET NULL
)

TABLESPACE pg_default;