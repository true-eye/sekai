package gov

import (
	"github.com/KiraCore/sekai/x/gov/keeper"
	"github.com/KiraCore/sekai/x/gov/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ApplyAssignPermissionProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyAssignPermissionProposalHandler(keeper keeper.Keeper) *ApplyAssignPermissionProposalHandler {
	return &ApplyAssignPermissionProposalHandler{keeper: keeper}
}

func (a ApplyAssignPermissionProposalHandler) ProposalType() string {
	return types.AssignPermissionProposalType
}

func (a ApplyAssignPermissionProposalHandler) Apply(ctx sdk.Context, proposal types.Content) {
	p := proposal.(*types.AssignPermissionProposal)

	actor, found := a.keeper.GetNetworkActorByAddress(ctx, p.Address)
	if !found {
		actor = types.NewDefaultActor(p.Address)
	}

	err := a.keeper.AddWhitelistPermission(ctx, actor, types.PermValue(p.Permission))
	if err != nil {
		panic("network actor has this permission")
	}
}

type ApplySetNetworkPropertyProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplySetNetworkPropertyProposalHandler(keeper keeper.Keeper) *ApplySetNetworkPropertyProposalHandler {
	return &ApplySetNetworkPropertyProposalHandler{keeper: keeper}
}

func (a ApplySetNetworkPropertyProposalHandler) ProposalType() string {
	return types.SetNetworkPropertyProposalType
}

func (a ApplySetNetworkPropertyProposalHandler) Apply(ctx sdk.Context, proposal types.Content) {
	p := proposal.(*types.SetNetworkPropertyProposal)

	err := a.keeper.SetNetworkProperty(ctx, p.NetworkProperty, p.Value)
	if err != nil {
		panic("error setting network property")
	}
}

type ApplyUpsertDataRegistryProposalHandler struct {
	keeper keeper.Keeper
}

func NewApplyUpsertDataRegistryProposalHandler(keeper keeper.Keeper) *ApplyUpsertDataRegistryProposalHandler {
	return &ApplyUpsertDataRegistryProposalHandler{keeper: keeper}
}

func (a ApplyUpsertDataRegistryProposalHandler) ProposalType() string {
	return types.UpsertDataRegistryProposalType
}

func (a ApplyUpsertDataRegistryProposalHandler) Apply(ctx sdk.Context, proposal types.Content) {
	p := proposal.(*types.UpsertDataRegistryProposal)
	entry := types.NewDataRegistryEntry(p.Hash, p.Reference, p.Encoding, p.Size_)
	a.keeper.UpsertDataRegistryEntry(ctx, p.Key, entry)
}