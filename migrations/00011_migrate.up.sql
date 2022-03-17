BEGIN;

CREATE INDEX IF NOT EXISTS collectible__id__asc__idx
    ON collectible(id ASC );
CREATE INDEX IF NOT EXISTS collectible__id__desc__idx
    ON collectible(id DESC );
CREATE INDEX IF NOT EXISTS collectible__title__idx
    ON collectible(title);
CREATE INDEX IF NOT EXISTS collectible__quote_token_id__idx
    ON collectible(quote_token_id);
CREATE INDEX IF NOT EXISTS collectible__token_id__idx
    ON collectible(token_id);


CREATE INDEX IF NOT EXISTS exchange_event_buy_token__nft_token_id__idx
    ON exchange_event_buy_token (nft_token_id);
CREATE INDEX IF NOT EXISTS exchange_event_buy_token__nft_price__idx
    ON exchange_event_buy_token (nft_price);
CREATE INDEX IF NOT EXISTS exchange_event_buy_token__nft_quote_token__idx
    ON exchange_event_buy_token (nft_quote_token);
CREATE INDEX IF NOT EXISTS exchange_event_buy_token__block_number__asc__idx
    ON exchange_event_buy_token (block_number ASC );
CREATE INDEX IF NOT EXISTS exchange_event_buy_token__block_number__desc__idx
    ON exchange_event_buy_token (block_number DESC );

COMMIT;
