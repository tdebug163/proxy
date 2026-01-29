import os
import asyncio
from mtprotoproxy import ProxyServer

# جلب السيكرت من ريندر
SECRET = os.getenv("MY_SECRET")

async def main():
    print("[-] Starting MTProto-PY Engine...")
    print(f"[-] Location: USA (Render Server)")
    
    # إعدادات البروكسي: المنفذ 443 والسيكرت الخاص بك
    # المكتبة تدعم الـ Fake TLS تلقائياً إذا بدأ السيكرت بـ ee
    server = ProxyServer(
        port=443,
        users={"user": SECRET}
    )
    
    await server.start()
    print("[+] Proxy is Live and Waiting for Telegram Connections!")

if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        pass
    except Exception as e:
        print(f"[!] Error: {e}")
