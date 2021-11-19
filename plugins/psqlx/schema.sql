CREATE TABLE ccdns_domains (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    disabled    BOOL DEFAULT 'f',

    soa_mname   VARCHAR(255) NOT NULL,
    soa_rname   VARCHAR(255) NOT NULL,
    soa_serial  BIGINT NOT NULL DEFAULT extract(epoch from now())::int,
    soa_refresh INT NOT NULL DEFAULT 10800,
    soa_retry   INT NOT NULL DEFAULT 2400,
    soa_expire  INT NOT NULL DEFAULT 604800,
    soa_minimum BIGINT NOT NULL DEFAULT 3600,

    CHECK (soa_serial > 0 AND soa_serial < 4294967295),
    CHECK (soa_refresh > 0),
    CHECK (soa_retry > 0),
    CHECK (soa_expire > 0),
    CHECK (soa_minimum > 0 AND soa_minimum < 4294967295),

    CONSTRAINT c_lowercase_name CHECK (((name)::TEXT = LOWER((name)::TEXT)))
);

CREATE UNIQUE INDEX name_index ON ccdns_domains(name);


CREATE TABLE ccdns_records_ns (
    id          BIGSERIAL PRIMARY KEY,
    domain_id   BIGINT NOT NULL,
    disabled    BOOL DEFAULT 'f',
    name        VARCHAR(255) NOT NULL,
    ttl         INT NOT NULL DEFAULT 300,
    host        VARCHAR(255) NOT NULL,

    CHECK (ttl >= 0),

    CONSTRAINT domain_exists
        FOREIGN KEY(domain_id) REFERENCES ccdns_domains(id)
        ON DELETE CASCADE,

    CONSTRAINT c_lowercase_name CHECK (((name)::TEXT = LOWER((name)::TEXT)))
);

CREATE TABLE ccdns_records_a (
    id          BIGSERIAL PRIMARY KEY,
    domain_id   BIGINT NOT NULL,
    disabled    BOOL DEFAULT 'f',
    name        VARCHAR(255) NOT NULL,
    ttl         INT NOT NULL DEFAULT 300,
    address     INET NOT NULL,

    CHECK (ttl >= 0),

    CONSTRAINT domain_exists
        FOREIGN KEY(domain_id) REFERENCES ccdns_domains(id)
        ON DELETE CASCADE,

    CONSTRAINT c_lowercase_name CHECK (((name)::TEXT = LOWER((name)::TEXT)))
);

CREATE TABLE ccdns_records_aaaa (
    id          BIGSERIAL PRIMARY KEY,
    domain_id   BIGINT NOT NULL,
    disabled    BOOL DEFAULT 'f',
    name        VARCHAR(255) NOT NULL,
    ttl         INT NOT NULL DEFAULT 300,
    address     INET NOT NULL,

    CHECK (ttl >= 0),

    CONSTRAINT domain_exists
        FOREIGN KEY(domain_id) REFERENCES ccdns_domains(id)
        ON DELETE CASCADE,

    CONSTRAINT c_lowercase_name CHECK (((name)::TEXT = LOWER((name)::TEXT)))
);

CREATE TABLE ccdns_records_cname (
    id          BIGSERIAL PRIMARY KEY,
    domain_id   BIGINT NOT NULL,
    disabled    BOOL DEFAULT 'f',
    name        VARCHAR(255) NOT NULL,
    ttl         INT NOT NULL DEFAULT 300,
    host        VARCHAR(255) NOT NULL,

    CHECK (ttl >= 0),

    CONSTRAINT domain_exists
        FOREIGN KEY(domain_id) REFERENCES ccdns_domains(id)
        ON DELETE CASCADE,

    CONSTRAINT c_lowercase_name CHECK (((name)::TEXT = LOWER((name)::TEXT)))
);

CREATE TABLE ccdns_records_aname (
    id          BIGSERIAL PRIMARY KEY,
    domain_id   BIGINT NOT NULL,
    disabled    BOOL DEFAULT 'f',
    name        VARCHAR(255) NOT NULL,
    ttl         INT NOT NULL DEFAULT 300,
    location    VARCHAR(255) NOT NULL,

    CHECK (ttl >= 0),

    CONSTRAINT domain_exists
        FOREIGN KEY(domain_id) REFERENCES ccdns_domains(id)
        ON DELETE CASCADE,

    CONSTRAINT c_lowercase_name CHECK (((name)::TEXT = LOWER((name)::TEXT)))
);

