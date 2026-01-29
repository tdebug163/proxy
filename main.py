import os
import subprocess
import requests
import tarfile
import sys

# 1. سحب السيكرت وتنظيفه من المسافات المخفية
raw_secret = os.getenv("G", "")
SECRET = raw_secret.strip() # هذا السطر يقتل أي مسافة زائدة

MTG_URL = "https://github.com/9seconds/mtg/releases/download/v2.1.7/mtg-2.1.7-linux-amd64.tar.gz"

def start_pure_proxy():
    # فحص دقيق جداً
    if not SECRET:
        print("[!] Error: Variable 'G' is empty or missing!", flush=True)
        return
    
    # طباعة معلومات للتأكد (بدون كشف السيكرت كامل)
    print(f"[-] Debug: Secret Length is {len(SECRET)} characters.", flush=True)
    print(f"[-] Debug: Secret starts with '{SECRET[:2]}...'", flush=True)

    print(f"[-] Fetching MTG Engine...", flush=True)

    try:
        r = requests.get(MTG_URL, stream=True)
        with open("mtg.tar.gz", "wb") as f:
            f.write(r.content)
            
        with tarfile.open("mtg.tar.gz", "r:gz") as tar:
            tar.extractall()
            
        binary_path = None
        for root, dirs, files in os.walk("."):
            for file in files:
                if file == "mtg":
                    binary_path = os.path.join(root, file)
                    break
        
        if binary_path:
            os.chmod(binary_path, 0o777)
            print(f"[-] Engine Ready. Executing Proxy Direct Command...", flush=True)
            
            # --- التغيير الجذري هنا ---
            # تمرير الأمر كقائمة (List) يمنع ضياع السيكرت
            command_list = [
                binary_path,
                "simple-run",
                "-b", "0.0.0.0:443",
                SECRET
            ]
            
            # تشغيل مباشر بدون shell
            subprocess.run(command_list)
        else:
            print("[!] Error: MTG binary not found.", flush=True)
            
    except Exception as e:
        print(f"[!] Crash: {e}", flush=True)

if __name__ == "__main__":
    start_pure_proxy()
