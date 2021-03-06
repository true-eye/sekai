package keeper

import (
	"github.com/KiraCore/sekai/x/staking/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper represents the keeper that maintains the Validator Registry.
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.LegacyAmino
}

// NewKeeper returns new keeper.
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.LegacyAmino) Keeper {
	return Keeper{storeKey: storeKey, cdc: cdc}
}

func (k Keeper) AddValidator(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&validator)
	store.Set(types.GetValidatorKey(validator.ValKey), bz)

	// Save by moniker
	store.Set(types.GetValidatorByMonikerKey(validator.Moniker), types.GetValidatorKey(validator.ValKey))
}

func (k Keeper) GetValidator(ctx sdk.Context, address sdk.ValAddress) types.Validator {
	return k.getValidatorByKey(ctx, types.GetValidatorKey(address))
}

func (k Keeper) GetValidatorByAccAddress(ctx sdk.Context, address sdk.AccAddress) types.Validator {
	return k.getValidatorByKey(ctx, types.GetValidatorKeyAcc(address))
}

func (k Keeper) GetValidatorByMoniker(ctx sdk.Context, moniker string) types.Validator {
	store := ctx.KVStore(k.storeKey)

	valKey := store.Get(types.GetValidatorByMonikerKey(moniker))

	return k.getValidatorByKey(ctx, valKey)
}

func (k Keeper) getValidatorByKey(ctx sdk.Context, key []byte) types.Validator {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(key)

	var validator types.Validator
	k.cdc.MustUnmarshalBinaryBare(bz, &validator)

	return validator
}

func (k Keeper) GetValidatorSet(ctx sdk.Context) []types.Validator {
	store := ctx.KVStore(k.storeKey)

	iter := sdk.KVStorePrefixIterator(store, types.ValidatorsKey)
	defer iter.Close()

	var validators []types.Validator
	for ; iter.Valid(); iter.Next() {
		var validator types.Validator
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &validator)
		validators = append(validators, validator)
	}

	return validators
}
