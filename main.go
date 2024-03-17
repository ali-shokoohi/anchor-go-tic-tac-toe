package main

import (
	"context"
	"log"

	"github.com/ali-shokoohi/anchor-go-tic-tac-toe/pkg/generated/tic_tac_toe"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
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
	ctx := context.TODO()
	// Create a new RPC client:
	cluster := rpc.DevNet_RPC
	wsCluster := rpc.DevNet_WS
	client := rpc.New(cluster)
	// Create a new WS client (used for confirming transactions)
	wsClient, err := ws.Connect(ctx, wsCluster)
	if err != nil {
		log.Fatal("Failed at connection into websocket")
	}

	keyPairs, err := loadKeyPairs()
	if err != nil {
		log.Fatal("Failed at getting key pairs: ", err)
		return
	}

	programID, err := solana.PublicKeyFromBase58(ProgramID)
	if err != nil {
		log.Fatal("Failed at getting programID: ", err)
	}
	tic_tac_toe.SetProgramID(programID)

	tSetupGame := tic_tac_toe.NewSetupGameInstruction(
		keyPairs.PlayerTwoPrivateKey.PublicKey(),
		keyPairs.GamePrivateKey.PublicKey(),
		keyPairs.PlayerOnePrivateKey.PublicKey(),
		solana.SystemProgramID,
	)

	recent, err := client.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		log.Fatal("Failed at getting recent block Hash: ", err)
		return
	}
	tx, err := solana.NewTransaction(
		[]solana.Instruction{tSetupGame.Build()},
		recent.Value.Blockhash,
		solana.TransactionPayer(keyPairs.GamePrivateKey.PublicKey()),
	)
	if err != nil {
		log.Fatal("Failed at calling SetupGame transaction: ", err)
		return
	}
	log.Println("SetupGame tx:", tx)
	signers := []solana.PrivateKey{keyPairs.GamePrivateKey, keyPairs.PlayerOnePrivateKey}
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			for _, signer := range signers {
				if signer.PublicKey().Equals(key) {
					return &signer
				}
			}
			return nil
		})
	if err != nil {
		log.Fatal("unable to sign SetupGame transaction: ", err)
		return
	}

	err = tx.VerifySignatures()
	if err != nil {
		log.Fatal("Error at Verifying tx SetupGame signatures: ", err)
		return
	}
	// Send transaction, and wait for confirmation:
	sig, err := confirm.SendAndConfirmTransaction(
		ctx,
		client,
		wsClient,
		tx,
	)
	if err != nil {
		log.Fatal("Failed at sending and confirm SetupGame transaction: ", err)
		return
	}
	log.Println("SetupGame Sig: ", sig)
}

func loadKeyPairs() (KeyPairs, error) {
	gamePrivateKey, err := solana.PrivateKeyFromSolanaKeygenFile("./config/keypairs/game2_id.json")
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
