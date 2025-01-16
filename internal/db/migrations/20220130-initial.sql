CREATE TABLE loan_provider (
    provider_id varchar(40) NOT NULL,
    provider varchar(40) NOT NULL,
    provider_type varchar(40) NOT NULL,
    provider_account_id varchar(40) NOT NULL,
    pre_expiration_token_refresh_mins integer NOT NULL,
    provider_name varchar(60) NOT NULL,
    api_key varchar(60) NOT NULL,
    api_secret varchar(100) NOT NULL,
    api_environment varchar(32) NOT NULL,
    api_base_url varchar(100) NOT NULL,
    labels jsonb not null default '{}'::jsonb,
    create_time timestamp without time zone NOT NULL,
    update_time timestamp without time zone NOT NULL,
    provider_config jsonb default '{}'::jsonb
);


CREATE TABLE lsp_token (
    provider_id varchar(40),
    auth_token varchar(256),
    expiration_time timestamp without time zone,
    create_time timestamp without time zone,
    update_time timestamp without time zone
);