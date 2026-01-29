import os
import subprocess
import requests
import tarfile
import sys

# --- التشغيل الإجباري (Hardcoded Secret) ---
# وضعنا السيكرت هنا مباشرة لنضمن أن المحرك يراه 100%
# ولن نعتمد على المتغيرات التي تصل فارغة
SECRET = "eeb83bb28ac66051d62d32557cde65e2"

# رابط المحرك
MTG_URL = "https://github.com/9seconds/mtg/releases/download/v2.1.7/mtg-2.1.7-linux-amd64.tar.gz"

def force_start_proxy():
    print(f"[-] Force Mode: Using Hardcoded Secret: {SECRET[:4]}...", flush=True)
    print(f"[-] Fetching MTG Engine...", flush=True)

    try:
        # 1. تحميل المحرك
        r = requests.get(MTG_URL, stream=True)
        with open("mtg.tar.gz", "wb") as f:
            f.write(r.content)
            
        with tarfile.open("mtg.tar.gz", "r:gz") as tar:
            tar.extractall()
            
        # 2. البحث عن الملف التنفيذي
        binary_path = None
        for root, dirs, files in os.walk("."):
            for file in files:
                if file == "mtg":
                    binary_path = os.path.join(root, file)
                    break
        
        if binary_path:
            os.chmod(binary_path, 0o777)
            print(f"[-] Engine Ready. Executing...", flush=True)
            
            # 3. التشغيل المباشر
            # نمرر الأمر كقائمة (List) لتجنب مشاكل المسافات والاقتباسات
            cmd = [
                binary_path,
                "simple-run",
                "-b", "0.0.0.0:443",
                SECRET
            ]
            
            # تشغيل العملية
            subprocess.run(cmd)
        else:
            print("[!] Error: Binary not found.", flush=True)
            
    except Exception as e:
        print(f"[!] Crash: {e}", flush=True)

if __name__ == "__main__":
    force_start_proxy()
