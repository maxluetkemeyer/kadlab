package cli

import (
	"bufio"
	"context"
	"d7024e_group04/internal/node"
	"fmt"
	"os"
	"strings"
)

func InputLoop(ctx context.Context, cancelCtx context.CancelFunc, node *node.Node) error {
	reader := bufio.NewReader(os.Stdin)
	errChan := make(chan error, 1)

	go func() {
		for {
			cliLogic(ctx, cancelCtx, errChan, reader, node)
		}
	}()

	for {
		select {
		case err := <-errChan:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func cliLogic(ctx context.Context, cancelCtx context.CancelFunc, errChan chan error, reader *bufio.Reader, node *node.Node) {
	fmt.Printf("$")

	input, err := reader.ReadString('\n')
	if err != nil {
		errChan <- err
	}

	if len(input) <= 0 {
		return
	}

	command := strings.Fields(strings.TrimSpace(input))
	switch command[0] {
	case "put":
		hash, err := node.PutObject(ctx, command[1])
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(hash)
		}
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
	case "forget":
		panic("TODO")
	default:
		fmt.Println("invalid command")
	}

}
