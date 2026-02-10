package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	// جلب البورت الديناميكي من ريندر
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000" // الافتراضي لبيئة ريندر
	}

	// المسار الذي سيزوره المراقب (UptimeRobot أو غيره)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// استجابة 200 OK صريحة لضمان حالة "Active"
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Status: Bridge Active | Location: US Server | Signal: Stable")
	})

	fmt.Printf("Health-Check Server started on port %s\n", port)

	// بدء الاستماع للبورت
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Critical Error: %v\n", err)
		os.Exit(1)
	}
}
