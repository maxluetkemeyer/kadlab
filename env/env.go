package env

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"
)

const IDLength = 20 // The static number of bytes in a KademliaID. 160 / 8 = 20

var ApiPort = 8080
var Port = 50051
var BucketSize = 20 // must be larger than 0
var NodesProxyDomain = "kademlianodes"
var Alpha = 3 // degree of parallelism
var RPCTimeout = 5 * time.Second
var BootstrapTimeout = 1 * time.Minute
var TRefresh = 3600 * time.Second
var TTL = 60 * 60 * time.Second

func init() {
	log.Println("Initialize environment variables")

	port := os.Getenv("PORT")
	bucketSize := os.Getenv("BUCKET_SIZE")
	apiPort := os.Getenv("API_PORT")
	nodesProxyDomain := os.Getenv("NODES_PROXY_DOMAIN")
	alpha := os.Getenv("ALPHA")
	rpcTimeoutInSeconds := os.Getenv("RPC_TIMEOUT_IN_SECONDS")
	ttlInSeconds := os.Getenv("TTLinSeconds")
	_, debug := os.LookupEnv("DEBUG")

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

		if bucketSizeInt <= 0 {
			log.Fatalf("invalid bucket size: %v", bucketSize)
		}

	}

	if nodesProxyDomain != "" {
		NodesProxyDomain = nodesProxyDomain
	}

	if alpha != "" {
		alphaInt, err := strconv.Atoi(alpha)
		if err == nil {
			Alpha = alphaInt
		}
	}

	if rpcTimeoutInSeconds != "" {
		rpcTimeoutInt, err := strconv.Atoi(rpcTimeoutInSeconds)
		if err == nil {
			RPCTimeout = time.Duration(rpcTimeoutInt) * time.Second
		}
	}

	if ttlInSeconds != "" {
		ttlInSecondsInt, err := strconv.Atoi(ttlInSeconds)
		if err != nil {
			log.Fatalf("invalid ttl size: %v", ttlInSeconds)
		}
		TTL = time.Duration(ttlInSecondsInt) * time.Second
		fmt.Printf("TTL %v", TTL)
	}

	if debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
}
