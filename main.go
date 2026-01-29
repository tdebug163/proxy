package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// متغير عالمي لحفظ السيكرت وعرضه في صفحة الويب
var LiveSecret = "Initializing... Please wait."
const MtgURL = "https://github.com/9seconds/mtg/releases/download/v2.1.7/mtg-2.1.7-linux-amd64.tar.gz"

func main() {
	// 1. تشغيل الويب سيرفر (لعرض السيكرت لك)
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// تنسيق الصفحة لتكون واضحة
			fmt.Fprintf(w, "=== MTG Proxy Auto-Generated ===\n\n")
			fmt.Fprintf(w, "STATUS: Running 🔥\n")
			fmt.Fprintf(w, "PORT: 443\n")
			fmt.Fprintf(w, "SECRET: %s\n\n", LiveSecret)
			fmt.Fprintf(w, "Make sure to copy the secret above!")
		})
		
		port := os.Getenv("PORT")
		if port == "" {
			port = "10000"
		}
		fmt.Printf("[-] Web Server listening on port %s\n", port)
		http.ListenAndServe(":"+port, nil)
	}()

	// 2. البدء في عملية التجهيز
	if err := startSystem(); err != nil {
		fmt.Printf("[!] Fatal Error: %v\n", err)
		select {}
	}
}

func startSystem() error {
	fmt.Println("[-] Downloading MTG Engine...")
	
	// تحميل
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

	// فك ضغط
	fmt.Println("[-] Extracting...")
	exec.Command("tar", "-xvf", "mtg.tar.gz").Run()

	binaryPath := "./mtg-2.1.7-linux-amd64/mtg"
	os.Chmod(binaryPath, 0777)

	// --- الخطوة الحاسمة: توليد السيكرت ---
	fmt.Println("[-] Asking Engine to Generate Secret (FakeTLS - google.com)...")
	
	// نطلب من المحرك توليد سيكرت خاص بـ google.com عشان التمويه
	genCmd := exec.Command(binaryPath, "generate-secret", "--hex", "google.com")
	var outBuf bytes.Buffer
	genCmd.Stdout = &outBuf
	
	if err := genCmd.Run(); err != nil {
		return fmt.Errorf("failed to generate secret: %v", err)
	}

	// تنظيف السيكرت الناتج
	LiveSecret = strings.TrimSpace(outBuf.String())
	fmt.Printf("[-] Secret Generated Successfully: %s\n", LiveSecret)

	// --- كتابة ملف الإعدادات بالسيكرت الجديد ---
	fmt.Println("[-] Creating Config File...")
	configContent := fmt.Sprintf(`
bind-to = "0.0.0.0:443"

[users]
name = "auto_user"
secret = "%s"
`, LiveSecret)

	if err := os.WriteFile("mtg.toml", []byte(configContent), 0644); err != nil {
		return err
	}

	fmt.Println("[-] Engine Ready. Starting Proxy...")

	// تشغيل المحرك بملف الإعدادات
	cmd := exec.Command(binaryPath, "run", "mtg.toml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
