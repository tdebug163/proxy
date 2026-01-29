import os
import asyncio
from MTPyProxy.proxy import MTProtoProxy

# جلب السيكرت من ريندر (الذي يبدأ بـ ee)
SECRET = os.getenv("MY_SECRET")

async def run_proxy():
    # المنفذ 443 هو المنفذ القياسي للبروكسي
    # المكتبة ستتعامل مع السيكرت (ee) وتفعل التشفير تلقائياً
    config = {
        "port": 443,
        "users": {"user1": SECRET},
        "display_stats": True
    }
    
    proxy = MTProtoProxy(config)
    print("[-] MTProto Proxy is starting...")
    print(f"[-] Global State: Active (USA Server)")
    await proxy.run()

if __name__ == "__main__":
    try:
        asyncio.run(run_proxy())
    except Exception as e:
        print(f"[!] Error: {e}")
