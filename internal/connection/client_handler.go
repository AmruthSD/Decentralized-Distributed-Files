package connection

import (
	"bufio"
	"fmt"
	"os"

	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/client"
)

func (node *Node) Handle_Client() {
	fmt.Println("Client Started")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		if scanner.Scan() {
			fmt.Println("Your file path:", scanner.Text())
			hashes, err := client.HashFile(scanner.Text())
			if err != nil {
				fmt.Println(err)
				continue
			}
			for i := 0; i < len(hashes); i++ {

			}
		}
	}
}