CREATE TABLE ccdns_records_txt (
    id          BIGSERIAL PRIMARY KEY,
    domain_id   BIGINT NOT NULL,
    disabled    BOOL DEFAULT 'f',
    name        VARCHAR(255) NOT NULL,
    ttl         INT NOT NULL DEFAULT 300,
    text        VARCHAR(65535) NOT NULL,

    CHECK (ttl >= 0),

    CONSTRAINT domain_exists
        FOREIGN KEY(domain_id) REFERENCES ccdns_domains(id)
        ON DELETE CASCADE,

    CONSTRAINT c_lowercase_name CHECK (((name)::TEXT = LOWER((name)::TEXT)))
);

CREATE TABLE ccdns_records_mx (
    id          BIGSERIAL PRIMARY KEY,
    domain_id   BIGINT NOT NULL,
    disabled    BOOL DEFAULT 'f',
    name        VARCHAR(255) NOT NULL,
    ttl         INT NOT NULL DEFAULT 300,
    host        VARCHAR(255) NOT NULL,
    preference  INT NOT NULL,

    CHECK (ttl >= 0),
    CHECK (preference >= 0 AND preference < 65535),

    CONSTRAINT domain_exists
        FOREIGN KEY(domain_id) REFERENCES ccdns_domains(id)
        ON DELETE CASCADE,

    CONSTRAINT c_lowercase_name CHECK (((name)::TEXT = LOWER((name)::TEXT)))
);

CREATE TABLE ccdns_records_srv (
    id          BIGSERIAL PRIMARY KEY,
    domain_id   BIGINT NOT NULL,
    disabled    BOOL DEFAULT 'f',
    name        VARCHAR(255) NOT NULL,
    ttl         INT NOT NULL DEFAULT 300,
    target      VARCHAR(255) NOT NULL,
    port        INT NOT NULL,
    priority    INT NOT NULL,
    weight      INT NOT NULL,

    CHECK (ttl >= 0),

    CONSTRAINT domain_exists
        FOREIGN KEY(domain_id) REFERENCES ccdns_domains(id)
        ON DELETE CASCADE,

    CONSTRAINT c_lowercase_name CHECK (((name)::TEXT = LOWER((name)::TEXT)))
);

CREATE TABLE ccdns_records_caa (
    id          BIGSERIAL PRIMARY KEY,
    domain_id   BIGINT NOT NULL,
    disabled    BOOL DEFAULT 'f',
    name        VARCHAR(255) NOT NULL,
    ttl         INT NOT NULL DEFAULT 300,
    tag         VARCHAR(255) NOT NULL,
    value       VARCHAR(255) NOT NULL,
    flag        INT NOT NULL,

    CHECK (ttl >= 0),

    CONSTRAINT domain_exists
        FOREIGN KEY(domain_id) REFERENCES ccdns_domains(id)
        ON DELETE CASCADE,

    CONSTRAINT c_lowercase_name CHECK (((name)::TEXT = LOWER((name)::TEXT)))
);

CREATE TABLE ccdns_records_ptr (
    id          BIGSERIAL PRIMARY KEY,
    domain_id   BIGINT NOT NULL,
    disabled    BOOL DEFAULT 'f',
    name        VARCHAR(255) NOT NULL,
    ttl         INT NOT NULL DEFAULT 300,
    host        VARCHAR(255) NOT NULL,

    CHECK (ttl >= 0),

    CONSTRAINT domain_exists
        FOREIGN KEY(domain_id) REFERENCES ccdns_domains(id)
        ON DELETE CASCADE,

    CONSTRAINT c_lowercase_name CHECK (((name)::TEXT = LOWER((name)::TEXT)))
);

CREATE TABLE ccdns_records_tlsa (
    id          BIGSERIAL PRIMARY KEY,
    domain_id   BIGINT NOT NULL,
    disabled    BOOL DEFAULT 'f',
    name        VARCHAR(255) NOT NULL,
    ttl         INT NOT NULL DEFAULT 300,
    usage       INT NOT NULL,
    selector    INT NOT NULL,
    matching_type   INT NOT NULL,
    certificate VARCHAR(255) NOT NULL,

    CHECK (ttl >= 0),

    CONSTRAINT domain_exists
        FOREIGN KEY(domain_id) REFERENCES ccdns_domains(id)
        ON DELETE CASCADE,

    CONSTRAINT c_lowercase_name CHECK (((name)::TEXT = LOWER((name)::TEXT)))
);
