package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

const (
	defaultMemoryDumpPath = "memory_dump.mem"
	outputPath            = "extracted_fvek.bin"
)

// readMemoryDump reads the memory dump file at the given path
func readMemoryDump(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read memory dump: %v", err)
	}
	return data, nil
}

type FVEMetadata struct {
	Signature [8]byte
	Size      uint32
	Version   uint32
}

// searchFVEMetadata looks for BitLocker FVE structures
func searchFVEMetadata(dump []byte) ([]FVEMetadata, error) {
	var results []FVEMetadata
	signature := []byte{0x2D, 0x46, 0x56, 0x45, 0x2D, 0x46, 0x53, 0x2D} // "-FVE-FS-"

	for i := 0; i < len(dump)-len(signature); i++ {
		if bytesEqual(dump[i:i+len(signature)], signature) {
			log.Printf("Found FVE metadata at offset %d, \thex:%x", i, i)
			metadata := FVEMetadata{}
			copy(metadata.Signature[:], signature)
			metadata.Size = binary.LittleEndian.Uint32(dump[i+8 : i+12])
			metadata.Version = binary.LittleEndian.Uint32(dump[i+12 : i+16])
			results = append(results, metadata)
		}
	}
	return results, nil
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type FVEKeyData struct {
	Offset uint64
	Data   []byte
}

// extractPotentialFVEK extracts potential FVEKs from identified structures
func extractPotentialFVEK(dump []byte, metadata []FVEMetadata) ([]FVEKeyData, error) {
	var potentialKeys []FVEKeyData

	for _, m := range metadata {
		keyOffset := uint64(m.Size) + 0x100 // Example offset, adjust as needed
		if keyOffset < uint64(len(dump)) {
			keyData := dump[keyOffset : keyOffset+32] // Assuming 256-bit key
			if isValidFVEK(keyData) {
				potentialKeys = append(potentialKeys, FVEKeyData{
					Offset: keyOffset,
					Data:   keyData,
				})
			}
		}
	}
	return potentialKeys, nil
}

// validateAndSaveFVEK validates and saves extracted FVEKs
func validateAndSaveFVEK(keys []FVEKeyData) error {
	validKeyCount := 0
	for _, key := range keys {
		if isValidFVEK(key.Data) {
			filename := fmt.Sprintf("%s_%d", outputPath, validKeyCount)
			err := saveFVEK(key.Data, filename)
			if err != nil {
				return fmt.Errorf("failed to save FVEK: %v", err)
			}
			log.Printf("Potential FVEK saved: %s, containing %x", filename, key.Data)
			validKeyCount++
		}
	}
	if validKeyCount == 0 {
		log.Printf("No valid FVEK found.")
	} else {
		log.Printf("%d valid FVEKs extracted.", validKeyCount)
	}
	return nil
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
	if zeroCount > 25 {
		return false
	}

	return true
}

// saveFVEK writes a valid FVEK to a file
func saveFVEK(data []byte, path string) error {
	return os.WriteFile(path, data, 0600)
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
	sizeOfDump := len(dump)
	log.Printf("Memory dump read successfully. Size: %d bytes", sizeOfDump)

	log.Printf("Searching for FVE metadata...")
	metadata, err := searchFVEMetadata(dump)
	if err != nil {
		log.Fatalf("Failed to search for FVE metadata: %v", err)
	}
	log.Printf("Found %d FVE metadata structures", len(metadata))

	log.Printf("Extracting potential FVEK...")
	keys, err := extractPotentialFVEK(dump, metadata)
	if err != nil {
		log.Fatalf("Failed to extract potential FVEK: %v", err)
	}
	log.Printf("Extracted %d potential FVEK", len(keys))

	for _, key := range keys {
		log.Printf("Potential FVEK offset: %d, \thex:%x", key.Offset, key.Offset)
		log.Printf("Potential FVEK data: %x", key.Data)
		log.Printf("Potential FVEK data length: %d", len(key.Data))
	}

	log.Printf("Validating and saving potential FVEK...")
	err = validateAndSaveFVEK(keys)
	if err != nil {
		log.Fatalf("Failed to validate and save FVEK: %v", err)
	}

	log.Printf("Potential FVEK extraction completed successfully.")
	log.Printf("You can try to decrypt the disk using the extracted FVEK with a tool like dislocker on Linux:")
	log.Printf("sudo dislocker -V /dev/sdX -k extracted_fvek.bin_0 --dislocker-file your_dislocker.img")
}
