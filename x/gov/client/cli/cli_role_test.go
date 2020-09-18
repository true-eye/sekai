package cli_test

import (
	"context"
	"fmt"

	"github.com/KiraCore/sekai/x/gov/client/cli"
	customgovtypes "github.com/KiraCore/sekai/x/gov/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	types3 "github.com/cosmos/cosmos-sdk/types"
)

func (s IntegrationTestSuite) TestWhitelistRolePermission() {
	// Query permissions for role Validator
	val := s.network.Validators[0]

	cmd := cli.GetCmdQueryRolePermissions()
	_, out := testutil.ApplyMockIO(cmd)
	clientCtx := val.ClientCtx.WithOutput(out).WithOutputFormat("json")

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	cmd.SetArgs([]string{
		"0", // RoleInTest
	})

	err := cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	var perms customgovtypes.Permissions
	val.ClientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &perms)

	s.Require().False(perms.IsWhitelisted(customgovtypes.PermSetPermissions))

	// Send Tx To Whitelist permission
	out.Reset()

	cmd = cli.GetTxWhitelistRolePermission()
	cmd.SetArgs([]string{
		"0", // Role created in test
		"2", // PermSetPermission
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, types3.NewCoins(types3.NewCoin(s.cfg.BondDenom, types3.NewInt(10))).String()),
	})

	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	// Query again to check if it has the new permission
	out.Reset()

	cmd = cli.GetCmdQueryRolePermissions()

	cmd.SetArgs([]string{
		"0", // RoleCreatedInTest
	})

	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	var newPerms customgovtypes.Permissions
	val.ClientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &newPerms)

	s.Require().True(newPerms.IsWhitelisted(customgovtypes.PermSetPermissions))
}

func (s IntegrationTestSuite) TestBlacklistRolePermission() {
	// Query permissions for role Validator
	val := s.network.Validators[0]

	cmd := cli.GetCmdQueryRolePermissions()
	_, out := testutil.ApplyMockIO(cmd)
	clientCtx := val.ClientCtx.WithOutput(out).WithOutputFormat("json")

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	cmd.SetArgs([]string{
		"2", // RoleValidator
	})

	err := cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	var perms customgovtypes.Permissions
	val.ClientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &perms)

	s.Require().True(perms.IsWhitelisted(customgovtypes.PermClaimValidator))
	s.Require().False(perms.IsBlacklisted(customgovtypes.PermClaimCouncilor))

	// Send Tx To Blacklist permission
	out.Reset()

	cmd = cli.GetTxBlacklistRolePermission()
	cmd.SetArgs([]string{
		"2", // RoleValidator
		"3", // PermClaimCouncilor
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, types3.NewCoins(types3.NewCoin(s.cfg.BondDenom, types3.NewInt(10))).String()),
	})

	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	// Query again to check if it has the new permission
	out.Reset()

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	cmd = cli.GetCmdQueryRolePermissions()

	cmd.SetArgs([]string{
		"2", // RoleValidator
	})

	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	var newPerms customgovtypes.Permissions
	val.ClientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &newPerms)

	s.Require().True(newPerms.IsWhitelisted(customgovtypes.PermClaimValidator))
	s.Require().True(newPerms.IsBlacklisted(customgovtypes.PermClaimCouncilor))
}

func (s IntegrationTestSuite) TestRemoveWhitelistRolePermission() {
	// Query permissions for role Validator
	val := s.network.Validators[0]

	cmd := cli.GetCmdQueryRolePermissions()
	_, out := testutil.ApplyMockIO(cmd)
	clientCtx := val.ClientCtx.WithOutput(out).WithOutputFormat("json")

	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &clientCtx)

	cmd.SetArgs([]string{
		"0", // RoleInTest
	})

	err := cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	var perms customgovtypes.Permissions
	val.ClientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &perms)

	s.Require().True(perms.IsWhitelisted(customgovtypes.PermClaimValidator))

	// Send Tx To Blacklist permission
	out.Reset()

	cmd = cli.GetTxRemoveWhitelistRolePermission()
	cmd.SetArgs([]string{
		"0", // RoleValidator
		"1", // PermClaimValidator
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, types3.NewCoins(types3.NewCoin(s.cfg.BondDenom, types3.NewInt(10))).String()),
	})

	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	// Query again to check if it has the new permission
	out.Reset()

	err = s.network.WaitForNextBlock()
	s.Require().NoError(err)

	cmd = cli.GetCmdQueryRolePermissions()

	cmd.SetArgs([]string{
		"0", // RoleInTest
	})

	err = cmd.ExecuteContext(ctx)
	s.Require().NoError(err)

	var newPerms customgovtypes.Permissions
	val.ClientCtx.JSONMarshaler.MustUnmarshalJSON(out.Bytes(), &newPerms)

	s.Require().False(newPerms.IsWhitelisted(customgovtypes.PermClaimValidator))
}
