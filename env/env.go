package env

import (
	"os"
	"strconv"
)

var Port = 50051

const IDLength = 20
const BucketSize = 20

func init() {
	port := os.Getenv("PORT")

	if port != "" {
		p, err := strconv.Atoi(port)
		if err == nil {
			Port = p
		}
	}

}
