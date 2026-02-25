package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) < 2 {
		runTUI()
		return
	}

	switch os.Args[1] {
	case "push":
		cmdPush(os.Args[2:])
	case "search":
		cmdSearch(os.Args[2:])
	case "list":
		cmdList()
	case "delete":
		cmdDelete(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\nusage: pino [push|search|list|delete]\n", os.Args[1])
		os.Exit(1)
	}
}

func runTUI() {
	pinos, err := listAllPinos()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	p := tea.NewProgram(newModel(pinos), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func cmdPush(args []string) {
	fs := flag.NewFlagSet("push", flag.ExitOnError)
	summary := fs.String("summary", "", "summary of the pino")
	prompt := fs.String("prompt", "", "original prompt")
	plan := fs.String("plan", "", "plan content (optional)")
	fs.Parse(args)

	if *summary == "" || *prompt == "" {
		fmt.Fprintln(os.Stderr, "usage: pino push --summary \"...\" --prompt \"...\" [--plan \"...\"]")
		os.Exit(1)
	}

	filename, err := savePino(*summary, *prompt, *plan)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(filename)
}

func cmdSearch(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: pino search <keyword>")
		os.Exit(1)
	}
	results, err := searchPinos(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(results)
}

func cmdList() {
	pinos, err := listAllPinos()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(pinos)
}

func cmdDelete(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: pino delete <filename>")
		os.Exit(1)
	}
	if err := deletePino(args[0]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("deleted")
}
