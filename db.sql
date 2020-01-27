-- Table: public.account

-- DROP TABLE public.account;

CREATE TABLE public.account
(
    id integer NOT NULL,
    total numeric(6,2) NOT NULL,
    CONSTRAINT account_pkey PRIMARY KEY (id)
)

ALTER TABLE public.account
    OWNER to dev;
)

-- Table: public.account_transaction

-- DROP TABLE public.account_transaction;

CREATE TABLE public.account_transaction
(
    id uuid NOT NULL,
    state character(255) COLLATE pg_catalog."default",
    amount numeric(6,2) NOT NULL,
    source character(255) COLLATE pg_catalog."default",
    CONSTRAINT account_transaction_pkey PRIMARY KEY (id)
)

ALTER TABLE public.account_transaction
    OWNER to dev;