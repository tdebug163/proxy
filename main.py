import os
import subprocess
import requests
import tarfile
import threading
from flask import Flask

# 1. سحب السيكرت
raw_secret = os.getenv("G", "")
SECRET = raw_secret.strip()

# رابط المحرك (Go)
MTG_URL = "https://github.com/9seconds/mtg/releases/download/v2.1.7/mtg-2.1.7-linux-amd64.tar.gz"

def setup_and_run_proxy():
    # فحص الأمان
    if not SECRET:
        print("[!] مصيبة: المتغير G فارغ أو غير موجود!", flush=True)
        return

    print(f"[-] Fetching MTG Binary...", flush=True)

    try:
        # تحميل وفك الضغط
        r = requests.get(MTG_URL, stream=True)
        with open("mtg.tar.gz", "wb") as f:
            f.write(r.content)
            
        with tarfile.open("mtg.tar.gz", "r:gz") as tar:
            tar.extractall()
            
        # البحث عن الملف
        binary_path = None
        for root, dirs, files in os.walk("."):
            for file in files:
                if file == "mtg":
                    binary_path = os.path.join(root, file)
                    break
        
        if binary_path:
            # إعطاء صلاحية التشغيل
            os.chmod(binary_path, 0o777)
            
            # --- الحل الجذري: إنشاء ملف تشغيل Shell ---
            # نكتب الأمر والسيكرت داخل ملف نصي تنفيذي
            # هذا يضمن أن النظام يرى السيكرت كنص ثابت 100%
            sh_content = f"""#!/bin/bash
# تشغيل البروكسي مع طباعة الأمر للتأكد
echo "[-] Executing Proxy Command..."
{binary_path} simple-run -b 0.0.0.0:443 "{SECRET}"
"""
            # حفظ ملف run.sh
            with open("run.sh", "w") as f:
                f.write(sh_content)
            
            # إعطاء صلاحية التشغيل للملف
            os.chmod("run.sh", 0o777)
            
            print("[-] run.sh created. Launching proxy via Shell...", flush=True)
            
            # تشغيل الملف
            subprocess.run(["./run.sh"])
            
        else:
            print("[!] Error: Binary not found.", flush=True)
            
    except Exception as e:
        print(f"[!] Crash: {e}", flush=True)

# --- الويب ---
app = Flask(__name__)

@app.route('/')
def home():
    return "Proxy Alive"

def run_web():
    app.run(host='0.0.0.0', port=10000)

if __name__ == "__main__":
    # تشغيل الويب
    threading.Thread(target=run_web).start()
    
    # تشغيل البروكسي
    setup_and_run_proxy()
