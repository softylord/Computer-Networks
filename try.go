package main

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/ethereum/go-ethereum/ethclient"
	"google.golang.org/api/option"
	"log"
	"math/big"
	"time"
	"strconv"
)

func main() {
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: "https://lab11-e48cc-default-rtdb.firebaseio.com",
	}
	// Fetch the service account key JSON file contents
	opt := option.WithCredentialsFile("C://Users//katy2//Downloads//lab11-e48cc-firebase-adminsdk-waycq-351da366b9.json")

	// Initialize the app with a service account, granting admin privileges
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("Error initializing app:", err)
	}

	// Create a database client from App.
	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}

	// Get a database reference to our blog.
	ref := client.NewRef("blockchain")
	// User is a json-serializable type.
	type Block struct {
		Number     uint64 `json:"number, omitempty"`
		Time       uint64 `json:"time, omitempty"`
		Difficulty uint64 `json:"difficulty,omitempty"`
		Hash       string `json:"hash,omitempty"`
		TransLen   int    `json:"transations_len,omitempty"`
	}
	type Tranact struct {
		ChainId  *big.Int `json:"chain id, omitempty"`
		Hash     string   `json:"hash, omitempty"`
		Value    *big.Int `json:"value, omitempty"`
		Cost     *big.Int `json:"cost, omitempty"`
		To       string   `json:"to,omitempty"`
		Gas      uint64   `json:"gas,omitempty"`
		GasPrice *big.Int `json:"gas price,omitempty"`
	}

	usersRef := ref.Child("blocks")

	cl, err := ethclient.Dial("https://mainnet.infura.io/v3/8133ff0c11dc491daac3f680d2f74d18")
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		for {

			header, err := cl.HeaderByNumber(context.Background(), nil)
			if err != nil {
				log.Fatal(err)
			}
			blockNumber := big.NewInt(header.Number.Int64())
			block, err := cl.BlockByNumber(context.Background(), blockNumber) //get block with this number
			if err != nil {
				log.Fatal(err)
			}

			if err := usersRef.Child("last block").Set(ctx, &Block{
				Number:     block.Number().Uint64(),
				Time:       block.Time(),
				Difficulty: block.Difficulty().Uint64(),
				Hash:       block.Hash().Hex(),
				TransLen:   len(block.Transactions()),
			}); err != nil {
				log.Fatalln("Error setting value:", err)
			}
		}
	}()
	go func() {
		blockNumber := big.NewInt(1236541)
		block, err := cl.BlockByNumber(context.Background(), blockNumber) //get block with this number
		if err != nil {
			log.Fatal(err)
		}

		if err := usersRef.Child("block number 1236541").Set(ctx, &Block{
			Number:     block.Number().Uint64(),
			Time:       block.Time(),
			Difficulty: block.Difficulty().Uint64(),
			Hash:       block.Hash().Hex(),
			TransLen:   len(block.Transactions()),
		}); err != nil {
			log.Fatalln("Error setting value:", err)
		}

		blockNumber = big.NewInt(15960495)
		block, err = cl.BlockByNumber(context.Background(), blockNumber) //get block with this number
		if err != nil {
			log.Fatal(err)
			
		}
		c:=0
		for _, tx := range block.Transactions() {
			c++
			if err := usersRef.Child("block number 15960495 Tranactions "+strconv.Itoa(c) ).Set(ctx, &Tranact{
				ChainId:  tx.ChainId(),
				Hash:     tx.Hash().String(),
				Value:    tx.Value(),
				Cost:     tx.Cost(),
				To:       tx.To().String(),
				Gas:      tx.Gas(),
				GasPrice: tx.GasPrice(),
			}); err != nil {
				log.Fatalln("Error setting value:", err)
			}
		}
	}()
	<-time.After(time.Minute * 120)
}
