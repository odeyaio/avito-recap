package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"slices"
	"strings"

	"avito-recap/internal/catalog"
	"avito-recap/internal/engine"
	"avito-recap/internal/model"
)

type DefaultActionResolver struct{}

func (DefaultActionResolver) Resolve(
	_ context.Context,
	action catalog.DefaultAction,
	dataset model.Dataset,
	evidence []engine.Evidence,
) (model.RecapNextAction, error) {
	const op = "service.DefaultActionResolver.Resolve"

	href := action.Href
	target := map[string]any{"type": action.TargetType}

	if categoryCode := topCategoryCode(dataset); categoryCode != "" && strings.Contains(href, "{category}") {
		href = strings.ReplaceAll(href, "{category}", url.QueryEscape(categoryCode))
		target["categoryCode"] = categoryCode
	} else {
		href = strings.ReplaceAll(href, "?q={category}", "")
	}

	targetJSON, err := json.Marshal(target)
	if err != nil {
		return model.RecapNextAction{}, fmt.Errorf("%s: marshal action target: %w", op, err)
	}
	evidenceJSON, err := json.Marshal(evidence)
	if err != nil {
		return model.RecapNextAction{}, fmt.Errorf("%s: marshal action evidence: %w", op, err)
	}

	return model.RecapNextAction{
		Code:     action.Code,
		Href:     href,
		Target:   targetJSON,
		Evidence: evidenceJSON,
	}, nil
}

func topCategoryCode(dataset model.Dataset) string {
	listingCategories := make(map[string]string, len(dataset.Listings))
	categoryCodes := make(map[string]string, len(dataset.Categories))
	for _, category := range dataset.Categories {
		categoryCodes[category.ID.String()] = category.Code
	}
	for _, listing := range dataset.Listings {
		listingCategories[listing.ID.String()] = listing.CategoryID.String()
	}

	counts := make(map[string]int)
	for _, event := range dataset.Events {
		if event.UserID != dataset.User.ID || !dataset.Period.Contains(event.OccurredAt) {
			continue
		}
		categoryID := ""
		if event.CategoryID != nil {
			categoryID = event.CategoryID.String()
		} else if event.ListingID != nil {
			categoryID = listingCategories[event.ListingID.String()]
		}
		if code := categoryCodes[categoryID]; code != "" {
			counts[code]++
		}
	}

	codes := make([]string, 0, len(counts))
	for code := range counts {
		codes = append(codes, code)
	}
	slices.Sort(codes)
	bestCode, bestCount := "", 0
	for _, code := range codes {
		if counts[code] > bestCount {
			bestCode, bestCount = code, counts[code]
		}
	}
	return bestCode
}
