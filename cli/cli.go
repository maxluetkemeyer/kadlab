package cli

import (
	"bufio"
	"context"
	"d7024e_group04/internal/node"
	"fmt"
	"io"
	"os"
	"strings"
)

func InputLoop(ctx context.Context, cancelCtx context.CancelFunc, stdin io.Reader, node *node.Node) error {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("$")

		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		if len(input) > 0 {
			command := strings.Fields(strings.TrimSpace(input))
			switch command[0] {
			case "put":
				node.PutObject()
				panic("TODO")
			case "get":
				hash := command[1]
				if hash == "" {
					fmt.Println("no hash was provided")
					break
				}
				node.GetObject(ctx, hash)
				panic("TODO")
			case "exit":
				cancelCtx()
				return ctx.Err()
			case "forget":
				panic("TODO")
			default:
				fmt.Println("invalid command")
			}
		}

	}

}
