package main

import (
	"context"
	"log"

	"github.com/ali-shokoohi/anchor-go-tic-tac-toe/pkg/generated/tic_tac_toe"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
)

// TODO: Getting ProgramID from config
const ProgramID = "43tJ4dM7NL1zTnTmfLQmGxSAh4XHE6iHzuwEaDK8xGDG"

type KeyPairs struct {
	GamePrivateKey      solana.PrivateKey
	PlayerOnePrivateKey solana.PrivateKey
	PlayerTwoPrivateKey solana.PrivateKey
}

func main() {
	// Create a new RPC client:
	cluster := rpc.DevNet_RPC
	client := rpc.New(cluster)

	_ = client

	keyPairs, err := loadKeyPairs()
	if err != nil {
		log.Fatal("Failed at getting key pairs: ", err)
		return
	}

	tSetupGame := tic_tac_toe.NewSetupGameInstruction(
		keyPairs.PlayerTwoPrivateKey.PublicKey(),
		keyPairs.GamePrivateKey.PublicKey(),
		keyPairs.PlayerOnePrivateKey.PublicKey(),
		keyPairs.GamePrivateKey.PublicKey(),
	)
	tile := tic_tac_toe.Tile{Row: 3, Column: 3}
	tPlay := tic_tac_toe.NewPlayInstruction(tile, tSetupGame.GetGameAccount().PublicKey, *tSetupGame.PlayerTwo)
	_ = tPlay
	programID, err := solana.PublicKeyFromBase58(ProgramID)
	if err != nil {
		log.Fatal("Failed at getting programID: ", err)
	}
	tic_tac_toe.SetProgramID(programID)

	recent, err := client.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		log.Fatal("Failed at getting recent block Hash: ", err)
		return
	}
	tx, err := solana.NewTransaction([]solana.Instruction{tSetupGame.Build()}, recent.Value.Blockhash, solana.TransactionPayer(keyPairs.PlayerOnePrivateKey.PublicKey()))
	if err != nil {
		log.Fatal("Failed at calling SetupGame transaction: ", err)
		return
	}
	log.Println("SetupGame tx:", tx)
}

func loadKeyPairs() (KeyPairs, error) {
	gamePrivateKey, err := solana.PrivateKeyFromSolanaKeygenFile("./config/keypairs/game_id.json")
	if err != nil {
		err = errors.New("Failed at getting gamePrivateKey: " + err.Error())
		return KeyPairs{}, err
	}
	playerOnePrivateKey, err := solana.PrivateKeyFromSolanaKeygenFile("./config/keypairs/player1_id.json")
	if err != nil {
		err = errors.New("Failed at getting playerOnePrivateKey: " + err.Error())
		return KeyPairs{}, err
	}
	playerTwoPrivateKey, err := solana.PrivateKeyFromSolanaKeygenFile("./config/keypairs/player2_id.json")
	if err != nil {
		err = errors.New("Failed at getting playerTwoPrivateKey: " + err.Error())
		return KeyPairs{}, err
	}

	keyPairs := KeyPairs{
		GamePrivateKey:      gamePrivateKey,
		PlayerOnePrivateKey: playerOnePrivateKey,
		PlayerTwoPrivateKey: playerTwoPrivateKey,
	}
	return keyPairs, nil

}
