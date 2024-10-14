package main

import (
	"fmt"
	manager "github.com/jf17/daemon-manager-systemd/service"
	"log"
	"net/http"
	"os"
)

var systemDRecord manager.SystemDRecord

func startDaemon() (string, error) {
	return systemDRecord.Start()
}

func stopDaemon() (string, error) {
	return systemDRecord.Stop()
}

func restartDaemon() (string, error) {
	return systemDRecord.Restart()
}

func checkDaemonStatus() (string, error) {
	return systemDRecord.Status()
}

func handleStartRequest(w http.ResponseWriter, _ *http.Request) {
	log.Printf("Stopping daemon: %s", systemDRecord.Name) // Log the action

	start, err := startDaemon()
	if err != nil {
		log.Printf("Error started daemon %s: %v", systemDRecord.Name, err) // Log error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Daemon %s started successfully", systemDRecord.Name) // Log success
	fmt.Fprintln(w, start)
}

func handleStopRequest(w http.ResponseWriter, _ *http.Request) {
	log.Printf("Stopping daemon: %s", systemDRecord.Name) // Log the action

	stop, err := stopDaemon()
	if err != nil {
		log.Printf("Error stopping daemon %s: %v", systemDRecord.Name, err) // Log error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Daemon %s stopped successfully", systemDRecord.Name) // Log success
	fmt.Fprintln(w, stop)
}

func handleRestartRequest(w http.ResponseWriter, _ *http.Request) {
	log.Printf("Restarting daemon: %s", systemDRecord.Name) // Log the action

	restart, err := restartDaemon()
	if err != nil {
		log.Printf("Error restarting daemon %s: %v", systemDRecord.Name, err) // Log error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Daemon %s restart successfully", systemDRecord.Name) // Log success
	fmt.Fprintln(w, restart)
}

func handleStatusRequest(w http.ResponseWriter, _ *http.Request) {
	log.Printf("Check status daemon: %s", systemDRecord.Name) // Log the action

	status, err := checkDaemonStatus()
	if err != nil {
		log.Printf("Error check status daemon %s: %v", systemDRecord.Name, err) // Log error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, status)
}

func main() {

	logFile, err := os.OpenFile("daemon-manager.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log logFile: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile) // Set output to the log logFile

	name := ""
	path := ""
	port := ""

	for i := 1; i < len(os.Args); i += 2 {
		switch os.Args[i] {
		case "-name":
			name = os.Args[i+1]
		case "-path":
			path = os.Args[i+1]
		case "-port":
			port = os.Args[i+1]
		default:
			fmt.Println("Invalid argument:", os.Args[i])
			return
		}
	}
	if name == "" || path == "" || port == "" {
		fmt.Println("Problems with arguments when launching the application.\n Example: sudo daemon-manager-systemd -port 8080 -name apache2 -path /lib/systemd/system/")
		log.Println("Problems with arguments when launching the application.")
		return
	}

	fmt.Println("name=" + name)
	fmt.Println("path=" + path)
	fmt.Println("port=" + port)

	systemDRecord = manager.SystemDRecord{Name: name, Path: path}

	http.HandleFunc("/"+name+"/start", handleStartRequest)
	http.HandleFunc("/"+name+"/stop", handleStopRequest)
	http.HandleFunc("/"+name+"/restart", handleRestartRequest)
	http.HandleFunc("/"+name+"/status", handleStatusRequest)
	fmt.Printf("Server listening on port %s...", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
