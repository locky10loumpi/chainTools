package keeper_test

import (
	"testing"

	testkeeper "chainTools/testutil/keeper"
	"chainTools/x/tools/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.ToolsKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
