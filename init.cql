CREATE KEYSPACE IF NOT EXISTS gravatar WITH replication = {
    'class': 'NetworkTopologyStrategy',
    'replication_factor': 1
    };

USE gravatar;

-- md5

CREATE TABLE IF NOT EXISTS md5_qq_mapping
(
    email_md5 text PRIMARY KEY,
    qq_id     bigint,
);

CREATE INDEX IF NOT EXISTS md5_qq_mapping_qq_id_idx ON md5_qq_mapping (qq_id);

-- sha256

CREATE TABLE IF NOT EXISTS sha256_qq_mapping
(
    email_sha256 text PRIMARY KEY,
    qq_id        bigint,
);

CREATE INDEX IF NOT EXISTS sha256_qq_mapping_qq_id_idx ON sha256_qq_mapping (qq_id);
