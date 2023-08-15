package util

import (
	"go-gin-restful-service/log"
	"time"

	"github.com/bwmarrin/snowflake"
)

func GenerateSnowID() int64 {
	// Create a new Node with a Node number of 1
	time.Sleep(1 * time.Nanosecond)
	node, err := snowflake.NewNode(1000)
	if err != nil {
		log.Logger.Error(err)
		return 0
	}
	// Generate a snowflake ID.
	id := node.Generate()
	// Print out the ID in a few different ways.
	return id.Int64()
}
