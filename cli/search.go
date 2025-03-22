package cli_interface

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/hydrocode-de/datailama/internal/db"
	"github.com/hydrocode-de/datailama/internal/ollama"
	"github.com/hydrocode-de/datailama/internal/sql"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

func searchAction(c *cli.Context) error {
	// extract the title argument
	title := strings.Join(c.Args().Slice(), " ")

	// get the optional arguments
	author := c.String("author")
	orderBy := c.String("order")
	direction := c.String("direction")
	limit := c.Int("limit")

	dbURL := c.String("database-url")
	dbManager, err := db.New(c.Context, dbURL)
	if err != nil {
		return err
	}
	defer dbManager.Close()

	result, err := dbManager.SearchPaperByTitle(c.Context, sql.SearchPaperByTitleParams{
		Title:     title,
		Author:    author,
		OrderBy:   orderBy,
		Direction: direction,
		Limit:     int32(limit),
	})
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetStyle(table.StyleColoredCyanWhiteOnBlack)
	t.AppendHeader(table.Row{"ID", "Title", "Author", "Year", "Citations", "Citations / Year"})
	for _, paper := range result {
		title := paper.Title
		if len(title) > 50 {
			title = title[:50] + "..."
		}
		t.AppendRow(table.Row{paper.ID, title, paper.Author, paper.Published.Time.Year(), paper.Citations, fmt.Sprintf("%.2f", paper.CitationsYear)})
	}

	fmt.Println(t.Render())

	return nil
}

func searchBodyAction(c *cli.Context) error {
	// extract the prompt argument
	prompt := c.Args().First()

	// get the optional arguments
	limit := c.Int("limit")
	colorCode := c.String("color")

	// Get database connection URL
	dbURL := c.String("database-url")
	dbManager, err := db.New(c.Context, dbURL)
	if err != nil {
		return err
	}
	defer dbManager.Close()

	// Get the ollama URL and create a new context with it
	ollamaURL := c.String("ollama-url")
	ctx := context.WithValue(c.Context, ollama.OllamaURLKey, ollamaURL)

	// Call embedding with the new context
	embedding, err := ollama.EmbedText(ctx, prompt)
	if err != nil {
		return err
	}

	// Use the same context for the database search
	result, err := dbManager.SearchPaperBody(ctx, sql.SearchPaperBodyParams{
		Embedding: embedding,
		Limit:     int32(limit),
	})
	if err != nil {
		return err
	}

	// Get terminal width for formatting
	termWidth := 100 // Default width if we can't determine it
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		termWidth = w
	}

	// Create appropriate table style based on color name
	tableStyle := table.StyleColoredDark
	colorNum := ""

	switch strings.ToLower(colorCode) {
	case "red":
		tableStyle = table.StyleColoredRedWhiteOnBlack
		colorNum = "31"
	case "green":
		tableStyle = table.StyleColoredGreenWhiteOnBlack
		colorNum = "32"
	case "yellow":
		tableStyle = table.StyleColoredYellowWhiteOnBlack
		colorNum = "33"
	case "blue":
		tableStyle = table.StyleColoredBlueWhiteOnBlack
		colorNum = "34"
	case "magenta":
		tableStyle = table.StyleColoredMagentaWhiteOnBlack
		colorNum = "35"
	case "cyan":
		tableStyle = table.StyleColoredCyanWhiteOnBlack
		colorNum = "36"
	default:
		// Default to cyan if an unrecognized color is provided
		tableStyle = table.StyleColoredCyanWhiteOnBlack
		colorNum = "36"
	}

	t := table.NewWriter()
	t.SetStyle(tableStyle)
	t.SetTitle("Search Results for: " + prompt)

	// Make the table more readable
	t.Style().Options.SeparateRows = true
	t.Style().Options.DrawBorder = true
	t.Style().Box.PaddingLeft = "  "
	t.Style().Box.PaddingRight = "  "

	t.AppendHeader(table.Row{"ID", "Score", "Title", "Year", "Citations", "Citations/Year"})

	// Configure the table columns
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "ID", Hidden: true},
		{Name: "Score", WidthMax: 8},
		{Name: "Title", WidthMax: 50},
		{Name: "Year", WidthMax: 6},
		{Name: "Citations", WidthMax: 10},
		{Name: "Citations/Year", WidthMax: 15},
	})

	// Group results by paper ID to avoid duplicate entries
	paperGroups := make(map[int64][]sql.SearchPaperBodyRow)
	for _, paper := range result {
		paperGroups[paper.ID] = append(paperGroups[paper.ID], paper)
	}

	// Process each paper
	for id, papers := range paperGroups {
		// Sort matches by cosine distance (most relevant first)
		sort.Slice(papers, func(i, j int) bool {
			// Get cosine distance for i
			var cosineI float64
			switch v := papers[i].CosineDistance.(type) {
			case float64:
				cosineI = v
			case float32:
				cosineI = float64(v)
			default:
				cosineStr := fmt.Sprintf("%v", papers[i].CosineDistance)
				parsedVal, err := strconv.ParseFloat(cosineStr, 64)
				if err == nil {
					cosineI = parsedVal
				}
			}

			// Get cosine distance for j
			var cosineJ float64
			switch v := papers[j].CosineDistance.(type) {
			case float64:
				cosineJ = v
			case float32:
				cosineJ = float64(v)
			default:
				cosineStr := fmt.Sprintf("%v", papers[j].CosineDistance)
				parsedVal, err := strconv.ParseFloat(cosineStr, 64)
				if err == nil {
					cosineJ = parsedVal
				}
			}

			return cosineI < cosineJ
		})

		paper := papers[0] // Use first match for metadata
		title := paper.Title
		if len(title) > 50 {
			title = title[:50] + "..."
		}

		// Add paper metadata row
		var cosineFormatted string
		switch v := paper.CosineDistance.(type) {
		case float64:
			cosineFormatted = fmt.Sprintf("%.4f", v)
		case float32:
			cosineFormatted = fmt.Sprintf("%.4f", float64(v))
		default:
			cosineFormatted = fmt.Sprintf("%v", paper.CosineDistance)
		}

		t.AppendRow(table.Row{
			id,
			cosineFormatted,
			title,
			paper.Published.Time.Year(),
			paper.Citations,
			fmt.Sprintf("%.2f", paper.CitationsYear),
		})

		// Format and add match excerpts
		matchWidth := termWidth - 20 // Give some margin
		for i, p := range papers {
			if i >= 2 { // Limit to 2 excerpts per paper to avoid clutter
				break
			}

			// Format the match text to wrap at line boundaries
			matchText := p.Match
			matchText = strings.TrimSpace(matchText)

			// Truncate very long matches
			if len(matchText) > 500 {
				matchText = matchText[:500] + "..."
			}

			// Format the match text into multiple lines of appropriate width
			formattedMatch := formatTextBlock(matchText, matchWidth)

			// Highlight search terms with bold formatting and accent color
			highlightedMatch := highlightSearchTerms(formattedMatch, prompt, colorNum)

			// For match text, use a different style to distinguish it
			coloredMatch := "\033[3m" + highlightedMatch + "\033[0m" // Italics for match text in terminals that support it

			t.AppendRow(table.Row{
				id,
				"",
				coloredMatch,
				"",
				"",
				"",
			}, table.RowConfig{})

			// Add a thin separator between matches
			if i < len(papers)-1 && i < 1 {
				t.AppendSeparator()
			}
		}

		// Add a strong separator between different papers
		t.AppendSeparator()
	}

	fmt.Println(t.Render())
	return nil
}

