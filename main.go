package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"os"

	"ethereum/models"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	dbHost     = "your-db-host"
	dbPort     = "your-db-port"
	dbUser     = "your-db-user"
	dbPassword = "your-db-password"
	dbName     = "your-db-name"
)
const numWorkers = 10

func worker(ctx context.Context, client *ethclient.Client, startBlock uint64, endBlock uint64, db *sql.DB, resultCh chan<- error) {
	for blockNum := startBlock; blockNum <= endBlock; blockNum++ {
		block, err := client.BlockByNumber(ctx, big.NewInt(int64(blockNum)))
		if err != nil {
			resultCh <- err
			return
		}

		select {
		case <-ctx.Done():
			resultCh <- nil
			return
		default:
			SaveBlockDataToDB(block, db)
		}
	}

	resultCh <- nil
}

func main() {
	ctx := context.Background()
	url := os.Getenv("ethereum_node_url")
	// Create a new Ethereum client
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Connect to the PostgreSQL database
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Get the latest block number
	latestBlock, err := client.BlockNumber(ctx)
	if err != nil {
		log.Fatalf("Failed to get the latest block number: %v", err)
	}

	// Create a channel to receive worker results
	resultCh := make(chan error)

	// Calculate the number of blocks each worker should process
	blocksPerWorker := (latestBlock + 1) / numWorkers

	// Start the workers in parallel
	for i := 0; i < numWorkers; i++ {
		startBlock := uint64(i) * blocksPerWorker
		endBlock := startBlock + blocksPerWorker - 1

		// The last worker handles any remaining blocks
		if i == numWorkers-1 {
			endBlock = latestBlock
		}

		go worker(ctx, client, startBlock, endBlock, db, resultCh)
	}

	// Wait for all workers to finish
	for i := 0; i < numWorkers; i++ {
		if err := <-resultCh; err != nil {
			log.Printf("Worker error: %v", err)
		}
	}

	log.Println("Indexing completed.")

	g := gin.New()
	g.Use(gin.Logger())
	Router(g)
	_ = g.Run(os.Getenv("address"))
}

func SaveBlockDataToDB(block *types.Block, db *sql.DB) {

}

func Router(g *gin.Engine) {
	{
		g.GET("/blocks", models.LatestBlocks)
		g.GET("/blocks/:id", models.GetBlockByID)
		g.GET("/transaction/:txHash", models.GetTrans)
	}
}
