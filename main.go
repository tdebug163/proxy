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

var LiveSecret = "Initializing..."
const MtgURL = "https://github.com/9seconds/mtg/releases/download/v2.1.7/mtg-2.1.7-linux-amd64.tar.gz"

func main() {
	// تشغيل الويب
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "=== MTG Proxy (Go Edition) ===\n\n")
			fmt.Fprintf(w, "STATUS: Running 🔥\n")
			fmt.Fprintf(w, "SECRET: %s\n\n", LiveSecret)
			fmt.Fprintf(w, "Copy the secret above to Telegram.")
		})
		
		port := os.Getenv("PORT")
		if port == "" {
			port = "10000"
		}
		fmt.Printf("[-] Web Server listening on port %s\n", port)
		http.ListenAndServe(":"+port, nil)
	}()

	if err := startSystem(); err != nil {
		fmt.Printf("[!] Fatal Error: %v\n", err)
		select {}
	}
}

func startSystem() error {
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

	// توليد السيكرت
	fmt.Println("[-] Generating Secret...")
	genCmd := exec.Command(binaryPath, "generate-secret", "--hex", "google.com")
	var outBuf bytes.Buffer
	genCmd.Stdout = &outBuf
	
	if err := genCmd.Run(); err != nil {
		return fmt.Errorf("failed to generate secret: %v", err)
	}

	LiveSecret = strings.TrimSpace(outBuf.String())
	fmt.Printf("[-] Secret Generated: %s\n", LiveSecret)

	// --- التصحيح الجذري هنا ---
	// استخدام [users.name] بدلاً من [[users]]
	// هذا التنسيق هو الوحيد الذي يقبله الإصدار 2.1.7
	fmt.Println("[-] Creating Config File (Fixed Format)...")
	
	configContent := fmt.Sprintf(`
bind-to = "0.0.0.0:443"

[users.auto_user]
secret = "%s"
`, LiveSecret)

	// كتابة الملف
	if err := os.WriteFile("mtg.toml", []byte(configContent), 0644); err != nil {
		return err
	}

	// طباعة محتوى الملف للتاكد في اللوج
	fmt.Println("[-] Config Content Preview:")
	fmt.Println(configContent)

	fmt.Println("[-] Engine Ready. Starting Proxy...")

	cmd := exec.Command(binaryPath, "run", "mtg.toml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
