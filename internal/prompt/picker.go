// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package prompt

import (
	"fmt"

	"github.com/charmbracelet/huh"

	"github.com/amrelsagaei/updock/internal/hub"
)

const chooseVersionSentinel = "__choose_version__"

// SelectFunc runs an interactive select and returns the chosen value.
// Replaceable in tests.
var SelectFunc = runHuhSelect

func runHuhSelect(title string, options []huh.Option[string]) (string, error) {
	var selected string
	err := huh.NewSelect[string]().
		Title(title).
		Options(options...).
		Value(&selected).
		Run()
	return selected, err
}

// ImageSelection holds the result of PickImage.
type ImageSelection struct {
	Result        hub.RankedResult
	ChooseVersion bool
}

// PickImage shows an interactive list of ranked search results.
func PickImage(results []hub.RankedResult) (ImageSelection, error) {
	if len(results) == 0 {
		return ImageSelection{}, fmt.Errorf("no results to pick from")
	}

	options := BuildImageOptions(results)
	selected, err := SelectFunc("Pick an image", options)
	if err != nil {
		return ImageSelection{}, fmt.Errorf("image selection: %w", err)
	}

	return ResolveImageSelection(selected, results), nil
}

// ResolveImageSelection maps a selected value back to an ImageSelection.
func ResolveImageSelection(selected string, results []hub.RankedResult) ImageSelection {
	if selected == chooseVersionSentinel {
		return ImageSelection{Result: results[0], ChooseVersion: true}
	}

	for _, r := range results {
		if r.RepoName == selected {
			return ImageSelection{Result: r}
		}
	}

	return ImageSelection{Result: results[0]}
}

// BuildImageOptions creates the huh option list from ranked results.
func BuildImageOptions(results []hub.RankedResult) []huh.Option[string] {
	opts := make([]huh.Option[string], 0, len(results)+1)

	best := results[0]
	opts = append(opts,
		huh.NewOption(
			fmt.Sprintf("%s:latest  [best match] [%s]", best.RepoName, best.Badge),
			best.RepoName,
		),
		huh.NewOption(
			fmt.Sprintf("%s (choose version)", best.RepoName),
			chooseVersionSentinel,
		),
	)

	for _, r := range results[1:] {
		label := fmt.Sprintf("%s:latest  [%s]", r.RepoName, r.Badge)
		if r.PullCount > 0 {
			label = fmt.Sprintf("%s:latest  %s pulls  [%s]", r.RepoName, FormatCount(r.PullCount), r.Badge)
		}
		opts = append(opts, huh.NewOption(label, r.RepoName))
	}

	return opts
}

// PickTag shows an interactive list of tags for an image.
func PickTag(tags []hub.Tag) (hub.Tag, error) {
	if len(tags) == 0 {
		return hub.Tag{}, fmt.Errorf("no tags to pick from")
	}

	options := BuildTagOptions(tags)
	selected, err := SelectFunc("Pick a version", options)
	if err != nil {
		return hub.Tag{}, fmt.Errorf("tag selection: %w", err)
	}

	return ResolveTagSelection(selected, tags), nil
}

// ResolveTagSelection maps a selected tag name back to a Tag.
func ResolveTagSelection(selected string, tags []hub.Tag) hub.Tag {
	for _, t := range tags {
		if t.Name == selected {
			return t
		}
	}
	return tags[0]
}

// BuildTagOptions creates the huh option list from tags.
func BuildTagOptions(tags []hub.Tag) []huh.Option[string] {
	options := make([]huh.Option[string], 0, len(tags))
	for _, t := range tags {
		label := t.Name
		if t.FullSize > 0 {
			label = fmt.Sprintf("%-20s  %s", t.Name, FormatBytes(t.FullSize))
		}
		options = append(options, huh.NewOption(label, t.Name))
	}
	return options
}

// FormatCount formats a number as a human-readable count string.
func FormatCount(n int) string {
	switch {
	case n >= 1_000_000_000:
		return fmt.Sprintf("%.1fB", float64(n)/1_000_000_000)
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	case n >= 1_000:
		return fmt.Sprintf("%.1fK", float64(n)/1_000)
	default:
		return fmt.Sprintf("%d", n)
	}
}

// FormatBytes formats bytes as a human-readable size string.
func FormatBytes(b int64) string {
	switch {
	case b >= 1<<30:
		return fmt.Sprintf("%.1f GB", float64(b)/(1<<30))
	case b >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(b)/(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(b)/(1<<10))
	default:
		return fmt.Sprintf("%d B", b)
	}
}
