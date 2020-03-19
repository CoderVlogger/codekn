CREATE DATABASE IF NOT EXISTS profx CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
SELECT * FROM profx.migrations;

SELECT * FROM profx.sources;
SELECT COUNT(*) FROM profx.sources;
SELECT * FROM profx.rules;
SELECT * FROM profx.rules sr WHERE sr.source = 'discord';
SELECT COUNT(*) FROM profx.rules;
TRUNCATE TABLE profx.sources;

SELECT * FROM profx.links ORDER BY created LIMIT 50;
SELECT * FROM profx.links ORDER BY created DESC LIMIT 50;
SELECT COUNT(*) FROM profx.links;
TRUNCATE TABLE profx.links;

SELECT * FROM profx.resources ORDER BY created DESC;
SELECT created, source, url FROM profx.resources ORDER BY created DESC;
SELECT * FROM profx.resources ORDER BY created DESC LIMIT 50;
SELECT COUNT(*) FROM profx.resources;
TRUNCATE TABLE profx.resources;

SELECT * FROM profx.sys_logs ORDER BY created DESC LIMIT 50;
TRUNCATE TABLE profx.sys_logs;
