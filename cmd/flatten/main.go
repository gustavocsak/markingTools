package main

import (
	"flag"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/gustavocsak/markingTools/internal/flattener"
	"os"
	"strings"
)

var (
	doc = lipgloss.NewStyle().Padding(1, 2)

	header = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("212")).
		MarginBottom(1)

	// command name style
	cmd = lipgloss.NewStyle().
		Bold(true).
		Background(lipgloss.Color("57")).
		Foreground(lipgloss.Color("230"))

	flagStyle = lipgloss.NewStyle().Bold(true).MarginRight(1)

	descStyle = lipgloss.NewStyle().Faint(true)

	exampleCmd  = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).MarginLeft(2)
	exampleDesc = lipgloss.NewStyle().Faint(true).MarginLeft(2).MarginBottom(1)

	// example file structure
	exampleHeader = lipgloss.NewStyle().Bold(true).MarginLeft(2)
	fileTreeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244")).
			MarginLeft(4).
			MarginBottom(1)
)

func printHelp() {
	var b strings.Builder

	b.WriteString("A tool to flatten nested subdirectories, moving all files to the parent directory.\n\n")

	usage := fmt.Sprintf("%s %s\n\n",
		cmd.Render("flatten"),
		"[OPTIONS] [DIRECTORIES...]",
	)
	b.WriteString(header.Render("USAGE:") + usage)

	b.WriteString(header.Render("FLAGS:") + "\n")

	flags := [][2]string{}
	flag.CommandLine.VisitAll(func(f *flag.Flag) {
		// e and exclude will be 2 different flags
		// must add support for long flags
		name := fmt.Sprintf("  -%s", f.Name)
		if f.Name == "e" {
			name = fmt.Sprintf("  -%s, --exclude", f.Name)
		}
		if f.Name == "help" {
			name = "  -h, --help"
		}

		flags = append(flags, [2]string{name, f.Usage})
	})

	for _, f := range flags {
		b.WriteString(flagStyle.Render(f[0]) + descStyle.Render(f[1]) + "\n")
	}

	b.WriteString("\n" + header.Render("EXAMPLES:") + "\n")

	// visual structure example
	b.WriteString(exampleHeader.Render("Before:") + "\n")
	b.WriteString(fileTreeStyle.Render("bob_project/\n" +
		"└── part1/\n" +
		"    └── project.cpp\n" +
		"    └── project.h\n" +
		"└── todo/\n" +
		"    └── part2/\n" +
		"        └── README.md\n" +
		"    └── main.go\n"))

	b.WriteString(exampleCmd.Render("flatten bob_project/") + "\n\n")

	b.WriteString(exampleHeader.Render("After:") + "\n")
	b.WriteString(fileTreeStyle.Render("bob_project/\n" +
		"├── main.go\n" +
		"├── project.cpp\n" +
		"├── project.h\n" +
		"└── README.md\n"))

	b.WriteString("\n")

	// exclude flag
	ex2Cmd := exampleCmd.Render("flatten -e \"zip,log\" .")
	ex2Desc := exampleDesc.Render("Flattens the current directory, but skips all '.zip' and '.log' files.")
	b.WriteString(ex2Cmd + "\n" + ex2Desc + "\n")

	ex3Cmd := exampleCmd.Render("flatten")
	ex3Desc := exampleDesc.Render("Flattens all subdirectories in the current directory (./).")
	b.WriteString(ex3Cmd + "\n" + ex3Desc + "\n")

	fmt.Fprintln(flag.CommandLine.Output(), doc.Render(b.String()))
}

func main() {
	var exclude string
	flag.Usage = printHelp
	flag.StringVar(&exclude, "e", "", "exclude file extension")

	var help bool
	flag.BoolVar(&help, "help", false, "Show this help message.")

	flag.Parse()

	log.SetLevel(log.DebugLevel)
	log.Info("Flat: starting flattening directories")

	var excludeList []string
	if exclude != "" {
		log.Debugf("exclude: %s", exclude)
		excludeList = strings.Split(exclude, ",")
	}

	excludeSet := make(map[string]struct{})
	for _, ext := range excludeList {
		cleanedExt := strings.TrimPrefix(ext, ".")
		if cleanedExt != "" {
			excludeSet[cleanedExt] = struct{}{}
		}
	}

	dirsToFlatten := flag.Args()

	if len(dirsToFlatten) == 0 {
		// Flatten all subdirectories in the current directory.
		log.Info("No directory specified. Scanning current directory for subfolders...")

		parentEntries, err := os.ReadDir(".")
		if err != nil {
			log.Errorf("Failed to read current directory: %v", err)
			os.Exit(1)
		}

		processedCount := 0
		for _, entry := range parentEntries {
			if entry.IsDir() {
				dirPath := entry.Name()
				log.Infof("Running flattener on: %s", dirPath)
				flattener.Run(dirPath, excludeSet)
				processedCount++
			}
		}

		if processedCount == 0 {
			log.Info("No subdirectories found to flatten.")
		}

	} else {
		// Flatten only the directories listed by the user.
		log.Infof("Processing %d specified director(y/ies)...", len(dirsToFlatten))

		for _, dirPath := range dirsToFlatten {
			info, err := os.Stat(dirPath)
			if err != nil {
				log.Errorf("Cannot access path %s: %v", dirPath, err)
				continue
			}

			if !info.IsDir() {
				log.Warnf("Skipping %s, it is not a directory.", dirPath)
				continue
			}

			log.Infof("Running flattener on: %s", dirPath)
			flattener.Run(dirPath, excludeSet)
		}
	}

	log.Info("Flat: finished processing directories")
}
