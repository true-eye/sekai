package types

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgClaimValidator{}

func NewMsgClaimValidator(
	moniker string,
	website string,
	social string,
	identity string,
	comission sdk.Dec,
	valKey sdk.ValAddress,
	pubKey crypto.PubKey,
) (*MsgClaimValidator, error) {
	if valKey == nil {
		return nil, fmt.Errorf("validator not set")
	}

	if pubKey == nil {
		return nil, fmt.Errorf("public key not set")
	}

	pkStr := sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, pubKey)

	return &MsgClaimValidator{
		Moniker:    moniker,
		Website:    website,
		Social:     social,
		Identity:   identity,
		Commission: comission,
		ValKey:     valKey,
		PubKey:     pkStr,
	}, nil
}

func (m MsgClaimValidator) Route() string {
	return ModuleName
}

func (m MsgClaimValidator) Type() string {
	return ClaimValidator
}

func (m MsgClaimValidator) ValidateBasic() error {
	return nil
}

func (m MsgClaimValidator) GetSignBytes() []byte {
	//bz := ModuleCdc.MustMarshalJSON(m)
	//return sdk.MustSortJSON(bz)
	return nil
}

func (m MsgClaimValidator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		sdk.AccAddress(m.ValKey),
	}
}
