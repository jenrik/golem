CREATE TABLE pages (jobId varchar(36), link text, done bool);
CREATE UNIQUE INDEX ON pages(jobId, link);

CREATE TABLE jobs (jobId varchar(36), defs json);
CREATE INDEX ON jobs (jobId);

CREATE TABLE data (jobId varchar(36), link text, data jsonb, error text);
