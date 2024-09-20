package env

import (
	"os"
	"strconv"
	"time"
)

const IDLength = 20
const Alpha = 3 // degree of parallelism
const RPCTimeout = 5 * time.Second

var ApiPort = 8080
var Port = 50051
var BucketSize = 20
var KnownDomain = "kademlianodes"

func init() {
	port := os.Getenv("PORT")
	bucketSize := os.Getenv("BUCKET_SIZE")
	apiPort := os.Getenv("API_PORT")

	if port != "" {
		portInt, err := strconv.Atoi(port)
		if err == nil {
			Port = portInt
		}
	}

	if apiPort != "" {
		apiPortInt, err := strconv.Atoi(apiPort)
		if err == nil {
			ApiPort = apiPortInt
		}
	}

	if bucketSize != "" {
		bucketSizeInt, err := strconv.Atoi(bucketSize)
		if err == nil {
			BucketSize = bucketSizeInt
		}
	}

}
