package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

// إعدادات السيكرت والرابط
const (
	MySecret = "eeb83bb28ac66051d62d32557cde65e2" // السيكرت هنا ثابت ومضمون
	MtgURL   = "https://github.com/9seconds/mtg/releases/download/v2.1.7/mtg-2.1.7-linux-amd64.tar.gz"
)

func main() {
	// 1. تشغيل ويب سيرفر في الخلفية (Goroutine)
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Go Proxy Runner is Live!")
		})
		port := os.Getenv("PORT")
		if port == "" {
			port = "10000"
		}
		fmt.Printf("[-] Web Server listening on port %s\n", port)
		http.ListenAndServe(":"+port, nil)
	}()

	// 2. تحميل المحرك وتشغيله
	if err := runProxy(); err != nil {
		fmt.Printf("[!] Fatal Error: %v\n", err)
		os.Exit(1)
	}
}

func runProxy() error {
	fmt.Println("[-] Downloading MTG Engine...")
    
    // تحميل الملف
	resp, err := http.Get(MtgURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create("mtg.tar.gz")
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(out, resp.Body)

	// فك الضغط (باستخدام أمر tar للنظام لأنه أسرع وأسهل في go)
	fmt.Println("[-] Extracting...")
	exec.Command("tar", "-xvf", "mtg.tar.gz").Run()

    // البحث عن الملف التنفيذي
    // في الغالب بعد فك الضغط يكون في مجلد، لكن سنبحث عنه أو نشغل المسار المتوقع
    // لتسهيل الأمر سنستخدم المسار المتوقع للإصدار 2.1.7
    binaryPath := "./mtg-2.1.7-linux-amd64/mtg"

    // إعطاء صلاحية تنفيذ
    os.Chmod(binaryPath, 0777)

	fmt.Println("[-] Starting MTG Proxy on Port 443...")

    // تشغيل البروكسي
    // نمرر السيكرت كـ Argument مباشر، Go يتعامل معها بذكاء وما يضيعها
	cmd := exec.Command(binaryPath, "simple-run", "-b", "0.0.0.0:443", MySecret)
	
    // توجيه المخرجات للوج عشان تشوف الأخطاء
    cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
