package tools_test

import (
	"testing"

	keepertest "chainTools/testutil/keeper"
	"chainTools/testutil/nullify"
	"chainTools/x/tools"
	"chainTools/x/tools/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.ToolsKeeper(t)
	tools.InitGenesis(ctx, *k, genesisState)
	got := tools.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
