# GoBitExtract 🔐  
**BitLocker FVEK Extractor from Memory Dumps**  

GoBitExtract is a simple program written in Go that allows you to search a memory dump file for a BitLocker FVEK (Full Volume Encryption Key). It is designed for **legitimate recovery purposes**, such as forensic investigations or system administration or try recovering your own data when you have lost the original key...

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
```bash
git clone https://github.com/lao-tseu-is-alive/GoBitExtract.git
cd GoBitExtract
go build -o gobitextract gobitextract.go
