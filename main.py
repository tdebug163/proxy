import os
import subprocess
import requests
import tarfile
import sys

# 1. سحب السيكرت وتنظيفه
raw_secret = os.getenv("G", "")
SECRET = raw_secret.strip()

# رابط المحرك
MTG_URL = "https://github.com/9seconds/mtg/releases/download/v2.1.7/mtg-2.1.7-linux-amd64.tar.gz"

def create_config_file():
    # إنشاء ملف إعدادات بصيغة TOML
    # هذه الطريقة تجبر المحرك على قراءة السيكرت بشكل صحيح 100%
    config_content = f"""
bind-to = "0.0.0.0:443"

[[users]]
name = "render_admin"
secret = "{SECRET}"
"""
    with open("config.toml", "w") as f:
        f.write(config_content)
    print("[-] Config file 'config.toml' created successfully.", flush=True)

def start_proxy_via_config():
    if not SECRET:
        print("[!] Error: Variable 'G' is empty!", flush=True)
        return
    
    # طباعة تحقق
    print(f"[-] Debug: Secret starts with '{SECRET[:2]}...' (Length: {len(SECRET)})", flush=True)

    print(f"[-] Downloading MTG Engine...", flush=True)

    try:
        # تحميل المحرك
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
            
            # 2. إنشاء ملف الإعدادات بدلاً من الاعتماد على سطر الأوامر
            create_config_file()
            
            print(f"[-] Engine Ready. Running using Config File...", flush=True)
            
            # 3. تشغيل المحرك باستخدام ملف الإعدادات
            # الأمر أصبح: ./mtg run config.toml
            cmd = [binary_path, "run", "config.toml"]
            
            subprocess.run(cmd)
        else:
            print("[!] Error: MTG binary not found.", flush=True)
            
    except Exception as e:
        print(f"[!] Crash: {e}", flush=True)

if __name__ == "__main__":
    start_proxy_via_config()
