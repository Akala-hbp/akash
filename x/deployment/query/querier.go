package query

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ovrclk/akash/sdkutil"
	"github.com/ovrclk/akash/x/deployment/keeper"
	"github.com/ovrclk/akash/x/deployment/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier creates and returns a new deployment querier instance
func NewQuerier(keeper keeper.Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case deploymentsPath:
			return queryDeployments(ctx, path[1:], req, keeper, legacyQuerierCdc)
		case deploymentPath:
			return queryDeployment(ctx, path[1:], req, keeper, legacyQuerierCdc)
		case groupPath:
			return queryGroup(ctx, path[1:], req, keeper, legacyQuerierCdc)
		}
		return []byte{}, sdkerrors.ErrUnknownRequest
	}
}

func queryDeployments(ctx sdk.Context, path []string, _ abci.RequestQuery, keeper keeper.Keeper,
	legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	// isValidState denotes whether given state flag is valid or not
	filters, isValidState, err := parseDepFiltersPath(path)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInternal, err.Error())
	}

	var values Deployments
	keeper.WithDeployments(ctx, func(deployment types.Deployment) bool {
		if filters.Accept(deployment, isValidState) {
			value := Deployment{
				Deployment: deployment,
				Groups:     keeper.GetGroups(ctx, deployment.ID()),
			}
			values = append(values, value)
		}

		return false
	})

	return sdkutil.RenderQueryResponse(legacyQuerierCdc, values)
}

func queryDeployment(ctx sdk.Context, path []string, _ abci.RequestQuery, keeper keeper.Keeper,
	legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

	id, err := ParseDeploymentPath(path)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInternal, err.Error())
	}

	deployment, ok := keeper.GetDeployment(ctx, id)
	if !ok {
		return nil, types.ErrDeploymentNotFound
	}

	value := Deployment{
		Deployment: deployment,
		Groups:     keeper.GetGroups(ctx, deployment.ID()),
	}

	return sdkutil.RenderQueryResponse(legacyQuerierCdc, value)
}

func queryGroup(ctx sdk.Context, path []string, _ abci.RequestQuery, keeper keeper.Keeper,
	legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

	id, err := ParseGroupPath(path)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "internal error")
	}

	group, ok := keeper.GetGroup(ctx, id)
	if !ok {
		return nil, sdkerrors.Wrap(err, "group not found")
	}

	value := Group(group)

	return sdkutil.RenderQueryResponse(legacyQuerierCdc, value)
}
