package env

import (
	"os"
	"strconv"
)

var Port = 50051
var IDLength = 20
var BucketSize = 20

func init() {
	port := os.Getenv("PORT")
	idLength := os.Getenv("ID_LENGTH")
	bucketSize := os.Getenv("BUCKET_SIZE")

	if port != "" {
		portInt, err := strconv.Atoi(port)
		if err == nil {
			Port = portInt
		}
	}

	if idLength != "" {
		idLengthInt, err := strconv.Atoi(idLength)
		if err == nil {
			IDLength = idLengthInt
		}
	}

	if bucketSize != "" {
		bucketSizeInt, err := strconv.Atoi(bucketSize)
		if err == nil {
			BucketSize = bucketSizeInt
		}
	}

}
