package bank

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/shanev/cosmos-record-keeper/recordkeeper"
)

// Keeper is the model object for the package bank module
type Keeper struct {
	RecordKeeper
	storeKey   sdk.StoreKey
	codec      *codec.Codec
	paramStore params.Subspace
	bankKeeper bank.Keeper
	codespace  sdk.CodespaceType
}

// NewKeeper creates a bank keeper.
func NewKeeper(codec *codec.Codec, storeKey sdk.StoreKey, bankKeeper bank.Keeper,
	paramStore params.Subspace, codespace sdk.CodespaceType) Keeper {
	return Keeper{
		RecordKeeper: recordkeeper.NewRecordKeeper(storeKey, codec),
		storeKey:     storeKey,
		codec:        codec,
		bankKeeper:   bankKeeper,
		paramStore:   paramStore.WithKeyTable(ParamKeyTable()),
		codespace:    codespace,
	}
}

// Codespace returns the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// Transaction gets a specific transaction by id.
func (k Keeper) Transaction(ctx sdk.Context, id uint64) (Transaction, sdk.Error) {
	tx := Transaction{}
	err := k.Get(ctx, id, &tx)
	return tx, err
}

// AddCoin adds a coin to an address and adds the transaction to the association list.
func (k Keeper) AddCoin(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin,
	referenceID uint64, txType TransactionType) (sdk.Coins, sdk.Error) {
	if !txType.AllowedForAddition() {
		return sdk.Coins{}, ErrInvalidTransactionType(txType)
	}
	coins, err := k.bankKeeper.AddCoins(ctx, addr, sdk.Coins{amt})
	if err != nil {
		return coins, err
	}

	tx := Transaction{
		ID:                k.IncrementID(ctx),
		Type:              txType,
		ReferenceID:       referenceID,
		Amount:            amt,
		AppAccountAddress: addr,
		CreatedTime:       ctx.BlockHeader().Time,
	}
	k.Set(ctx, tx.ID, tx)
	k.PushWithAddress(ctx, k.storeKey, accountKey, tx.ID, addr)
	return coins, nil
}

// SubtractCoin subtracts a coin from an address and adds the transaction to the association list.
func (k Keeper) SubtractCoin(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin,
	referenceID uint64, txType TransactionType) (sdk.Coins, sdk.Error) {
	if !txType.AllowedForDeduction() {
		return sdk.Coins{}, ErrInvalidTransactionType(txType)
	}
	coins, err := k.bankKeeper.SubtractCoins(ctx, addr, sdk.Coins{amt})
	if err != nil {
		return coins, err
	}

	tx := Transaction{
		ID:                k.IncrementID(ctx),
		Type:              txType,
		ReferenceID:       referenceID,
		Amount:            amt,
		AppAccountAddress: addr,
		CreatedTime:       ctx.BlockHeader().Time,
	}
	k.Set(ctx, tx.ID, tx)
	k.PushWithAddress(ctx, k.storeKey, accountKey, tx.ID, addr)
	return coins, nil
}

func (k Keeper) GetCoins(ctx sdk.Context, address sdk.AccAddress) sdk.Coins {
	return k.bankKeeper.GetCoins(ctx, address)
}

func (k Keeper) rewardBrokerAddress(ctx sdk.Context) sdk.AccAddress {
	address := sdk.AccAddress{}
	k.paramStore.GetIfExists(ctx, ParamKeyRewardBrokerAddress, &address)
	return address
}

func (k Keeper) payReward(ctx sdk.Context,
	sender sdk.AccAddress, recipient sdk.AccAddress,
	amount sdk.Coin, inviteID uint64) sdk.Error {
	if !k.rewardBrokerAddress(ctx).Equals(sender) {
		return ErrInvalidRewardBrokerAddress(sender)
	}
	_, err := k.AddCoin(ctx, recipient, amount, inviteID, TransactionRewardPayout)
	if err != nil {
		return err
	}
	return nil
}

// Transactions gets all the transactions
func (k Keeper) Transactions(ctx sdk.Context) []Transaction {
	txs := make([]Transaction, 0)
	err := k.Each(ctx, func(val []byte) bool {
		var tx Transaction
		k.codec.MustUnmarshalBinaryLengthPrefixed(val, &tx)
		txs = append(txs, tx)
		return true
	})
	if err != nil {
		return nil
	}
	return txs
}

// TransactionsByAddress gets transactions for a given address and applies sent filters.
func (k Keeper) TransactionsByAddress(ctx sdk.Context, address sdk.AccAddress, filterSetters ...Filter) []Transaction {
	filters := GetFilters(filterSetters...)
	transactions := make([]Transaction, 0)
	filterByType := len(filters.TransactionTypes) > 0

	offsetCount := filters.Offset
	count := 0
	mapFunc := func(txID uint64) bool     {
		if offsetCount > 0 {
			offsetCount = offsetCount - 1
			return true
		}

		if filters.Limit > 0 && count == filters.Limit {
			return false
		}

		tx, err := k.Transaction(ctx, txID)
		if err != nil {
			panic(err)
		}
		count++
		if filterByType && !tx.Type.OneOf(filters.TransactionTypes) {
			return true
		}
		transactions = append(transactions, tx)
		return true
	}
	if filters.SortOrder == SortDesc {
		k.ReverseMapByAddress(ctx, accountKey, address, mapFunc)
		return transactions
	}
	k.MapByAddress(ctx, accountKey, address, mapFunc)
	return transactions
}
