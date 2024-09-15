package env

import (
	"os"
	"strconv"
)

const IDLength = 20

var Port = 50051
var BucketSize = 20

func init() {
	port := os.Getenv("PORT")
	bucketSize := os.Getenv("BUCKET_SIZE")

	if port != "" {
		portInt, err := strconv.Atoi(port)
		if err == nil {
			Port = portInt
		}
	}

	if bucketSize != "" {
		bucketSizeInt, err := strconv.Atoi(bucketSize)
		if err == nil {
			BucketSize = bucketSizeInt
		}
	}

}
