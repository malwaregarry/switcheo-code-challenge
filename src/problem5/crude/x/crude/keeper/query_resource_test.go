package keeper_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "crude/testutil/keeper"
	"crude/testutil/nullify"
	"crude/x/crude/types"
)

func TestResourceQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.CrudeKeeper(t)
	msgs := createNResource(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetResourceRequest
		response *types.QueryGetResourceResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetResourceRequest{Id: msgs[0].Id},
			response: &types.QueryGetResourceResponse{Resource: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetResourceRequest{Id: msgs[1].Id},
			response: &types.QueryGetResourceResponse{Resource: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetResourceRequest{Id: uint64(len(msgs))},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Resource(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestResourceQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.CrudeKeeper(t)
	msgs := createNResource(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllResourceRequest {
		return &types.QueryAllResourceRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.ResourceAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Resource), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Resource),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.ResourceAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Resource), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Resource),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.ResourceAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Resource),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.ResourceAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