// formatTextBlock formats a text block to fit within a specific width
func formatTextBlock(text string, width int) string {
	if width <= 0 {
		width = 80
	}

	// Clean up the text by removing extra spaces, line breaks, etc.
	text = cleanupText(text)

	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		// If adding this word would exceed the width
		if len(currentLine)+len(word)+1 > width {
			if currentLine != "" {
				lines = append(lines, strings.TrimSpace(currentLine))
				currentLine = word
			} else {
				// Word is longer than width, need to break it
				lines = append(lines, word[:width-3]+"...")
			}
		} else {
			if currentLine == "" {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		}
	}

	if currentLine != "" {
		lines = append(lines, strings.TrimSpace(currentLine))
	}

	return strings.Join(lines, "\n")
}

// cleanupText removes extra spaces, normalizes line breaks, and cleans hyphenation
func cleanupText(text string) string {
	// Replace multiple spaces with a single space
	text = strings.Join(strings.Fields(text), " ")

	// Fix hyphenated words split across lines (common in PDF extracts)
	text = regexp.MustCompile(`(\w+)-\s+(\w+)`).ReplaceAllString(text, "$1$2")

	// Some specific cleanup for academic papers
	text = strings.ReplaceAll(text, "- ", "")
	text = strings.ReplaceAll(text, "−1", "-1") // Normalize unicode minus
	text = strings.ReplaceAll(text, "◦", "°")   // Replace degree symbol

	return text
}

// highlightSearchTerms highlights search terms in the text with bold formatting and accent color
func highlightSearchTerms(text string, terms string, colorCode string) string {
	// Color code should already be a number at this point
	if colorCode == "" {
		colorCode = "36" // Default cyan
	}

	words := strings.Fields(terms)
	for _, word := range words {
		if len(word) < 3 {
			continue // Skip very short words
		}
		// Case insensitive pattern with word boundaries
		pattern := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(word) + `\b`)
		text = pattern.ReplaceAllStringFunc(text, func(match string) string {
			return "\033[1;" + colorCode + "m" + match + "\033[0m" // Bold + Color
		})
	}
	return text
}
