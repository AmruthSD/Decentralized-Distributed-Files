package connection

import (
	"bufio"
	"fmt"
	"os"
)

func (node *Node) Handle_Client() {

	scanner := bufio.NewScanner(os.Stdin)
	for {
		if scanner.Scan() {
			fmt.Println("Your file path:", scanner.Text())

		}
	}
}
