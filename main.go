package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ali-shokoohi/anchor-go-tic-tac-toe/pkg/generated/tic_tac_toe"
	bin "github.com/gagliardetto/binary"
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

	gameState, err := getGameState(ctx, client, keyPairs.GamePrivateKey.PublicKey())
	if err != nil {
		log.Fatal("Failed at getting game state:", err)
		return
	}
	log.Println("GameState:", gameState)

	// Check if game is already started
	if gameState.Turn == 0 {
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
			solana.TransactionPayer(keyPairs.PlayerOnePrivateKey.PublicKey()),
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
	err = play(ctx, client, wsClient, keyPairs, gameState)
	if err != nil {
		log.Fatal("Failed at calling play:", err)
		return
	}
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

func play(ctx context.Context, client *rpc.Client, wsClient *ws.Client, keyPairs KeyPairs, gameState *tic_tac_toe.Game) error {
	var player solana.PrivateKey
	if gameState.Turn%2 == 0 {
		player = keyPairs.PlayerTwoPrivateKey
	} else {
		player = keyPairs.PlayerOnePrivateKey
	}
	log.Println("\nPlayer:\n", player.String(), "\nTurn:", gameState.Turn)

	// Choosing the tile cell
	var cell uint
	fmt.Print("Enter a tile cell position: ")
	_, err := fmt.Scan(&cell)
	if err != nil {
		return err
	}

	if cell > 8 {
		err = errors.New("The cell is out of box! (0 : 8)")
		return err
	}

	// Setting the tile
	var tile *tic_tac_toe.Tile
	switch cell {
	case 0:
		tile = &tic_tac_toe.Tile{Row: 0, Column: 0}
	case 1:
		tile = &tic_tac_toe.Tile{Row: 0, Column: 1}
	case 2:
		tile = &tic_tac_toe.Tile{Row: 0, Column: 2}
	case 3:
		tile = &tic_tac_toe.Tile{Row: 1, Column: 0}
	case 4:
		tile = &tic_tac_toe.Tile{Row: 1, Column: 1}
	case 5:
		tile = &tic_tac_toe.Tile{Row: 1, Column: 2}
	case 6:
		tile = &tic_tac_toe.Tile{Row: 2, Column: 0}
	case 7:
		tile = &tic_tac_toe.Tile{Row: 2, Column: 1}
	case 8:
		tile = &tic_tac_toe.Tile{Row: 2, Column: 2}
	default:
		err = errors.New("Invalid cell. You should choose 0:8")
		return err
	}

	tPlay := tic_tac_toe.NewPlayInstruction(
		*tile,
		keyPairs.GamePrivateKey.PublicKey(),
		player.PublicKey(),
	)
	recent, err := client.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		err = errors.New("Failed at getting recent block Hash: " + err.Error())
		return err
	}
	tx, err := solana.NewTransaction(
		[]solana.Instruction{tPlay.Build()},
		recent.Value.Blockhash,
		solana.TransactionPayer(player.PublicKey()),
	)
	if err != nil {
		err = errors.New("Failed at calling SetupGame transaction: " + err.Error())
		return err
	}
	log.Println("Play tx:", tx)
	signers := []solana.PrivateKey{keyPairs.GamePrivateKey, player}
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
		log.Fatal("unable to sign play transaction: ", err)
		return nil
	}

	err = tx.VerifySignatures()
	if err != nil {
		log.Fatal("Error at Verifying tx play signatures: ", err)
		return nil
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
		return nil
	}
	log.Println("Play Sig: ", sig)

	gameState, err = getGameState(ctx, client, keyPairs.GamePrivateKey.PublicKey())
	if err != nil {
		log.Fatal("Failed at getting game state:", err)
		return err
	}
	log.Println("GameState:", gameState)

	return play(ctx, client, wsClient, keyPairs, gameState)
}

func getGameState(ctx context.Context, client *rpc.Client, gameAccountPublicKey solana.PublicKey) (*tic_tac_toe.Game, error) {
	// Fetch the account data corresponding to the game account's public key
	accountInfo, err := client.GetAccountInfo(ctx, gameAccountPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch game account info")
	}

	// Check if the account data exists
	if accountInfo.Value == nil {
		return nil, errors.New("game account data is nil")
	}

	// Convert DataBytesOrJSON to a byte slice
	dataBytes := accountInfo.Value.Data.GetBinary()

	// Define your Game struct
	var gameState tic_tac_toe.Game
	err = gameState.UnmarshalWithDecoder(bin.NewBinDecoder(dataBytes))
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode game state")
	}

	return &gameState, nil
}
