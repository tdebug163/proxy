import os
import subprocess
import requests
import tarfile
import sys

# سحب السيكرت من المتغير G
SECRET = os.getenv("G")
MTG_URL = "https://github.com/9seconds/mtg/releases/download/v2.1.7/mtg-2.1.7-linux-amd64.tar.gz"

def start_pure_proxy():
    if not SECRET:
        print("[!] Error: Environment variable 'G' is missing!", flush=True)
        return

    print(f"[-] Fetching MTG Engine...", flush=True)

    # 1. تحميل المحرك
    try:
        r = requests.get(MTG_URL, stream=True)
        with open("mtg.tar.gz", "wb") as f:
            f.write(r.content)
            
        # 2. فك الضغط
        with tarfile.open("mtg.tar.gz", "r:gz") as tar:
            tar.extractall()
            
        # 3. البحث عن الملف التنفيذي
        binary_path = None
        for root, dirs, files in os.walk("."):
            for file in files:
                if file == "mtg":
                    binary_path = os.path.join(root, file)
                    break
        
        if binary_path:
            os.chmod(binary_path, 0o777)
            print(f"[-] Engine Ready. Running Proxy on Port 443 with Secret from G...", flush=True)
            
            # 4. التشغيل
            cmd = f"{binary_path} simple-run -b 0.0.0.0:443 {SECRET}"
            subprocess.run(cmd, shell=True)
        else:
            print("[!] Error: MTG binary not found.", flush=True)
            
    except Exception as e:
        print(f"[!] Crash: {e}", flush=True)

if __name__ == "__main__":
    start_pure_proxy()
