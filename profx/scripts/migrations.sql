/* DATABASE */

CREATE DATABASE IF NOT EXISTS profx CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
DROP DATABASE profx;

/* MIGRATION TABLES */

CREATE TABLE IF NOT EXISTS profx.migrations
(
    id        SMALLINT UNSIGNED PRIMARY KEY,
    migration SMALLINT UNSIGNED NOT NULL
);
INSERT IGNORE INTO profx.migrations (id, migration)
VALUES (1, 0);

/* MIGRATION 1 */

-- migrate

CREATE DATABASE IF NOT EXISTS profx CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE TABLE IF NOT EXISTS profx.sources
(
    url     VARCHAR(1000)                             NOT NULL,
    rule    VARCHAR(1000)                             NOT NULL,
    created TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6) NOT NULL
);
CREATE TABLE IF NOT EXISTS profx.links
(
    hash    VARCHAR(64) PRIMARY KEY                   NOT NULL,
    url     VARCHAR(1000)                             NOT NULL,
    created TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6) NOT NULL
);
CREATE TABLE IF NOT EXISTS profx.resources
(
    hash        VARCHAR(64) PRIMARY KEY                   NOT NULL,
    kind        VARCHAR(64)                               NOT NULL,
    url         VARCHAR(1000)                             NOT NULL,
    title       VARCHAR(1000),
    description VARCHAR(1000),
    created     TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6) NOT NULL
);
CREATE INDEX idx_resources_kind ON profx.resources (kind);
CREATE TABLE IF NOT EXISTS profx.sys_logs
(
    message VARCHAR(1000)                             NOT NULL,
    created TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6) NOT NULL
);

-- revert

DROP TABLE IF EXISTS profx.sources;
DROP TABLE IF EXISTS profx.links;
DROP TABLE IF EXISTS profx.resources;
DROP TABLE IF EXISTS profx.sys_logs;

/* MIGRATION 2 */

-- migrate

ALTER TABLE profx.sources
    ADD COLUMN name VARCHAR(64) NOT NULL FIRST,
    ADD COLUMN kind VARCHAR(64) NOT NULL AFTER url;
ALTER TABLE profx.sources
    DROP COLUMN rule;
CREATE TABLE IF NOT EXISTS profx.source_rules
(
    source  VARCHAR(64)                               NOT NULL,
    type    TINYINT UNSIGNED                          NOT NULL,
    rule    VARCHAR(1000)                             NOT NULL,
    created TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6) NOT NULL
);

-- revert

ALTER TABLE profx.sources
    DROP COLUMN name,
    DROP COLUMN kind;
ALTER TABLE profx.sources
    ADD COLUMN rule VARCHAR(1000) NOT NULL;
DROP TABLE IF EXISTS profx.source_rules;

/* MIGRATION 3 */

-- migrate

RENAME TABLE profx.source_rules TO profx.rules;
ALTER TABLE profx.links
    ADD COLUMN source   VARCHAR(64) AFTER url,
    ADD COLUMN from_url VARCHAR(1000) AFTER source;
ALTER TABLE profx.resources
    ADD COLUMN source   VARCHAR(64) AFTER kind,
    ADD COLUMN from_url VARCHAR(1000) AFTER source;

CREATE INDEX idx_links_source ON profx.links (source);
CREATE INDEX idx_resources_source ON profx.resources (source);

-- revert

ALTER TABLE profx.links
    DROP COLUMN source,
    DROP COLUMN from_url;
ALTER TABLE profx.resources
    DROP COLUMN source,
    DROP COLUMN from_url;
RENAME TABLE profx.rules TO profx.source_rules;

# No need to drop indexes, they are dropped when column is removed.
# DROP INDEX idx_links_source ON profx.links;
# DROP INDEX idx_resources_source ON profx.resources;
