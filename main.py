import os
from mtprotoproxy import proxy

# جلب السيكرت من ريندر
SECRET = os.getenv("MY_SECRET")

# تشغيل البروكسي على المنفذ 443
if __name__ == "__main__":
    # المكتبة تتولى كل شيء داخلياً (التشفير، التوصيل، التحميل)
    proxy.run(port=443, users={ "user": SECRET })
