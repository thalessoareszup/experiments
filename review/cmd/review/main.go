package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//go:embed embedded/SKILL.md
var embeddedSkill string

//go:embed embedded/web/*
var embeddedWeb embed.FS

func main() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "skill":
		skillCmd(os.Args[2:])
	case "serve":
		serveCmd(os.Args[2:])
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "review - diff review skill + viewer")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  review skill [--out dir]")
	fmt.Fprintln(os.Stderr, "  review serve [--addr host:port] [--review path] [--patch path]")
	fmt.Fprintln(os.Stderr, "")
}

func skillCmd(args []string) {
	fs := flag.NewFlagSet("skill", flag.ExitOnError)
	outDir := fs.String("out", filepath.FromSlash("skill/review"), "output directory")
	_ = fs.Parse(args)

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		log.Fatalf("mkdir %s: %v", *outDir, err)
	}
	outPath := filepath.Join(*outDir, "SKILL.md")
	if err := os.WriteFile(outPath, []byte(embeddedSkill), 0o644); err != nil {
		log.Fatalf("write %s: %v", outPath, err)
	}
	fmt.Println(outPath)
}

func serveCmd(args []string) {
	flagSet := flag.NewFlagSet("serve", flag.ExitOnError)
	addr := flagSet.String("addr", "127.0.0.1:6767", "listen address")
	reviewPath := flagSet.String("review", "", "path to review JSON file")
	patchPath := flagSet.String("patch", "", "path to patch diff file (overrides diff.patch in review)")
	_ = flagSet.Parse(args)

	// Load review and patch if specified
	var reviewJSON, patchText []byte
	var err error

	if *reviewPath != "" {
		reviewJSON, err = os.ReadFile(*reviewPath)
		if err != nil {
			log.Fatalf("read review: %v", err)
		}

		// Parse review to get patch path if not overridden
		if *patchPath == "" {
			var review struct {
				Diff struct {
					Patch   string `json:"patch"`
					Unified string `json:"unified"`
				} `json:"diff"`
			}
			if err := json.Unmarshal(reviewJSON, &review); err != nil {
				log.Fatalf("parse review: %v", err)
			}
			if review.Diff.Patch != "" {
				// Resolve patch path relative to review file
				reviewDir := filepath.Dir(*reviewPath)
				*patchPath = filepath.Join(reviewDir, review.Diff.Patch)
			}
		}

		if *patchPath != "" {
			patchText, err = os.ReadFile(*patchPath)
			if err != nil {
				log.Fatalf("read patch: %v", err)
			}
		}
	}

	// Get embedded HTML
	htmlBytes, err := embeddedWeb.ReadFile("embedded/web/index.html")
	if err != nil {
		log.Fatalf("read embedded html: %v", err)
	}

	// Inject data if we have a review
	if reviewJSON != nil {
		injection := buildInjection(reviewJSON, patchText)
		htmlBytes = injectIntoHTML(htmlBytes, injection)
	}

	mux := http.NewServeMux()

	// Serve injected HTML at root
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "/index.html" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(htmlBytes)
	})

	// Serve local files under /files/ (for loading other reviews)
	mux.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("."))))

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("listen %s: %v", *addr, err)
	}

	actualAddr := ln.Addr().String()
	fmt.Printf("http://%s/\n", actualAddr)

	srv := &http.Server{Handler: mux}
	if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server: %v", err)
	}
}

func buildInjection(reviewJSON, patchText []byte) []byte {
	var buf bytes.Buffer
	buf.WriteString("<script>window.__REVIEW__ = ")
	buf.Write(reviewJSON)
	buf.WriteString(";\nwindow.__PATCH__ = ")
	// Encode patch as JSON string
	patchJSONBytes, _ := json.Marshal(string(patchText))
	buf.Write(patchJSONBytes)
	buf.WriteString(";</script>")
	return buf.Bytes()
}

func injectIntoHTML(html, injection []byte) []byte {
	// Insert just before </head>
	marker := []byte("</head>")
	idx := bytes.Index(html, marker)
	if idx == -1 {
		// Fallback: prepend to html
		return append(injection, html...)
	}
	result := make([]byte, 0, len(html)+len(injection))
	result = append(result, html[:idx]...)
	result = append(result, injection...)
	result = append(result, html[idx:]...)
	return result
}

func toSlashClean(p string) string {
	p = filepath.Clean(p)
	p = filepath.ToSlash(p)
	p = strings.TrimPrefix(p, "./")
	p = strings.TrimPrefix(p, "/")
	return p
}

func fsSub(fsys fs.FS, dir string) (fs.FS, error) {
	return fs.Sub(fsys, dir)
}
