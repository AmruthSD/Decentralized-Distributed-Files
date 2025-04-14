package connection

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/client"
	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/config"
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

			filePath := scanner.Text()
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Println("failed to open file: %w", err)
				continue
			}
			f := 0
			buffer := make([]byte, config.MetaData.ChunkSize)
			for i := 0; i < len(hashes); i++ {
				n, err := file.Read(buffer)
				if err != nil && err != io.EOF {
					fmt.Println("read error: %w", err)
					f = 1
					break
				}
				if n == 0 {
					break
				}
				key, _ := hex.DecodeString(hashes[i])
				// for each hash get nodes
				nodes := node.get_closest_nodes(key)

				var wg sync.WaitGroup
				for j := 0; j < len(nodes); j++ {
					wg.Add(1)
					fmt.Println(nodes[j].Address)
					go func(hash string) {
						defer wg.Done()
						// send file chunk
						node.send_chunk(buffer[:n], hash, nodes[j].Address)
					}(hashes[i])
				}
				wg.Wait()

				fmt.Println("Done with hash number", i)
			}
			file.Close()
			if f == 1 {
				continue
			}
			fmt.Println("File sent", filePath)
		}
	}
}

func (node *Node) send_chunk(buffer []byte, hash string, peer_address string) {
	conn, err := net.Dial("tcp", peer_address)
	if err != nil {
		fmt.Println("Connection err at ", peer_address, err)
		return
	}
	readbuff := make([]byte, config.MetaData.ChunkSize)
	conn.Write([]byte(fmt.Sprintf("STORE %s\n", hash)))
	conn.Read(readbuff)
	conn.Write(buffer)
}

func (node *Node) handel_store(parts []string, conn net.Conn) string {
	hash := parts[1]
	dirPath := "./files/" + strconv.Itoa(int(config.MetaData.Port)) + "/storage/"
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	outputFile, _ := os.Create(dirPath + hash + ".hash")
	buffer := make([]byte, config.MetaData.ChunkSize)
	conn.Write([]byte("OKAY SEND"))
	n, _ := conn.Read(buffer)
	outputFile.Write(buffer[:n])
	outputFile.Close()
	return "DONE"
}
