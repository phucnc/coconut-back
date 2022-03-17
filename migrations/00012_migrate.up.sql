BEGIN;

CREATE TABLE IF NOT EXISTS exchange_event_sell_token
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
    created_at      timestamptz     NOT NULL DEFAULT now(),
    updated_at      timestamptz     NOT NULL DEFAULT now(),
    deleted_at      timestamptz              DEFAULT NULL,
    PRIMARY KEY (block_number, tx_index, log_index)
);

CREATE INDEX IF NOT EXISTS exchange_event_sell_token_pkey
    ON exchange_event_sell_token (block_number, tx_index, log_index);
CREATE INDEX IF NOT EXISTS exchange_event_sell_token__nft_token_id__idx
    ON exchange_event_sell_token (nft_token_id);
CREATE INDEX IF NOT EXISTS exchange_event_sell_token__nft_price__idx
    ON exchange_event_sell_token (nft_price);
CREATE INDEX IF NOT EXISTS exchange_event_sell_token__nft_quote_token__idx
    ON exchange_event_sell_token (nft_quote_token);
CREATE INDEX IF NOT EXISTS exchange_event_sell_token__block_number__asc__idx
    ON exchange_event_sell_token (block_number ASC );
CREATE INDEX IF NOT EXISTS exchange_event_sell_token__block_number__desc__idx
    ON exchange_event_sell_token (block_number DESC );
CREATE INDEX IF NOT EXISTS exchange_event_sell_token__deleted_at__idx
    ON exchange_event_sell_token (deleted_at);


CREATE INDEX IF NOT EXISTS collectible__deleted_at__idx
    ON collectible(deleted_at);
CREATE INDEX IF NOT EXISTS collectible_category__deleted_at__idx
    ON collectible_category(deleted_at);


ALTER TABLE public.exchange_event_buy_token
    ADD COLUMN IF NOT EXISTS created_at timestamptz NOT NULL DEFAULT now();
ALTER TABLE public.exchange_event_buy_token
    ADD COLUMN IF NOT EXISTS updated_at timestamptz NOT NULL DEFAULT now();
ALTER TABLE public.exchange_event_buy_token
    ADD COLUMN IF NOT EXISTS deleted_at timestamptz DEFAULT NULL;

CREATE INDEX IF NOT EXISTS exchange_event_buy_token__deleted_at__idx
    ON exchange_event_buy_token (deleted_at);

COMMIT;
