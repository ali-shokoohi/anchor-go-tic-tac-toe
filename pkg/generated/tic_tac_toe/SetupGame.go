// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package tic_tac_toe

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// SetupGame is the `setupGame` instruction.
type SetupGame struct {
	PlayerTwo *ag_solanago.PublicKey

	// [0] = [WRITE, SIGNER] game
	//
	// [1] = [WRITE, SIGNER] playerOne
	//
	// [2] = [] systemProgram
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewSetupGameInstructionBuilder creates a new `SetupGame` instruction builder.
func NewSetupGameInstructionBuilder() *SetupGame {
	nd := &SetupGame{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 3),
	}
	return nd
}

// SetPlayerTwo sets the "playerTwo" parameter.
func (inst *SetupGame) SetPlayerTwo(playerTwo ag_solanago.PublicKey) *SetupGame {
	inst.PlayerTwo = &playerTwo
	return inst
}

// SetGameAccount sets the "game" account.
func (inst *SetupGame) SetGameAccount(game ag_solanago.PublicKey) *SetupGame {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(game).WRITE().SIGNER()
	return inst
}

// GetGameAccount gets the "game" account.
func (inst *SetupGame) GetGameAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetPlayerOneAccount sets the "playerOne" account.
func (inst *SetupGame) SetPlayerOneAccount(playerOne ag_solanago.PublicKey) *SetupGame {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(playerOne).WRITE().SIGNER()
	return inst
}

// GetPlayerOneAccount gets the "playerOne" account.
func (inst *SetupGame) GetPlayerOneAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetSystemProgramAccount sets the "systemProgram" account.
func (inst *SetupGame) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *SetupGame {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "systemProgram" account.
func (inst *SetupGame) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

func (inst SetupGame) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_SetupGame,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst SetupGame) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *SetupGame) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.PlayerTwo == nil {
			return errors.New("PlayerTwo parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.Game is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.PlayerOne is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
	}
	return nil
}

func (inst *SetupGame) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("SetupGame")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=1]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("PlayerTwo", *inst.PlayerTwo))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=3]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("         game", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("    playerOne", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("systemProgram", inst.AccountMetaSlice.Get(2)))
					})
				})
		})
}

func (obj SetupGame) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `PlayerTwo` param:
	err = encoder.Encode(obj.PlayerTwo)
	if err != nil {
		return err
	}
	return nil
}
func (obj *SetupGame) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `PlayerTwo`:
	err = decoder.Decode(&obj.PlayerTwo)
	if err != nil {
		return err
	}
	return nil
}

// NewSetupGameInstruction declares a new SetupGame instruction with the provided parameters and accounts.
func NewSetupGameInstruction(
	// Parameters:
	playerTwo ag_solanago.PublicKey,
	// Accounts:
	game ag_solanago.PublicKey,
	playerOne ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey) *SetupGame {
	return NewSetupGameInstructionBuilder().
		SetPlayerTwo(playerTwo).
		SetGameAccount(game).
		SetPlayerOneAccount(playerOne).
		SetSystemProgramAccount(systemProgram)
}
