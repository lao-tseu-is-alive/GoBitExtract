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
2025/02/18 17:38:37 Reading memory dump from: /dev/shm/memory_dump.mem
2025/02/18 17:38:40 Memory dump read successfully. Size: 8606650803 bytes
2025/02/18 17:38:40 Searching for FVEK...
2025/02/18 17:38:40 Found FVE metadata at offset 62083496 (0x3b351a8)
2025/02/18 17:38:40 Skipping structure at 0x3b351a8: version mismatch (262148)
2025/02/18 17:38:40 Found FVE metadata at offset 80012600 (0x4c4e538)
2025/02/18 17:38:40 Skipping structure at 0x4c4e538: version mismatch (0)
2025/02/18 17:38:40 Found FVE metadata at offset 93606312 (0x59451a8)
2025/02/18 17:38:40 Skipping structure at 0x59451a8: version mismatch (262148)
2025/02/18 17:38:40 Found FVE metadata at offset 150348294 (0x8f62206)
2025/02/18 17:38:40 Skipping structure at 0x8f62206: version mismatch (3226897737)
2025/02/18 17:38:40 Found FVE metadata at offset 150381277 (0x8f6a2dd)
2025/02/18 17:38:40 Skipping structure at 0x8f6a2dd: version mismatch (944146760)
2025/02/18 17:38:40 Found FVE metadata at offset 150399663 (0x8f6eaaf)
2025/02/18 17:38:40 Potential FVEK found at offset 0x17032817: 89477033c049893e4c8d5c24ccccccbbbbbbbbaaaaaaaa7330498b7b38498be3
2025/02/18 17:38:40 Found FVE metadata at offset 150484616 (0x8f83688)
2025/02/18 17:38:40 Skipping structure at 0x8f83688: version mismatch (0)
2025/02/18 17:38:40 Found FVE metadata at offset 150484680 (0x8f836c8)
2025/02/18 17:38:40 Skipping structure at 0x8f836c8: version mismatch (0)
2025/02/18 17:38:41 Found FVE metadata at offset 2181622840 (0x8208ec38)
2025/02/18 17:38:41 Skipping structure at 0x8208ec38: version mismatch (262148)
2025/02/18 17:38:41 Found FVE metadata at offset 2373472029 (0x8d784f1d)
2025/02/18 17:38:41 Skipping structure at 0x8d784f1d: version mismatch (2696842572)
2025/02/18 17:38:41 Extracted 1 valid FVEKs
2025/02/18 17:38:41 Saved valid FVEK to: extracted_fvek.bin_0
2025/02/18 17:38:41 Potential FVEK extraction completed successfully.
2025/02/18 17:38:41 You can try to decrypt the disk using the extracted FVEK with a tool like dislocker on Linux:
2025/02/18 17:38:41 sudo dislocker -V /dev/sdX -k extracted_fvek.bin_0 --dislocker-file your_dislocker.img
```

## 📚 References
[Neodyme findings of early 2025](https://neodyme.io/en/blog/bitlocker_screwed_without_a_screwdriver/#step-3b-exploiting-the-linux-kernel)
