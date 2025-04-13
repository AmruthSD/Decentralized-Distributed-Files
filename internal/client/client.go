package client

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/AmruthSD/Decentralized-Distributed-Files/internal/config"
)

func UploadFile(filePath string) ([]string, error) {
	// open file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

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
	outputFile, err := os.Create(filePath + ".hash")
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	for i := 0; i < len(hashes); i++ {
		outputFile.WriteString(hashes[i] + "\n")
	}
	outputFile.Close()

	// for each hash get nodes

	// send file chunk

	return hashes, nil
}
