package connection

import (
	"bufio"
	"compress/gzip"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/client"
	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/config"
)

// reading and executing user commands
func (node *Node) Handle_Client() {
	fmt.Println("Client Started")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if scanner.Scan() {
			parts := strings.Split(scanner.Text(), " ")
			if len(parts) != 2 {
				fmt.Println("INVALID INPUT")
				continue
			}
			filePath := parts[1]
			action := parts[0]
			if action == "UPLOAD" {
				node.UploadFile(filePath)
			} else if action == "DOWNLOAD" {
				node.DownLoadFile(filePath)
			} else if action == "DELETE" {
				node.DeleteFile(filePath)
			} else {
				fmt.Println("INVALID INPUT")
				continue
			}

		}
	}
}

// actually makes chunks and sends them to the neighbours
func (node *Node) UploadFile(filePath string) {
	fmt.Println("Your file path:", filePath)
	hashes, err := client.HashFile(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	compressedDirPath := "./files/" + strconv.Itoa(int(config.MetaData.Port)) + "/compressed/"
	file, err := os.Open(compressedDirPath + filePath)
	if err != nil {
		fmt.Println("failed to open file: %w", err)
		return
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
			fmt.Println("sending to", hex.EncodeToString(nodes[j].Node_id), nodes[j].Address)
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
		return
	}
	fmt.Println("File sent", filePath)
	os.Remove(compressedDirPath + filePath)
}

// makes connection and sends the chunk to store
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

// stores the file in directory
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
	UpdateTimeStamp(hash)
	return "DONE"
}

// downloads the file after querying each chunk
func (node *Node) DownLoadFile(fileName string) {
	dir := "./files/" + strconv.Itoa(int(config.MetaData.Port)) + "/hashed/"

	file, err := os.Open(dir + fileName + ".hash")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	scanner := bufio.NewScanner(file)
	compressedDirPath := "./files/" + strconv.Itoa(int(config.MetaData.Port)) + "/compressed/"
	outdir := "./files/" + strconv.Itoa(int(config.MetaData.Port)) + "/downloaded/"
	os.MkdirAll(outdir, os.ModePerm)
	outputFile, _ := os.Create(outdir + fileName)
	compressedFile, err := os.Create(compressedDirPath + fileName)
	if err != nil {
		log.Fatal(err)
	}

	for scanner.Scan() {
		hashVal := scanner.Text()
		fmt.Println(hashVal)
		cid, _ := hex.DecodeString(hashVal)
		nodes := node.get_closest_nodes(cid)
		written := 0
		for _, v := range nodes {
			conn, err := net.Dial("tcp", v.Address)
			if err != nil {
				continue
			}

			conn.Write([]byte(fmt.Sprintf("DOYOUHAVE %s DOWNLOAD\n", hashVal)))
			fmt.Println("Sent do you have download")
			readbuff := make([]byte, config.MetaData.ChunkSize)
			n, _ := conn.Read(readbuff)
			if string(readbuff[:n]) == "YES\n" {
				n, _ = conn.Read(readbuff)
				compressedFile.Write(readbuff[:n])
				written = 1
				break
			}
		}
		if written == 0 {
			fmt.Println("ERROR RIP")
			return
		}
	}

	compressedFile.Seek(0, 0)

	gzipReader, err := gzip.NewReader(compressedFile)
	if err != nil {
		return
	}
	defer gzipReader.Close()
	_, _ = io.Copy(outputFile, gzipReader)
	os.Remove(compressedDirPath + fileName)

	file.Close()
	outputFile.Close()
	fmt.Println("DOWNLOAD COMPLETE")
}

/*
if its a download type then we should actually send it else chunk
if check just says yes or no
*/
func (node *Node) handle_doyouhave(parts []string, conn net.Conn) string {
	id := parts[1]
	dir := "./files/" + strconv.Itoa(int(config.MetaData.Port)) + "/storage/"
	kind := parts[2]
	if kind == "DOWNLOAD" {
		buffer := make([]byte, config.MetaData.ChunkSize)
		file, err := os.Open(dir + id + ".hash")
		if err != nil {
			conn.Write([]byte("NO\n"))
			return "STOP"
		}
		conn.Write([]byte("YES\n"))
		n, _ := file.Read(buffer)
		conn.Write(buffer[:n])
		return "STOP"
	} else if kind == "CHECK" {
		buffer := make([]byte, config.MetaData.ChunkSize)
		_, err := os.Open(dir + id)
		if err == nil {
			conn.Write([]byte("YES\n"))
			return "STOP"
		}
		conn.Write([]byte("NO\n"))
		n, _ := conn.Read(buffer)
		outputFile, _ := os.Create(dir + id)
		outputFile.Write(buffer[:n])
		UpdateTimeStamp(id)
	}
	return "STOP"
}

// deletes file from the directory i.e the hashes
func (node *Node) DeleteFile(fileName string) {
	dir := "./files/" + strconv.Itoa(int(config.MetaData.Port)) + "/hashed/"
	os.Remove(dir + fileName + ".hash")
	fmt.Println("Deleted File")
}
