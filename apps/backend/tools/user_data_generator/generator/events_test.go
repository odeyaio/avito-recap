package generator

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"avito-recap/internal/model"
)

func TestEventTimeOrder(t *testing.T) {
	categories := DefaultCategories()
	profiles := DefaultUserConfig()
	
	// Generate a small dataset
	users := GenerateUserData(42, 10, categories)
	sellers := SelectSellerUsers(42, users)
	listings := GenerateListingData(42, 50, categories, sellers, time.Now().UTC().AddDate(0, 0, 30))
	
	activityEvents, deals, reviews := GenerateUserEvents(42, users, profiles, listings, categories)
	
	// Test 1: All events should have OccurredAt after user registration
	for _, event := range activityEvents {
		userIdx := -1
		for i, user := range users {
			if user.ID == event.UserID {
				userIdx = i
				break
			}
		}
		if userIdx >= 0 && !users[userIdx].RegisteredAt.IsZero() {
			if event.OccurredAt.Before(users[userIdx].RegisteredAt) {
				t.Errorf("Event %s occurred before user registration for user %s", event.Type, event.UserID)
			}
		}
	}
	
	// Test 2: Edit events should occur after listing publish and before close
	editEvents := GenerateListingEditEvents(100, listings)
	for _, event := range editEvents {
		if event.Type != model.EventType(model.EventTypeEdit) {
			continue
		}
		
		listingIdx := -1
		for i, listing := range listings {
			if listing.ID == *event.ListingID {
				listingIdx = i
				break
			}
		}
		if listingIdx >= 0 {
			listing := listings[listingIdx]
			if event.OccurredAt.Before(listing.PublishedAt) {
				t.Errorf("Edit event occurred before listing publish for listing %s", listing.ID)
			}
			if listing.ClosedAt != nil && event.OccurredAt.After(*listing.ClosedAt) {
				t.Errorf("Edit event occurred after listing close for listing %s", listing.ID)
			}
		}
	}
	
	// Test 3: Share events should occur after view for the same listing
	viewEvents := make(map[uuid.UUID]time.Time) // listingID -> earliest view time
	for _, event := range activityEvents {
		if event.Type == model.EventType(model.EventTypeView) && event.ListingID != nil {
			if existing, ok := viewEvents[*event.ListingID]; !ok || event.OccurredAt.Before(existing) {
				viewEvents[*event.ListingID] = event.OccurredAt
			}
		}
	}
	
	for _, event := range activityEvents {
		if event.Type == model.EventType(model.EventTypeShare) && event.ListingID != nil {
			if viewTime, ok := viewEvents[*event.ListingID]; ok {
				if event.OccurredAt.Before(viewTime) {
					t.Errorf("Share event occurred before view for listing %s", *event.ListingID)
				}
			}
		}
	}
	
	// Test 4: Cancel deal events should occur after contact for the same listing
	contactEvents := make(map[uuid.UUID]time.Time) // listingID -> earliest contact time
	for _, event := range activityEvents {
		if event.Type == model.EventType(model.EventTypeContact) && event.ListingID != nil {
			if existing, ok := contactEvents[*event.ListingID]; !ok || event.OccurredAt.Before(existing) {
				contactEvents[*event.ListingID] = event.OccurredAt
			}
		}
	}
	
	for _, event := range activityEvents {
		if event.Type == model.EventType(model.EventTypeCancelDeal) && event.ListingID != nil {
			if contactTime, ok := contactEvents[*event.ListingID]; ok {
				if event.OccurredAt.Before(contactTime) {
					t.Errorf("Cancel deal event occurred before contact for listing %s", *event.ListingID)
				}
			}
		}
	}
	
	// Test 5: Deals should occur after contact for the same listing
	for _, deal := range deals {
		contactTime, ok := contactEvents[deal.ListingID]
		if !ok {
			t.Logf("No contact event found for deal listing %s", deal.ListingID)
			continue
		}
		if deal.CreatedAt.Before(contactTime) {
			t.Errorf("Deal created before contact for listing %s", deal.ListingID)
		}
	}
	
	// Test 6: Reviews should occur after deal creation
	for _, review := range reviews {
		if review.DealID == nil {
			continue
		}
		
		dealIdx := -1
		for i, deal := range deals {
			if deal.ID == *review.DealID {
				dealIdx = i
				break
			}
		}
		if dealIdx >= 0 {
			deal := deals[dealIdx]
			if review.CreatedAt.Before(deal.CreatedAt) {
				t.Errorf("Review created before deal for deal %s", deal.ID)
			}
		}
	}
	
	// Test 7: Within the same intent, events should be in order
	// Group events by user and check per-intent ordering
	userEvents := make(map[uuid.UUID][]model.ActivityEvent)
	for _, event := range activityEvents {
		userEvents[event.UserID] = append(userEvents[event.UserID], event)
	}
	
	for userID, events := range userEvents {
		// Sort by time
		sorted := make([]model.ActivityEvent, len(events))
		copy(sorted, events)
		
		// Simple bubble sort for testing
		for i := 0; i < len(sorted); i++ {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[i].OccurredAt.After(sorted[j].OccurredAt) {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}
		
		// Check that Search comes before View for the same category if they're close in time
		var lastSearch *model.ActivityEvent
		for i := 0; i < len(sorted); i++ {
			if sorted[i].Type == model.EventType(model.EventTypeSearch) {
				lastSearch = &sorted[i]
			} else if sorted[i].Type == model.EventType(model.EventTypeView) && lastSearch != nil {
				// If view is within 1 hour of search, it should be after
				if sorted[i].OccurredAt.Sub(lastSearch.OccurredAt) < time.Hour {
					if sorted[i].OccurredAt.Before(lastSearch.OccurredAt) {
						t.Errorf("View occurred before search for user %s", userID)
					}
				}
			}
		}
	}
	
	t.Logf("Passed time order tests with %d users, %d listings, %d events, %d deals, %d reviews",
		len(users), len(listings), len(activityEvents), len(deals), len(reviews))
}

func TestBusinessLogic(t *testing.T) {
	categories := DefaultCategories()
	profiles := DefaultUserConfig()
	
	// Generate a larger dataset to ensure we have enough data for business logic tests
	users := GenerateUserData(42, 50, categories)
	sellers := SelectSellerUsers(42, users)
	listings := GenerateListingData(42, 200, categories, sellers, time.Now().UTC().AddDate(0, 0, 30))
	
	activityEvents, deals, reviews := GenerateUserEvents(42, users, profiles, listings, categories)
	editEvents := GenerateListingEditEvents(100, listings)
	
	// Test 1: Deal should be unique per listing (one deal per listing maximum)
	listingDealCount := make(map[uuid.UUID]int)
	for _, deal := range deals {
		listingDealCount[deal.ListingID]++
	}
	
	for listingID, count := range listingDealCount {
		if count > 1 {
			t.Errorf("Listing %s has %d deals, expected at most 1", listingID, count)
		}
	}
	
	// Test 2: CancelDeal should occur before listing is closed
	for _, event := range activityEvents {
		if event.Type != model.EventType(model.EventTypeCancelDeal) || event.ListingID == nil {
			continue
		}
		
		listingIdx := -1
		for i, listing := range listings {
			if listing.ID == *event.ListingID {
				listingIdx = i
				break
			}
		}
		if listingIdx >= 0 {
			listing := listings[listingIdx]
			if listing.ClosedAt != nil && event.OccurredAt.After(*listing.ClosedAt) {
				t.Errorf("CancelDeal occurred after listing close for listing %s", listing.ID)
			}
		}
	}
	
	// Test 3: Edit should occur before listing is closed (explicit check)
	for _, event := range editEvents {
		if event.Type != model.EventType(model.EventTypeEdit) || event.ListingID == nil {
			continue
		}
		
		listingIdx := -1
		for i, listing := range listings {
			if listing.ID == *event.ListingID {
				listingIdx = i
				break
			}
		}
		if listingIdx >= 0 {
			listing := listings[listingIdx]
			if listing.ClosedAt != nil && event.OccurredAt.After(*listing.ClosedAt) {
				t.Errorf("Edit occurred after listing close for listing %s", listing.ID)
			}
		}
	}
	
	// Test 4: Edit should only occur for listings that are not closed at edit time
	for _, event := range editEvents {
		if event.Type != model.EventType(model.EventTypeEdit) || event.ListingID == nil {
			continue
		}
		
		listingIdx := -1
		for i, listing := range listings {
			if listing.ID == *event.ListingID {
				listingIdx = i
				break
			}
		}
		if listingIdx >= 0 {
			listing := listings[listingIdx]
			// The generator should not create edit events for closed listings
			if listing.ClosedAt != nil {
				t.Errorf("Edit event created for already closed listing %s", listing.ID)
			}
		}
	}
	
	// Test 5: CancelDeal should not lead to a Deal for the same listing BY THE SAME USER
	cancelledListingsPerUser := make(map[uuid.UUID]map[uuid.UUID]bool) // userID -> listingID -> cancelled
	for _, event := range activityEvents {
		if event.Type == model.EventType(model.EventTypeCancelDeal) && event.ListingID != nil {
			if cancelledListingsPerUser[event.UserID] == nil {
				cancelledListingsPerUser[event.UserID] = make(map[uuid.UUID]bool)
			}
			cancelledListingsPerUser[event.UserID][*event.ListingID] = true
		}
	}
	
	for _, deal := range deals {
		if cancelledListingsPerUser[deal.BuyerID] != nil && cancelledListingsPerUser[deal.BuyerID][deal.ListingID] {
			t.Errorf("Deal created for listing %s by user %s that had a CancelDeal event", deal.ListingID, deal.BuyerID)
		}
	}
	
	// Test 6: Reviews should be unique per deal
	dealReviewCount := make(map[uuid.UUID]int)
	for _, review := range reviews {
		if review.DealID != nil {
			dealReviewCount[*review.DealID]++
		}
	}
	
	for dealID, count := range dealReviewCount {
		if count > 1 {
			t.Errorf("Deal %s has %d reviews, expected at most 1", dealID, count)
		}
	}
	
	t.Logf("Passed business logic tests with %d users, %d listings, %d events, %d deals, %d reviews, %d edit events",
		len(users), len(listings), len(activityEvents), len(deals), len(reviews), len(editEvents))
}
