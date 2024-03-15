// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package tic_tac_toe

import (
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
)

type Tile struct {
	Row    uint8
	Column uint8
}

func (obj Tile) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Row` param:
	err = encoder.Encode(obj.Row)
	if err != nil {
		return err
	}
	// Serialize `Column` param:
	err = encoder.Encode(obj.Column)
	if err != nil {
		return err
	}
	return nil
}

func (obj *Tile) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Row`:
	err = decoder.Decode(&obj.Row)
	if err != nil {
		return err
	}
	// Deserialize `Column`:
	err = decoder.Decode(&obj.Column)
	if err != nil {
		return err
	}
	return nil
}

type GameState interface {
	isGameState()
}

type gameStateContainer struct {
	Enum   ag_binary.BorshEnum `borsh_enum:"true"`
	Active GameStateActive
	Tie    GameStateTie
	Won    GameStateWon
}

type GameStateActive uint8

func (obj GameStateActive) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	return nil
}

func (obj *GameStateActive) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	return nil
}

func (_ *GameStateActive) isGameState() {}

type GameStateTie uint8

func (obj GameStateTie) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	return nil
}

func (obj *GameStateTie) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	return nil
}

func (_ *GameStateTie) isGameState() {}

type GameStateWon struct {
	Winner ag_solanago.PublicKey
}

func (obj GameStateWon) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Winner` param:
	err = encoder.Encode(obj.Winner)
	if err != nil {
		return err
	}
	return nil
}

func (obj *GameStateWon) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Winner`:
	err = decoder.Decode(&obj.Winner)
	if err != nil {
		return err
	}
	return nil
}

func (_ *GameStateWon) isGameState() {}

type Sign ag_binary.BorshEnum

const (
	SignX Sign = iota
	SignO
)

func (value Sign) String() string {
	switch value {
	case SignX:
		return "X"
	case SignO:
		return "O"
	default:
		return ""
	}
}