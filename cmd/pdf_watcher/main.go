// cmd/pdf_watcher/main.go
// Watcher de hot-reload para desarrollo de templates PDF.
// Uso: go run cmd/pdf_watcher/main.go
//
// Watches: internal/pdf/templates/
// On change: Runs go run cmd/pdf_test/main.go and opens the generated PDF
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

const (
	templatesPath   = "internal/pdf/templates"
	pdfTestPath     = "cmd/pdf_test/main.go"
	outputFileName  = "test_output.pdf"
	watchDebounceMs = 300
)

var (
	watchPath   = flag.String("path", templatesPath, "Path to watch for changes")
	pdfTestMain = flag.String("test", pdfTestPath, "Path to pdf_test main.go")
	empresa     = flag.String("empresa", "garfex", "Empresa to use for PDF generation")
	noOpen      = flag.Bool("no-open", false, "Don't open PDF automatically after generation")
	verbose     = flag.Bool("v", false, "Verbose output")
)

func main() {
	flag.Parse()

	log.SetFlags(0)
	log.SetPrefix("📺 ")

	log.Printf("🎯 PDF Template Watcher Started")
	log.Printf("   Watching: %s", *watchPath)
	log.Printf("   PDF Test: %s", *pdfTestMain)
	log.Printf("   Empresa:  %s", *empresa)
	log.Printf("")
	log.Printf("   Press Ctrl+C to stop")
	log.Printf("")

	// Verify watch path exists
	if _, err := os.Stat(*watchPath); os.IsNotExist(err) {
		log.Fatalf("❌ Watch path does not exist: %s", *watchPath)
	}

	// Verify pdf_test main.go exists
	if _, err := os.Stat(*pdfTestMain); os.IsNotExist(err) {
		log.Fatalf("❌ PDF test file does not exist: %s", *pdfTestMain)
	}

	// Create watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("❌ Error creating watcher: %v", err)
	}
	defer watcher.Close()

	// Watch the templates directory
	err = filepath.Walk(*watchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Watch directories only
		if info.IsDir() {
			if *verbose {
				log.Printf("🔍 Watching: %s", path)
			}
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("❌ Error adding watch: %v", err)
	}

	// Channel for debouncing
	changeCh := make(chan string, 10)

	// Start event handler goroutine
	go func() {
		var timer *time.Timer
		for path := range changeCh {
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(watchDebounceMs*time.Millisecond, func() {
				handleChange(path)
			})
		}
	}()

	// Start watching
	log.Printf("✅ Watching for .html and .css changes...")

	// Initial generation
	log.Printf("")
	handleChange("initial")

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Only process write and create events
			if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
				continue
			}

			// Only process .html and .css files
			ext := strings.ToLower(filepath.Ext(event.Name))
			if ext != ".html" && ext != ".css" {
				continue
			}

			if *verbose {
				log.Printf("📝 Change detected: %s", event.Name)
			}
			changeCh <- event.Name

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("❌ Watch error: %v", err)
		}
	}
}

func handleChange(trigger string) {
	if trigger != "initial" {
		log.Printf("")
		log.Printf("🔄 Regenerating PDF due to: %s", filepath.Base(trigger))
	}

	// Run pdf_test/main.go
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	args := []string{"run", *pdfTestMain, "-empresa=" + *empresa}
	cmd := exec.CommandContext(ctx, "go", args...)

	// Set working directory to project root
	cmd.Dir = getProjectRoot()

	// Capture output
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("❌ Error generating PDF: %v", err)
		if len(output) > 0 {
			log.Printf("   Output: %s", string(output))
		}
		return
	}

	if *verbose || trigger == "initial" {
		log.Printf("   %s", strings.ReplaceAll(string(output), "\n", "\n   "))
	}

	// Check if PDF was generated
	pdfPath := getProjectRoot() + "/" + outputFileName
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		log.Printf("⚠️  PDF file not found: %s", pdfPath)
		return
	}

	log.Printf("✅ PDF regenerated: %s", outputFileName)

	// Open PDF automatically (unless disabled)
	if !*noOpen {
		openPDF(pdfPath)
	}
}

func openPDF(path string) {
	var cmd *exec.Cmd

	// Detect OS and use appropriate command
	switch os := strings.ToLower(os.Getenv("GOOS")); os {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default: // linux and others
		// Try common PDF viewers in order
		viewers := []string{"xdg-open", "evince", "okular", "zathura", "mupdf", "gnome-open", "kde-open"}
		for _, viewer := range viewers {
			if _, err := exec.LookPath(viewer); err == nil {
				cmd = exec.Command(viewer, path)
				break
			}
		}
		if cmd == nil {
			log.Printf("⚠️  No PDF viewer found. Install xdg-open, evince, or okular")
			return
		}
	}

	if err := cmd.Start(); err != nil {
		log.Printf("⚠️  Error opening PDF: %v", err)
		return
	}

	log.Printf("📂 PDF opened in viewer")
}

func getProjectRoot() string {
	// Get the directory where the executable is running
	// Walk up to find go.mod
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "."
		}
		dir = parent
	}
}

// ANSI color codes for pretty output
var (
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	reset  = "\033[0m"
)

func init() {
	// Check if terminal supports colors
	if !isTerminal() {
		red, green, yellow, blue, reset = "", "", "", "", ""
	}
}

func isTerminal() bool {
	return isColorTerminal(os.Stdout.Fd())
}

func isColorTerminal(fd uintptr) bool {
	// Simple check - assume terminals are color capable
	// In production, use termenv or similar
	return true
}

func init() {
	// Override log to use colors
	log.SetFlags(0)
	log.SetPrefix("📺 ")

	// Add some color to the output functions
	fmt.Printf("%s🎯 PDF Template Watcher%s\n", blue, reset)
}
