package voting

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// store keys for voting params
var (
	KeyChallengerRewardPoolShare = []byte("challengerRewardPoolShare")
	KeyMajorityPercent           = []byte("majorityPercent")
	KeyQuorum                    = []byte("quorum")
)

// Params holds parameters for voting
type Params struct {
	ChallengerRewardPoolShare sdk.Dec
	MajorityPercent           sdk.Dec
	Quorum                    int
}

// DefaultParams is the default parameters for voting
func DefaultParams() Params {
	return Params{
		ChallengerRewardPoolShare: sdk.NewDecWithPrec(75, 2), // 75%
		MajorityPercent:           sdk.NewDecWithPrec(51, 2), // 51%
		Quorum:                    3,
	}
}

// KeyValuePairs implements params.ParamSet
func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{Key: KeyChallengerRewardPoolShare, Value: &p.ChallengerRewardPoolShare},
		{Key: KeyMajorityPercent, Value: &p.MajorityPercent},
		{Key: KeyQuorum, Value: &p.Quorum},
	}
}

// ParamTypeTable for story module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&Params{})
}

func (k Keeper) challengerRewardPoolShare(ctx sdk.Context) (res sdk.Dec) {
	k.paramStore.Get(ctx, KeyChallengerRewardPoolShare, &res)
	return
}

func (k Keeper) minQuorum(ctx sdk.Context) (res int) {
	k.paramStore.Get(ctx, KeyQuorum, &res)
	return
}

func (k Keeper) majorityPercent(ctx sdk.Context) (res sdk.Dec) {
	k.paramStore.Get(ctx, KeyMajorityPercent, &res)
	return
}

// SetParams sets the params for the expiration module
func (k Keeper) SetParams(ctx sdk.Context, params Params) {
	logger := ctx.Logger().With("module", "voting")
	k.paramStore.SetParamSet(ctx, &params)
	logger.Info(fmt.Sprintf("Loaded voting module params: %+v", params))
}