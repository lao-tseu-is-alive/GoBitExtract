package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// Define encryption type constants
const (
	XTS_AES_128 = 0x0480
	XTS_AES_256 = 0x0580
	AES_CBC_128 = 0x0280
	AES_CBC_256 = 0x0380
)

// findAllOccurrences returns all indices where pattern is found in data
func findAllOccurrences(data, pattern []byte) []int {
	var indices []int
	for i := 0; i <= len(data)-len(pattern); i++ {
		if bytes.Equal(data[i:i+len(pattern)], pattern) {
			indices = append(indices, i)
			log.Printf("Found FVE metadata at offset  \t(0x%x) : \t[%x] [%x] %v", i, data[i:i+4], data[i+4:i+8], data[i:i+8])
		}
	}
	return indices
}

type BitLockerKey struct {
	Offset         int64
	KeyData        []byte
	EncryptionType uint16
}

func main() {
	// Parse command line arguments
	inputFile := flag.String("input", "", "Memory dump file to analyze")
	outputFile := flag.String("output", "bitlocker_keys.bin", "Output file for extracted keys")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Please specify input memory dump file using -input")
		os.Exit(1)
	}

	// Open the memory dump file
	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Create output file for keys
	outFile, err := os.Create(*outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	// Define the patterns we're looking for
	fvePattern := []byte("-FVE-FS-")            // Main BitLocker identifier
	keyPrefix := []byte{0x03, 0x20, 0x01, 0x00} // Key start prefix

	// Possible encryption type prefixes (2 bytes each)
	encTypes := map[uint16]string{
		XTS_AES_128: "XTS-AES-128",
		XTS_AES_256: "XTS-AES-256",
		AES_CBC_128: "AES-CBC-128",
		AES_CBC_256: "AES-CBC-256",
	}

	// Buffer for reading the file (16KB)
	buffer := make([]byte, 16384)
	var offset int64 = 0
	var keysFound []BitLockerKey

	fmt.Println("Starting BitLocker key search...")
	startTime := time.Now()

	// Read file in chunks
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}
		if n == 0 {
			break
		}

		// Search for FVE pattern in current buffer
		fveIndices := findAllOccurrences(buffer[:n], fvePattern)

		for _, fveIndex := range fveIndices {
			// Calculate absolute offset in file
			absOffset := offset + int64(fveIndex)

			// Look for key prefix after FVE pattern
			// Typically within 100 bytes after FVE pattern
			searchWindow := 100
			if fveIndex+searchWindow >= n {
				// Handle buffer boundary case
				continue
			}

			remainingBuffer := buffer[fveIndex:n]
			keyStart := bytes.Index(remainingBuffer, keyPrefix)

			if keyStart >= 0 && keyStart < searchWindow {
				// Found potential key
				keyOffset := absOffset + int64(keyStart) + int64(len(fvePattern))

				// Check for encryption type prefix (2 bytes before key prefix)
				encTypeOffset := keyStart - 2
				if encTypeOffset >= 0 {
					encTypeBytes := remainingBuffer[encTypeOffset : encTypeOffset+2]
					encType := binary.BigEndian.Uint16(encTypeBytes)

					// Verify if it's a valid encryption type
					if encTypeName, exists := encTypes[encType]; exists {
						// Extract key (typically 32 bytes for AES-256, 16 for AES-128)
						keyLength := 32 // Default to maximum possible length
						if encType == AES_CBC_128 || encType == XTS_AES_128 {
							keyLength = 16
						}

						if keyStart+len(keyPrefix)+keyLength < len(remainingBuffer) {
							keyData := remainingBuffer[keyStart+len(keyPrefix) : keyStart+len(keyPrefix)+keyLength]

							keysFound = append(keysFound, BitLockerKey{
								Offset:         keyOffset,
								KeyData:        keyData,
								EncryptionType: encType,
							})

							fmt.Printf("Found key at offset 0x%x, Type: %s\n",
								keyOffset, encTypeName)
						}
					}
				}
			}
		}

		offset += int64(n)
	}

	// Write found keys to output file
	for i, key := range keysFound {
		// Write key metadata
		binary.Write(outFile, binary.BigEndian, key.Offset)
		binary.Write(outFile, binary.BigEndian, key.EncryptionType)
		binary.Write(outFile, binary.BigEndian, uint32(len(key.KeyData)))
		outFile.Write(key.KeyData)

		fmt.Printf("Wrote key %d to output file\n", i+1)
	}

	duration := time.Since(startTime)
	fmt.Printf("\nSearch completed in %v\n", duration)
	fmt.Printf("Total keys found: %d\n", len(keysFound))
	fmt.Printf("Keys saved to: %s\n", *outputFile)
}
