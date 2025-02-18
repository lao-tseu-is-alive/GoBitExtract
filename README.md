# GoBitExtract 🔐  
**BitLocker FVEK Extractor from Memory Dumps**  

GoBitExtract is a simple program written in [Go](https://go.dev/) that allows you to search a memory dump file for a BitLocker FVEK (Full Volume Encryption Key). It is designed for **legitimate recovery purposes**, such as forensic investigations or system administration or try recovering your own data when you have lost the original key...

## ⚠️ Disclaimer  
🚨 **Legal Use Only!** 🚨  
- This project is **strictly intended** for ethical and legal usage.  
- Extracting encryption keys without proper authorization **may be illegal** in many jurisdictions.  
- The author is not responsible for any misuse of this software.  
- **Verify the legality** of using this tool in your country before proceeding.  

🛑 **If you are unsure whether you are allowed to use this software, DO NOT USE IT.**  

## 📌 Features  
✅ Scans memory dump files (`.mem`, `.raw`, `.vmem`) for potential BitLocker FVEKs  
✅ Supports memory dumps from Proxmox, VMware, and other hypervisors  
✅ Designed for forensic professionals and system administrators  

## 🚀 Installation
### Requirements
install [Go](https://go.dev/doc/install) on your system, then clone the repository and build the program:

```bash
git clone https://github.com/lao-tseu-is-alive/GoBitExtract.git
cd GoBitExtract
go build -o gobitextract gobitextract.go
```

### Run the built program
```bash
./gobitextract /path/to/your/memorydump.mem
```

## 📝 Dev Usage
```bash
go run gobitextract.go /path/to/your/memorydump.mem
```

## example output (with modified keys)
```bash
2025/02/18 16:35:23 Reading memory dump from: /dev/shm/memory_dump.mem
2025/02/18 16:35:27 Memory dump read successfully. Size: 8606650803 bytes
2025/02/18 16:35:27 Searching for FVE metadata...
2025/02/18 16:35:27 Found FVE metadata at offset 62083496, hex:3b351a8
2025/02/18 16:35:27 Found FVE metadata at offset 80012600, hex:4c4e538
2025/02/18 16:35:27 Found FVE metadata at offset 93606312, hex:59451a8
2025/02/18 16:35:27 Found FVE metadata at offset 150348294, hex:8f62206
2025/02/18 16:35:27 Found FVE metadata at offset 150381277, hex:8f6a2dd
2025/02/18 16:35:27 Found FVE metadata at offset 150399663, hex:8f6eaaf
2025/02/18 16:35:27 Found FVE metadata at offset 150484616, hex:8f83688
2025/02/18 16:35:27 Found FVE metadata at offset 150484680, hex:8f836c8
2025/02/18 16:35:29 Found FVE metadata at offset 2181622840, hex:8208ec38
2025/02/18 16:35:30 Found FVE metadata at offset 2373472029, hex:8d784f1d
2025/02/18 16:35:38 Found 10 FVE metadata structures
2025/02/18 16:35:38 Extracting potential FVEK...
2025/02/18 16:35:38 Extracted 10 potential FVEK
2025/02/18 16:35:38 Potential FVEK offset: 135424, hex:21100
2025/02/18 16:35:38 Potential FVEK data: 0000000000000000000000000000000000000000000000000000000000000000
2025/02/18 16:35:38 Potential FVEK data length: 32
2025/02/18 16:35:38 Potential FVEK offset: 256, hex:100
2025/02/18 16:35:38 Potential FVEK data: a811000000000000000000000000000000000000000000000080010000000000
2025/02/18 16:35:38 Potential FVEK data length: 32
2025/02/18 16:35:38 Potential FVEK offset: 131400, hex:20148
2025/02/18 16:35:38 Potential FVEK data: 89477033c049893e4c8d5c24ccccccbbbbbbbbaaaaaaaa7330498b7b38498be3
2025/02/18 16:35:38 Potential FVEK data length: 32
2025/02/18 16:35:38 Potential FVEK offset: 207587905, hex:c5f8a41
2025/02/18 16:35:38 Potential FVEK data: 0000000000000000000000000000000000000000000000000000000000000000
2025/02/18 16:35:38 Potential FVEK data length: 32
2025/02/18 16:35:38 Potential FVEK offset: 55478856, hex:34e8a48
2025/02/18 16:35:38 Potential FVEK data: 0000000000000000000000000000000000000000000000000000000000000000
2025/02/18 16:35:38 Potential FVEK data length: 32
2025/02/18 16:35:38 Potential FVEK offset: 205703233, hex:c42c841
2025/02/18 16:35:38 Potential FVEK data: 0000000000000000000000000000000000000000000000000000000000000000
2025/02/18 16:35:38 Potential FVEK data length: 32
2025/02/18 16:35:38 Potential FVEK offset: 256, hex:100
2025/02/18 16:35:38 Potential FVEK data: a811000000000000000000000000000000000000000000000080010000000000
2025/02/18 16:35:38 Potential FVEK data length: 32
2025/02/18 16:35:38 Potential FVEK offset: 256, hex:100
2025/02/18 16:35:38 Potential FVEK data: a811000000000000000000000000000000000000000000000080010000000000
2025/02/18 16:35:38 Potential FVEK data length: 32
2025/02/18 16:35:38 Potential FVEK offset: 135424, hex:21100
2025/02/18 16:35:38 Potential FVEK data: 0000000000000000000000000000000000000000000000000000000000000000
2025/02/18 16:35:38 Potential FVEK data length: 32
2025/02/18 16:35:38 Potential FVEK offset: 54954568, hex:3468a48
2025/02/18 16:35:38 Potential FVEK data: 672e83512c29c1cc0ea8b4b798ef7441eaca54a3e46ced722f29e50f9930818d
2025/02/18 16:35:38 Potential FVEK data length: 32
2025/02/18 16:35:38 Validating and saving potential FVEK...
2025/02/18 16:35:38 Potential FVEK saved: extracted_fvek.bin_0, containing 0000000000000000000000000000000000000000000000000000000000000000
2025/02/18 16:35:38 Potential FVEK saved: extracted_fvek.bin_1, containing a811000000000000000000000000000000000000000000000080010000000000
2025/02/18 16:35:38 Potential FVEK saved: extracted_fvek.bin_2, containing 89477033c049893e4c8d5c2460498b5b20498b6b28498b7330498b7b38498be3
2025/02/18 16:35:38 Potential FVEK saved: extracted_fvek.bin_3, containing 0000000000000000000000000000000000000000000000000000000000000000
2025/02/18 16:35:38 Potential FVEK saved: extracted_fvek.bin_4, containing 0000000000000000000000000000000000000000000000000000000000000000
2025/02/18 16:35:38 Potential FVEK saved: extracted_fvek.bin_5, containing 0000000000000000000000000000000000000000000000000000000000000000
2025/02/18 16:35:38 Potential FVEK saved: extracted_fvek.bin_6, containing a811000000000000000000000000000000000000000000000080010000000000
2025/02/18 16:35:38 Potential FVEK saved: extracted_fvek.bin_7, containing a811000000000000000000000000000000000000000000000080010000000000
2025/02/18 16:35:38 Potential FVEK saved: extracted_fvek.bin_8, containing 0000000000000000000000000000000000000000000000000000000000000000
2025/02/18 16:35:38 Potential FVEK saved: extracted_fvek.bin_9, containing 672e83512c29c1cc0ea8b4b798ef7441eaca54a3e46ced722f29e50f9930818d
2025/02/18 16:35:38 Potential FVEK extraction completed successfully
2025/02/18 16:35:38 You can try to decrypt the disk using the extracted FVEK with a tool like dislocker on Linux:
2025/02/18 16:35:38 sudo dislocker -V /dev/sdX -k extracted_fvek.bin_2 --dislocker-file your_dislocker.img

```

## 📚 References

