import os
import asyncio
from mtprotoproxy import MTProtoProxy

# جلب السيكرت من ريندر
SECRET = os.getenv("MY_SECRET")

async def start_proxy():
    # هنا "user" هو اسم مستعار، والمهم هو السيكرت
    # المنفذ 443
    proxy = MTProtoProxy(users={"user": SECRET}, port=443)
    print("Starting Proxy on port 443...")
    await proxy.run()

if __name__ == "__main__":
    try:
        asyncio.run(start_proxy())
    except Exception as e:
        print(f"Error: {e}")
