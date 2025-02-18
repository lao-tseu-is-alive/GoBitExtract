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

// readMemoryDump Read the memory dump file at given path
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
	// Add more fields as needed
}

// searchFVEMetadata  search for BitLocker metadata, specifically looking for FVE (Full Volume Encryption) structures.
func searchFVEMetadata(dump []byte) ([]FVEMetadata, error) {
	var results []FVEMetadata
	signature := []byte{0x2D, 0x46, 0x56, 0x45, 0x2D, 0x46, 0x53, 0x2D} // "-FVE-FS-"

	for i := 0; i < len(dump)-len(signature); i++ {
		if bytesEqual(dump[i:i+len(signature)], signature) {
			log.Printf("Found FVE metadata at offset %d, hex:%x", i, i)
			metadata := FVEMetadata{}
			metadata.Signature = [8]byte(signature)
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

// extractPotentialFVEK will try to Extract potential FVEK from identified structures
func extractPotentialFVEK(dump []byte, metadata []FVEMetadata) ([]FVEKeyData, error) {
	var potentialKeys []FVEKeyData

	for _, m := range metadata {
		// This is a simplified example; actual FVEK extraction is more complex
		// and may require analyzing multiple related structures
		keyOffset := uint64(m.Size) + 0x100 // Example offset, adjust as needed
		if keyOffset < uint64(len(dump)) {
			keyData := dump[keyOffset : keyOffset+32] // Assuming 256-bit key
			potentialKeys = append(potentialKeys, FVEKeyData{
				Offset: keyOffset,
				Data:   keyData,
			})
		}
	}
	return potentialKeys, nil
}

// validateAndSaveFVEK validate the extracted data to ensure it's likely to be the FVEK.
func validateAndSaveFVEK(keys []FVEKeyData) error {
	for i, key := range keys {
		// Implement validation logic here
		// This could involve checking key properties, entropy, etc.
		if isValidFVEK(key.Data) {
			err := saveFVEK(key.Data, fmt.Sprintf("%s_%d", outputPath, i))
			if err != nil {
				return fmt.Errorf("failed to save FVEK: %v", err)
			}
			log.Printf("Potential FVEK saved: %s_%d, containing %x", outputPath, i, key.Data)
		}
	}
	return nil
}

func isValidFVEK(data []byte) bool {
	// Implement validation logic
	// This is a simplified example; actual validation is more complex
	return len(data) == 32 // Assuming 256-bit key
}

func saveFVEK(data []byte, path string) error {
	return os.WriteFile(path, data, 0600)
}

func main() {
	// retrieve the path to the memory dump file from first argument
	memoryDumpPath := defaultMemoryDumpPath
	if len(os.Args) > 1 {
		memoryDumpPath = os.Args[1]
	}
	// Read memory dump
	log.Printf("Reading memory dump from: %s", memoryDumpPath)
	dump, err := readMemoryDump(memoryDumpPath)
	if err != nil {
		log.Fatalf("Failed to read memory dump")
	}
	sizeOfDump := len(dump)
	log.Printf("Memory dump read successfully. Size: %d bytes", sizeOfDump)

	// Search for FVE metadata
	log.Printf("Searching for FVE metadata...")
	metadata, err := searchFVEMetadata(dump)
	if err != nil {
		log.Fatalf("Failed to search for FVE metadata")
	}
	log.Printf("Found %d FVE metadata structures", len(metadata))

	// Extract potential FVEK
	log.Printf("Extracting potential FVEK...")
	keys, err := extractPotentialFVEK(dump, metadata)
	if err != nil {
		log.Fatalf("Failed to extract potential FVEK")
	}
	log.Printf("Extracted %d potential FVEK", len(keys))
	for i := 0; i < len(keys); i++ {
		log.Printf("Potential FVEK offset: %d, hex:%x", keys[i].Offset, keys[i].Offset)
		log.Printf("Potential FVEK data: %x", keys[i].Data)
		log.Printf("Potential FVEK data length: %d", len(keys[i].Data))
	}

	// Validate and save potential FVEK
	log.Printf("Validating and saving potential FVEK...")
	err = validateAndSaveFVEK(keys)
	if err != nil {
		log.Fatalf("Failed to validate and save FVEK")
	}
	log.Printf("Potential FVEK extraction completed successfully")
	log.Printf("You can try to decrypt the disk using the extracted FVEK with a tool like dislocker on Linux:")
	log.Printf("sudo dislocker -V /dev/sdX -k extracted_fvek.bin_2 --dislocker-file your_dislocker.img")

}
