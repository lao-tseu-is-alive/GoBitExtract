package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

const (
	defaultMemoryDumpPath = "memory_dump.mem"
	outputPath            = "extracted_fvek.bin"
)

type FVEKeyData struct {
	Offset uint64
	Data   []byte
}

// FVE metadata search constants
var (
	fveSignature  = []byte("-FVE-FS-")             // Signature marking the start of the FVE structure
	vmkStartBytes = []byte{0x03, 0x20, 0x01, 0x00} // Unique pattern preceding VMK
)

// readMemoryDump reads the memory dump file at the given path
func readMemoryDump(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read memory dump: %v", err)
	}
	return data, nil
}

// searchFVEK scans memory for FVE structures and extracts the FVEK
func searchFVEK(dump []byte) []FVEKeyData {
	var potentialKeys []FVEKeyData

	offset := 0
	for {
		// Look for the FVE signature
		index := bytes.Index(dump[offset:], fveSignature)
		if index == -1 {
			break
		}
		index += offset // Adjust index relative to full dump
		log.Printf("Found FVE metadata at offset %d (0x%x)", index, index)

		// Check if version field at offset 4 is 1
		versionOffset := index + 8 + 4
		if versionOffset+4 > len(dump) {
			offset = index + len(fveSignature)
			continue
		}
		version := binary.LittleEndian.Uint32(dump[versionOffset : versionOffset+4])
		if version != 1 {
			log.Printf("Skipping structure at 0x%x: version mismatch (%d)", index, version)
			offset = index + len(fveSignature)
			continue
		}

		// Search for the known bytes `\x03\x20\x01\x00` within this structure
		vfkStartIndex := bytes.Index(dump[index:], vmkStartBytes)
		if vfkStartIndex == -1 {
			log.Printf("Skipping structure at 0x%x: VMK pattern not found", index)
			offset = index + len(fveSignature)
			continue
		}
		vfkStartIndex += index + len(vmkStartBytes) // Adjust position

		// Ensure we are within bounds before extracting the key
		if vfkStartIndex+32 > len(dump) {
			log.Printf("Skipping structure at 0x%x: FVEK location out of bounds", index)
			offset = index + len(fveSignature)
			continue
		}

		// Extract the potential FVEK
		keyData := dump[vfkStartIndex : vfkStartIndex+32]
		if isValidFVEK(keyData) {
			potentialKeys = append(potentialKeys, FVEKeyData{
				Offset: uint64(vfkStartIndex),
				Data:   keyData,
			})
			log.Printf("Potential FVEK found at offset 0x%x: %x", vfkStartIndex, keyData)
		} else {
			log.Printf("Discarding invalid FVEK at offset 0x%x", vfkStartIndex)
		}

		// Move to the next occurrence
		offset = index + len(fveSignature)
	}

	return potentialKeys
}

// isValidFVEK checks if the extracted key is likely a valid FVEK
func isValidFVEK(data []byte) bool {
	if len(data) != 32 {
		return false
	}
	zeroCount := 0
	for _, b := range data {
		if b == 0x00 {
			zeroCount++
		}
	}
	// Reject keys with more than 80% zeroes
	return zeroCount <= 25
}

// saveFVEK writes a valid FVEK to a file
func saveFVEK(data []byte, filename string) error {
	return os.WriteFile(filename, data, 0600)
}

func main() {
	memoryDumpPath := defaultMemoryDumpPath
	if len(os.Args) > 1 {
		memoryDumpPath = os.Args[1]
	}

	log.Printf("Reading memory dump from: %s", memoryDumpPath)
	dump, err := readMemoryDump(memoryDumpPath)
	if err != nil {
		log.Fatalf("Failed to read memory dump: %v", err)
	}
	log.Printf("Memory dump read successfully. Size: %d bytes", len(dump))

	log.Printf("Searching for FVEK...")
	keys := searchFVEK(dump)
	if len(keys) == 0 {
		log.Printf("No valid FVEK found.")
		return
	}

	log.Printf("Extracted %d valid FVEKs", len(keys))
	for i, key := range keys {
		filename := fmt.Sprintf("%s_%d", outputPath, i)
		err := saveFVEK(key.Data, filename)
		if err != nil {
			log.Fatalf("Failed to save FVEK: %v", err)
		}
		log.Printf("Saved valid FVEK to: %s", filename)
	}

	log.Printf("Potential FVEK extraction completed successfully.")
	log.Printf("You can try to decrypt the disk using the extracted FVEK with a tool like dislocker on Linux:")
	log.Printf("sudo dislocker -V /dev/sdX -k extracted_fvek.bin_0 --dislocker-file your_dislocker.img")
}
