package bnc

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"nft-backend/contracts"
	"nft-backend/internal/entities"
	"nft-backend/internal/repositories"
	"strings"
	"time"
)

func (c *Client) connect() {
	client, err := ethclient.Dial(c.config.NetworkEndpoint)
	if err != nil {
		c.conn <- err
		c.zapLogger.Sugar().Error(err)
		return
	}

	chainId, err := client.ChainID(context.Background())
	if err != nil {
		c.conn <- err
		c.zapLogger.Sugar().Error(err)
		return
	}
	c.zapLogger.Sugar().Info(chainId)

	//c.zapLogger.Sugar().Debug(crypto.Keccak256Hash(c.config.ExchangeBuyEventSignature))

	nftContractAddress := common.HexToAddress(c.config.NFTContractAddress)
	exchangeContractAddress := common.HexToAddress(c.config.ExchangeContractAddress)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			nftContractAddress,
			exchangeContractAddress,
		},
	}

	sub, err := client.SubscribeFilterLogs(context.Background(), query, c.logs)
	if err != nil {
		c.conn <- err
		c.zapLogger.Sugar().Error(err)
		return
	}
	nftContract, err := contracts.NewNFT(nftContractAddress, client)
	if err != nil {
		c.conn <- err
		c.zapLogger.Sugar().Error(err)
		return
	}
	exchangeContract, err := contracts.NewExchange(exchangeContractAddress, client)
	if err != nil {
		c.conn <- err
		c.zapLogger.Sugar().Error(err)
		return
	}

	watchOpts := &bind.WatchOpts{
		Start:   nil,
		Context: context.Background(),
	}

	buyTokenEvents := make(chan *contracts.ExchangeBuyToken)
	buyTokenEventSub, err := exchangeContract.WatchBuyToken(watchOpts, buyTokenEvents, nil, nil, nil)
	if err != nil {
		c.conn <- err
		c.zapLogger.Sugar().Error(err)
		return
	}

	sellTokenEvents := make(chan *contracts.ExchangeSellToken)
	sellTokenEventSub, err := exchangeContract.WatchSellToken(watchOpts, sellTokenEvents, nil, nil, nil)
	if err != nil {
		c.conn <- err
		c.zapLogger.Sugar().Error(err)
		return
	}

	/*go func() {
		time.Sleep(3*time.Second)
		client, err := ethclient.Dial(c.config.NetworkEndpoint)
		if err != nil {
			c.zapLogger.Sugar().Error(err)
			return
		}
		nftContract, err := contracts.NewNFT(nftContractAddress, client)
		if err != nil {
			c.zapLogger.Sugar().Error(err)
			return
		}
		exchangeContract, err := contracts.NewExchange(exchangeContractAddress, client)
		if err != nil {
			c.zapLogger.Sugar().Error(err)
			return
		}
		c.zapLogger.Sugar().Debug("start test")

		account1 := common.HexToAddress("0x177E42b7A1a124DaC17B61E1f9A30e3121067eF4")
		//account2 := common.HexToAddress("0xE2Fc0f5AF5C7B585CC74F14DF0A9e568E23a3960")

		privateKey1, err := crypto.HexToECDSA("00ec2ba4b2a0caa8b6d2a2dbde31318769dd7000b6a81df5f92ec73b43408d6c")
		if err != nil {
			c.zapLogger.Sugar().Debug(err)
			return
		}
		privateKey2, err := crypto.HexToECDSA("10d82453047626d3a7da10a47b6c7abcaf22de865eb400ae729aa0ac1cd034b1")
		if err != nil {
			c.zapLogger.Sugar().Debug(err)
			return
		}
		auth1, err := bind.NewKeyedTransactorWithChainID(privateKey1, chainId)
		if err != nil {
			c.zapLogger.Sugar().Debug(err)
			return
		}
		auth2, err := bind.NewKeyedTransactorWithChainID(privateKey2, chainId)
		if err != nil {
			c.zapLogger.Sugar().Debug(err)
			return
		}

		_, err = nftContract.Create(auth1, account1, "integration-test")
		if err != nil {
			c.zapLogger.Sugar().Debug(err)
			return
		}

		transferEvents := make(chan *contracts.NFTTransfer, 1)
		go func() {
			c.zapLogger.Sugar().Debug("WatchTransfer")
			sub, err := nftContract.WatchTransfer(watchOpts, transferEvents, nil, nil, nil)
			if err != nil {
				c.zapLogger.Sugar().Debug(err)
				return
			}
			select {
			case err := <-sub.Err():
				c.zapLogger.Sugar().Debug(err)
				return
			}
			//defer sub.Unsubscribe()
		}()
		tokenId := (<-transferEvents).TokenId
		c.zapLogger.Sugar().Debug(tokenId)

		sellNftPrice := contracts.SimpleExchangeNFTNftPrice{
			Price: decimal.NewFromFloat(1).BigInt(),
			Token: 0,
		}

		c.zapLogger.Sugar().Debug("SellToken")
		_, err = exchangeContract.SellToken(auth1, tokenId, sellNftPrice)
		if err != nil {
			c.zapLogger.Sugar().Debug(err)
			return
		}

		sellTokenEvent := make(chan *contracts.ExchangeSellToken)
		go func() {
			c.zapLogger.Sugar().Debug("WatchSellToken")
			sub, err := exchangeContract.WatchSellToken(watchOpts, sellTokenEvent, nil, nil, nil)
			if err != nil {
				return
			}
			select {
			case err := <-sub.Err():
				c.zapLogger.Sugar().Debug(err)
				return
			}
			//defer sub.Unsubscribe()
		}()
		<-sellTokenEvent

		c.zapLogger.Sugar().Debug("Approve")
		approveTx, err := nftContract.Approve(auth1, exchangeContractAddress, tokenId)
		if err != nil {
			c.zapLogger.Sugar().Debug(err)
			return
		}
		for {
			_, isPending, err := client.TransactionByHash(context.Background(), approveTx.Hash())
			if err != nil {
				c.zapLogger.Sugar().Debug(err)
				break
			}
			if !isPending {
				break
			}
		}

		c.zapLogger.Sugar().Debug("BuyToken")
		auth2.Value = decimal.NewFromFloat(1).BigInt()
		_, err = exchangeContract.BuyToken(auth2, tokenId)
		if err != nil {
			c.zapLogger.Sugar().Debug(err)
			return
		}

		buyTokenEvent := make(chan *contracts.ExchangeBuyToken)
		go func() {
			c.zapLogger.Sugar().Debug("WatchBuyToken")
			sub, err := exchangeContract.WatchBuyToken(watchOpts, buyTokenEvent, nil, nil, nil)
			if err != nil {
				return
			}
			select {
			case err := <-sub.Err():
				c.zapLogger.Sugar().Debug(err)
				return
			}
			//defer sub.Unsubscribe()
		}()
		c.zapLogger.Sugar().Debug(<-buyTokenEvent)
	}()*/

	c.conn <- nil

	c.zapLogger.Sugar().Info(fmt.Sprintf("connected to nft contract: %s", nftContractAddress))
	c.zapLogger.Sugar().Info(fmt.Sprintf("connected to nft exchange contract: %s", exchangeContractAddress))

	for {
		select {
		case err := <-sub.Err():
			c.zapLogger.Sugar().Error(err)
			select {
			case c.conn <- err:
			case <-time.After(5 * time.Second):
				return
			}
			return
		case err := <-buyTokenEventSub.Err():
			c.zapLogger.Sugar().Error(err)
			select {
			case c.conn <- err:
			case <-time.After(5 * time.Second):
				return
			}
			return
		case err := <-sellTokenEventSub.Err():
			c.zapLogger.Sugar().Error(err)
			select {
			case c.conn <- err:
			case <-time.After(5 * time.Second):
				return
			}
			return
		case event := <-sellTokenEvents:
			//log := <-c.logs
			c.zapLogger.Sugar().Debug("sell:", event)
			//c.zapLogger.Sugar().Debug("log:", log)
			//c.zapLogger.Sugar().Debug("seller", common.HexToAddress(log.Topics[2].Hex()).String())

			block, err := client.BlockByNumber(context.Background(), decimal.NewFromInt(int64(event.Raw.BlockNumber)).BigInt())
			if err != nil {
				c.zapLogger.Sugar().Error(err)
				continue
			}

			owner := (&repositories.CollectibleRepository{}).GetOwnerAddress(context.Background(), c.pgPool, decimal.NewFromBigInt(event.TokenId, 0))
			creator := (&repositories.CollectibleRepository{}).GetCreatorAddress(context.Background(), c.pgPool, decimal.NewFromBigInt(event.TokenId, 0))
			collectible_id := (&repositories.CollectibleRepository{}).GetIdByTokenId(context.Background(), c.pgPool, decimal.NewFromBigInt(event.TokenId, 0))

			sellTokenEventEnt := &entities.TokenEvent{
				TokenId: decimal.NewFromBigInt(event.TokenId, 0),
				NFTPrice: entities.NFTPrice{
					Token: int16(event.Token),
					Price: decimal.NewFromBigInt(event.Price, 0),
				},
				Block: entities.Block{
					Number:    event.Raw.BlockNumber,
					Hash:      event.Raw.BlockHash.String(),
					Timestamp: time.Unix(int64(block.Time()), 0),
				},
				Tx: entities.Tx{
					Hash:  event.Raw.TxHash.String(),
					Index: event.Raw.TxIndex,
				},
				Log: entities.Log{
					Index:   event.Raw.Index,
					Address: event.Raw.Address.String(),
					Data:    event.Raw.Data,
					Removed: event.Raw.Removed,
				},
				Type:    entities.TokenEvent_Sell,
				Account: owner,
			}

			err = (&repositories.ExchangeEventRepo{}).UpsertTokenEvent(context.Background(), c.pgPool, sellTokenEventEnt)
			if err != nil {
				c.zapLogger.Sugar().Error(err)
			}

			account_owner, err := (&repositories.AccountRepository{}).GetByAddress(
				context.Background(), c.pgPool, *owner)
			var collectible_id_64 sql.NullInt64
			collectible_id_64.Int64 = collectible_id
			collectible_id_64.Valid = true

			if err != nil {
				notice_sell := &entities.Notice{
					AccountID:     account_owner.Id,
					CollectibleID: collectible_id_64,
					Content:       entities.Notice_someone_bought_nft,
				}

				err = (&repositories.NoticeRepo{}).Insert(context.Background(), c.pgPool, notice_sell)

			}
			if !strings.EqualFold(*owner, *creator) { // send noti

				account_creator, err := (&repositories.AccountRepository{}).GetByAddress(
					context.Background(), c.pgPool, *creator)

				if err != nil {

					notice_resell := &entities.Notice{
						AccountID:     account_creator.Id,
						CollectibleID: collectible_id_64,
						Content:       entities.Notice_resell_nft_done,
					}
					err = (&repositories.NoticeRepo{}).Insert(context.Background(), c.pgPool, notice_resell)
				}

			}

			err = (&repositories.CollectibleRepository{}).UpdateLock(context.Background(), c.pgPool, sellTokenEventEnt.TokenId, false)
			if err != nil {
				c.zapLogger.Sugar().Error(err)
				continue
			}

		case event := <-buyTokenEvents:

			c.zapLogger.Sugar().Debug("buy:", event)

			//c.zapLogger.Sugar().Debug("event.Price", event.Price)
			//owner := common.HexToAddress(log.Topics[2].Hex()).String()
			//c.zapLogger.Sugar().Debug("buyer:", owner)

			block, err := client.BlockByNumber(context.Background(), decimal.NewFromInt(int64(event.Raw.BlockNumber)).BigInt())
			c.zapLogger.Sugar().Debug("block", block)

			if err != nil {
				c.zapLogger.Sugar().Error(err)
				continue
			}

			owner := (&repositories.CollectibleRepository{}).GetOwnerAddress(context.Background(), c.pgPool, decimal.NewFromBigInt(event.TokenId, 0))
			creator := (&repositories.CollectibleRepository{}).GetCreatorAddress(context.Background(), c.pgPool, decimal.NewFromBigInt(event.TokenId, 0))
			collectible_id := (&repositories.CollectibleRepository{}).GetIdByTokenId(context.Background(), c.pgPool, decimal.NewFromBigInt(event.TokenId, 0))

			buyTokenEventEnt := &entities.TokenEvent{
				TokenId: decimal.NewFromBigInt(event.TokenId, 0),
				NFTPrice: entities.NFTPrice{
					Token: int16(event.Token),
					Price: decimal.NewFromBigInt(event.Price, 0),
				},
				Block: entities.Block{
					Number:    event.Raw.BlockNumber,
					Hash:      event.Raw.BlockHash.String(),
					Timestamp: time.Unix(int64(block.Time()), 0),
				},
				Tx: entities.Tx{
					Hash:  event.Raw.TxHash.String(),
					Index: event.Raw.TxIndex,
				},
				Log: entities.Log{
					Index:   event.Raw.Index,
					Address: event.Raw.Address.String(),
					Data:    event.Raw.Data,
					Removed: event.Raw.Removed,
				},
				Type:    entities.TokenEvent_Buy,
				Account: owner,
			}

			//owner := (&repositories.CollectibleRepository{}).GetOwnerAddress(context.Background(), c.pgPool, buyTokenEventEnt.TokenId)

			err = (&repositories.ExchangeEventRepo{}).UpsertTokenEvent(context.Background(), c.pgPool, buyTokenEventEnt)
			if err != nil {
				c.zapLogger.Sugar().Error(err)
			}

			total_trade := (&repositories.ExchangeEventRepo{}).CountTrade(context.Background(), c.pgPool, buyTokenEventEnt.TokenId)

			(&repositories.CollectibleRepository{}).UpdateTotalTrade(context.Background(), c.pgPool, buyTokenEventEnt.TokenId, total_trade)

			err = (&repositories.CollectibleRepository{}).UpdateLock(context.Background(), c.pgPool, buyTokenEventEnt.TokenId, true)
			if err != nil {
				c.zapLogger.Sugar().Error(err)
				continue
			}

			if !strings.EqualFold(*owner, *creator) { // send noti

				account_creator, err := (&repositories.AccountRepository{}).GetByAddress(
					context.Background(), c.pgPool, *creator)
				var collectible_id_64 sql.NullInt64
				collectible_id_64.Int64 = collectible_id
				collectible_id_64.Valid = true

				if err != nil {
					notice := &entities.Notice{
						AccountID:     account_creator.Id,
						CollectibleID: collectible_id_64,
						Content:       entities.Notice_resell_nft,
					}
					err = (&repositories.NoticeRepo{}).Insert(context.Background(), c.pgPool, notice)
				}

			}

		case log := <-c.logs:
			if len(log.Topics) < 1 {
				continue
			}
			switch log.Topics[0] {
			case crypto.Keccak256Hash(c.config.NFTCreateEventSignature):
				if len(log.Topics) != 4 {
					c.zapLogger.Sugar().Warn("contract topics len is diff 4")
					continue
				}
				tokenURI, err := nftContract.TokenURI(&bind.CallOpts{}, log.Topics[3].Big())
				if err != nil {
					c.zapLogger.Sugar().Error(err)
					continue
				}
				tokenId := decimal.NewFromBigInt(log.Topics[3].Big(), 0)
				token := log.TxHash.String()
				tokenOwner := common.HexToAddress(log.Topics[2].Hex()).String()
				c.zapLogger.Sugar().Debug("tokenOwner: ", tokenOwner)
				c.zapLogger.Sugar().Debug("tokenURI: ", tokenURI)

				/*strList := strings.Split(tokenURI, "/")
				guid := strList[len(strList)-1]
				if len(strList) > 1 {
					err := (&repositories.CollectibleRepository{}).UpdateTokenInfoByGUID(
						context.Background(), c.pgPool, guid,
						tokenId,
						token,
						tokenOwner,
					)
					if err != nil {
						c.zapLogger.Sugar().Error(err)
						continue
					}
				}*/
				err = (&repositories.CollectibleRepository{}).UpdateTokenInfoByGUID(
					context.Background(),
					c.pgPool,
					tokenURI,
					tokenId,
					token,
					tokenOwner,
				)
				if err != nil {
					c.zapLogger.Sugar().Error(err)
					continue
				}

				/*account, err := (&repositories.AccountRepository{}).GetByAddress(
					context.Background(), c.pgPool, tokenOwner)
				if err != nil {
					notice := &entities.Notice{
						AccountID: account.Id,
						Content:   entities.Notice_create_nft,
					}
					err = (&repositories.NoticeRepo{}).Insert(context.Background(), c.pgPool, notice)
				}*/

				c.zapLogger.Sugar().Infow(
					"update collectible with contract info",
					"token", token,
					"tokenId", tokenId.String(),
					"tokenOwner", tokenOwner,
				)
			case crypto.Keccak256Hash(c.config.ExchangeBuyEventSignature):
				if len(log.Topics) < 3 {
					c.zapLogger.Sugar().Warn("topics in SellEvent less than 3")
					continue
				}

			}
		case <-c.stop:
			return
		}
	}
}

func (c *Client) run() {
	ticker := time.NewTicker(3 * time.Second)
	defer func() {
		ticker.Stop()
	}()

	var tryCount int
	var ethClient *ethclient.Client

	for {
		select {
		case err := <-c.conn:
			switch err {
			case nil:
				c.zapLogger.Sugar().Info("subscribe to nft contract success")
			default:
				c.zapLogger.Sugar().Info(fmt.Sprintf("try connecting after %v second(s)", time.Duration(tryCount)*time.Second))
				time.Sleep(time.Duration(tryCount) * time.Second)
				go c.connect()
			}
		case <-c.stop:
			if ethClient != nil {
				ethClient.Close()
			}
			return
		case <-ticker.C:
			//c.zapLogger.Sugar().Info("bnc running...")
		}
	}
}

func (c *Client) Shutdown() {
	close(c.stop)
}
