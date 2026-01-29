import os
import subprocess
import requests
import tarfile
from flask import Flask
import threading

# --- الإعدادات ---
SECRET = os.getenv("MY_SECRET")
# رابط تحميل أقوى نسخة من محرك MTG (ثابتة ولا تتغير)
MTG_URL = "https://github.com/9seconds/mtg/releases/download/v2.1.7/mtg-2.1.7-linux-amd64.tar.gz"

def download_and_run_mtg():
    print("[-] Downloading MTG Engine (The Beast)...")
    
    # تحميل المحرك
    r = requests.get(MTG_URL, stream=True)
    with open("mtg.tar.gz", "wb") as f:
        f.write(r.content)
    
    # فك الضغط
    with tarfile.open("mtg.tar.gz", "r:gz") as tar:
        tar.extractall()
    
    # البحث عن الملف التنفيذي وتشغيله
    # (الاسم قد يختلف قليلاً بعد فك الضغط لذا نبحث عنه)
    binary_path = None
    for root, dirs, files in os.walk("."):
        for file in files:
            if file == "mtg":
                binary_path = os.path.join(root, file)
                break
    
    if binary_path:
        print(f"[-] Engine found at: {binary_path}")
        os.chmod(binary_path, 0o777) # إعطاء صلاحية التشغيل
        
        # تشغيل البروكسي
        # الأمر: ./mtg simple-run -n 1.1.1.1:443 -b 0.0.0.0:443 SECRET
        print("[-] Starting Proxy on Port 443 with FakeTLS...")
        cmd = f"{binary_path} simple-run -b 0.0.0.0:443 {SECRET}"
        subprocess.run(cmd, shell=True)
    else:
        print("[!] Error: MTG binary not found inside the archive!")

# --- قسم الويب (عشان ريندر ما ينام) ---
app = Flask(__name__)

@app.route('/')
def home():
    return "MTG Proxy is Running 🔥"

def run_web():
    app.run(host='0.0.0.0', port=10000)

if __name__ == "__main__":
    # تشغيل الويب في خيط منفصل
    threading.Thread(target=run_web).start()
    
    # تشغيل البروكسي في الواجهة
    download_and_run_mtg()
