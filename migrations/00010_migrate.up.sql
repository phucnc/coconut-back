BEGIN;

CREATE TABLE IF NOT EXISTS exchange_event_buy_token
(
    block_number    bigint          NOT NULL,
    block_hash      text            NOT NULL,
    block_timestamp timestamptz     NOT NULL,
    tx_index        bigint          NOT NULL,
    tx_hash         text            NOT NULL,
    nft_token_id    numeric(120, 0) NOT NULL,
    nft_price       numeric(120, 0) NOT NULL,
    nft_quote_token smallint        NOT NULL,
    log_index       bigint          NOT NULL,
    log_address     text            NOT NULL,
    log_data        bytea           NOT NULL,
    log_removed     bool            NOT NULL,
    PRIMARY KEY (block_number, tx_index, log_index)
);

CREATE INDEX IF NOT EXISTS exchange_event_buy_token_pkey
    ON exchange_event_buy_token (block_number, tx_index, log_index);

COMMIT;
