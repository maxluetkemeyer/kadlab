package cli

import (
	"bufio"
	"context"
	"d7024e_group04/internal/node"
	"fmt"
	"os"
	"strings"
)

func InputLoop(ctx context.Context, cancelCtx context.CancelFunc, node node.NodeHandler) error {
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

func cliLogic(ctx context.Context, cancelCtx context.CancelFunc, errChan chan error, reader *bufio.Reader, node node.NodeHandler) {
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
	case "me":
		fmt.Println(node.Me())

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

		str, err := getCommand(ctx, node, hash)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(str)
		}

	case "exit":
		cancelCtx()
	case "forget":
		fmt.Println("Not implemented yet")
	default:
		fmt.Println("invalid command")
	}

}

func getCommand(ctx context.Context, node node.NodeHandler, hash string) (string, error) {
	val, candidates, err := node.GetObject(ctx, hash)

	if err != nil {
		return "", fmt.Errorf("failed to get value and candidates, err: %v", err)
	}

	if val != nil {
		return fmt.Sprintf("value: %v, found in node: %v", val.DataValue, val.NodeWithValue), nil
	}

	if candidates != nil {
		str := "could not find value, closest contacts are:"
		for _, contact := range candidates {
			str += fmt.Sprintln(contact)
		}
		return str, nil
	}

	return "", fmt.Errorf("wtf")
}
