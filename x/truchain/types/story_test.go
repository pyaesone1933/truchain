package types

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestValidSubmitStoryMsg(t *testing.T) {
	validBody := "This is a valid story body @shanev amirite?"
	validCategory := DEX
	validCreator := sdk.AccAddress([]byte{1, 2})
	validStoryType := Default
	msg := NewSubmitStoryMsg(validBody, validCategory, validCreator, validStoryType)
	err := msg.ValidateBasic()

	assert.Nil(t, err)
	assert.Equal(t, "truchain", msg.Type())
	assert.Equal(t, "submit_story", msg.Name())
	assert.Equal(
		t,
		`{"body":"This is a valid story body @shanev amirite?","category":3,"creator":"cosmos1qypq36vzru","story_type":0}`,
		fmt.Sprintf("%s", msg.GetSignBytes()),
	)
	assert.Equal(t, []sdk.AccAddress{validCreator}, msg.GetSigners())
}

func TestInValidBodySubmitStoryMsg(t *testing.T) {
	invalidBody := ""
	validCategory := DEX
	validCreator := sdk.AccAddress([]byte{1, 2})
	validStoryType := Default
	msg := NewSubmitStoryMsg(invalidBody, validCategory, validCreator, validStoryType)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(702), err.Code(), err.Error())
}

func TestInValidCreatorSubmitStoryMsg(t *testing.T) {
	validBody := "This is a valid story body @shanev amirite?"
	validCategory := DEX
	invalidCreator := sdk.AccAddress([]byte{})
	validStoryType := Default
	msg := NewSubmitStoryMsg(validBody, validCategory, invalidCreator, validStoryType)
	err := msg.ValidateBasic()

	assert.Equal(t, sdk.CodeType(7), err.Code(), err.Error())
}