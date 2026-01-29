package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// اسم ملف حفظ السيكرت
const SecretFile = "my_secret.txt"
const MtgURL = "https://github.com/9seconds/mtg/releases/download/v2.1.7/mtg-2.1.7-linux-amd64.tar.gz"

// متغير عالمي للسيكرت
var CurrentSecret = ""

func main() {
	// 1. تشغيل ويب سيرفر (يعرض السيكرت المحفوظ)
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "=== MTG Proxy Persistent ===\n\n")
			if CurrentSecret != "" {
				fmt.Fprintf(w, "STATUS: Running ✅\n")
				fmt.Fprintf(w, "SECRET: %s\n\n", CurrentSecret)
				fmt.Fprintf(w, "(This secret is saved and will be reused on restart)")
			} else {
				fmt.Fprintf(w, "STATUS: Initializing...\n")
			}
		})
		
		port := os.Getenv("PORT")
		if port == "" {
			port = "10000"
		}
		fmt.Printf("[-] Web Server listening on port %s\n", port)
		http.ListenAndServe(":"+port, nil)
	}()

	// 2. تشغيل النظام
	if err := runSystem(); err != nil {
		fmt.Printf("[!] Fatal Error: %v\n", err)
		// إبقاء البرنامج يعمل لكي نرى الخطأ في الويب لو احتجنا
		select {}
	}
}

func runSystem() error {
	// --- الخطوة 1: تجهيز المحرك ---
	binaryPath := "./mtg-2.1.7-linux-amd64/mtg"
	
	// التحقق مما إذا كان المحرك موجوداً (لتجنب التحميل المتكرر)
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		fmt.Println("[-] Downloading MTG Engine...")
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

		fmt.Println("[-] Extracting...")
		exec.Command("tar", "-xvf", "mtg.tar.gz").Run()
	} else {
		fmt.Println("[-] Engine already exists. Skipping download.")
	}

	os.Chmod(binaryPath, 0777)

	// --- الخطوة 2: إدارة السيكرت (الحفظ والاسترجاع) ---
	// هل الملف موجود؟
	if content, err := os.ReadFile(SecretFile); err == nil && len(content) > 0 {
		fmt.Println("[-] Found saved secret!")
		CurrentSecret = strings.TrimSpace(string(content))
	} else {
		fmt.Println("[-] No saved secret found. Generating new one...")
		genCmd := exec.Command(binaryPath, "generate-secret", "--hex", "google.com")
		var outBuf bytes.Buffer
		genCmd.Stdout = &outBuf
		if err := genCmd.Run(); err != nil {
			return fmt.Errorf("generation failed: %v", err)
		}
		CurrentSecret = strings.TrimSpace(outBuf.String())
		
		// حفظ السيكرت للمستقبل
		os.WriteFile(SecretFile, []byte(CurrentSecret), 0644)
		fmt.Println("[-] New secret generated and saved to 'my_secret.txt'")
	}

	fmt.Printf("[-] Using Secret: %s\n", CurrentSecret)
	fmt.Println("[-] Starting Proxy via Direct Command (No Config File)...")

	// --- الخطوة 3: التشغيل المباشر ---
	// نستخدم simple-run لأنه لا يحتاج ملف إعدادات (وهذا يحل مشكلة الخطأ السابق)
	// Go يمرر السيكرت بشكل آمن جداً
	cmd := exec.Command(binaryPath, "simple-run", "-b", "0.0.0.0:443", CurrentSecret)
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
