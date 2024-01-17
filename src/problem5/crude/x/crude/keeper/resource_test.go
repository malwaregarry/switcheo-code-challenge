package keeper_test

import (
	"context"
	"testing"

	keepertest "crude/testutil/keeper"
	"crude/testutil/nullify"
	"crude/x/crude/keeper"
	"crude/x/crude/types"
	"github.com/stretchr/testify/require"
)

func createNResource(keeper keeper.Keeper, ctx context.Context, n int) []types.Resource {
	items := make([]types.Resource, n)
	for i := range items {
		items[i].Id = keeper.AppendResource(ctx, items[i])
	}
	return items
}

func TestResourceGet(t *testing.T) {
	keeper, ctx := keepertest.CrudeKeeper(t)
	items := createNResource(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetResource(ctx, item.Id)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestResourceRemove(t *testing.T) {
	keeper, ctx := keepertest.CrudeKeeper(t)
	items := createNResource(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveResource(ctx, item.Id)
		_, found := keeper.GetResource(ctx, item.Id)
		require.False(t, found)
	}
}

func TestResourceGetAll(t *testing.T) {
	keeper, ctx := keepertest.CrudeKeeper(t)
	items := createNResource(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllResource(ctx)),
	)
}

func TestResourceCount(t *testing.T) {
	keeper, ctx := keepertest.CrudeKeeper(t)
	items := createNResource(keeper, ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, keeper.GetResourceCount(ctx))
}
