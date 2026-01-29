package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
    // تم حذف "time" من هنا لأنه كان سبب المشكلة
)

// السيكرت الثابت
const MySecret = "eeb83bb28ac66051d62d32557cde65e2"

// رابط المحرك
const MtgURL = "https://github.com/9seconds/mtg/releases/download/v2.1.7/mtg-2.1.7-linux-amd64.tar.gz"

func main() {
	// 1. تشغيل الويب سيرفر
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Go Proxy is Running 🔥")
		})
		
		port := os.Getenv("PORT")
		if port == "" {
			port = "10000"
		}
		fmt.Printf("[-] Web Server listening on port %s\n", port)
		http.ListenAndServe(":"+port, nil)
	}()

	// 2. تشغيل البروكسي
	if err := runProxy(); err != nil {
		fmt.Printf("[!] Fatal Error: %v\n", err)
		// نمنع البرنامج من الإغلاق لكي يبقى الويب شغالاً
		select {}
	}
}

func runProxy() error {
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

	binaryPath := "./mtg-2.1.7-linux-amd64/mtg"
	os.Chmod(binaryPath, 0777)

	fmt.Println("[-] Engine Ready. Starting Proxy on Port 443...")

	// تشغيل البروكسي
	cmd := exec.Command(binaryPath, "simple-run", "-b", "0.0.0.0:443", MySecret)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
