package keeper

import (
	"context"

	"chainTools/x/tools/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k Keeper) TotalWallet(goCtx context.Context, req *types.QueryTotalWalletRequest) (*types.QueryTotalWalletResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	bank := k.bankKeeper.GetAllBalances(ctx, sdk.AccAddress(req.Address))

	ubonding := sdk.Coins{}
	for _, unbond := range k.stakingKeeper.GetAllUnbondingDelegations(ctx, sdk.AccAddress(req.Address)) {
		for _, v := range unbond.Entries {
			ubonding.Add(sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), v.Balance))
		}
	}

	delegated := sdk.Coins{}
	rewards := sdk.Coins{}
	for _, delegation := range k.stakingKeeper.GetDelegatorDelegations(ctx, sdk.AccAddress(req.Address), 50000) {
		val, found := k.stakingKeeper.GetValidator(ctx, delegation.GetValidatorAddr())
		if !found {
			return nil, stakingtypes.ErrNoValidatorFound
		}

		endingPeriod := k.distrKeeper.IncrementValidatorPeriod(ctx, val)

		delegated.Add(sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), val.TokensFromShares(delegation.Shares).TruncateInt()))
		for _, c := range k.distrKeeper.CalculateDelegationRewards(ctx, val, delegation, endingPeriod) {
			rewards.Add(sdk.NewCoin(c.Denom, c.Amount.TruncateInt()))
		}
	}

	total := sdk.Coins{}
	total.Add(bank...)
	total.Add(ubonding...)
	total.Add(delegated...)
	total.Add(rewards...)

	return &types.QueryTotalWalletResponse{
		Bank:      bank,
		Unbond:    ubonding,
		Delegated: delegated,
		Rewards:   rewards,
		Total:     total,
	}, nil
}
