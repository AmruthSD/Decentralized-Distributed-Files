package client

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/config"
)

func HashFile(filePath string) ([]string, error) {
	// open file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	filename := filepath.Base(filePath)
	// make hashes
	hashes := make([]string, 0)
	buffer := make([]byte, config.MetaData.ChunkSize)
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("read error: %w", err)
		}

		if n == 0 {
			break
		}

		hash := sha256.Sum256(buffer[:n])
		hashes = append(hashes, hex.EncodeToString(hash[:]))
	}

	file.Close()

	// write hashes list to another file
	dirPath := "./files/" + strconv.Itoa(int(config.MetaData.Port)) + "/"
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	outputFile, err := os.Create(dirPath + filename + ".hash")
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	for i := 0; i < len(hashes); i++ {
		outputFile.WriteString(hashes[i] + "\n")
	}
	outputFile.Close()
	return hashes, nil
}
