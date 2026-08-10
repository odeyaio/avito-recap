package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	gen "avito-recap/tools/user_data_generator/generator"
)

func main() {
	var (
		dsn          string
		seed         int64
		userCount    int
		listingCount int
	)

	flag.StringVar(&dsn, "dsn", os.Getenv("DATABASE_URL"), "PostgreSQL DSN")
	flag.Int64Var(&seed, "seed", 42, "Random seed")
	flag.IntVar(&userCount, "users", 100, "Number of users to create")
	flag.IntVar(&listingCount, "listings", 200, "Number of listings to create")
	flag.Parse()

	if dsn == "" {
		log.Fatal("DATABASE_URL is empty")
	}

	if userCount <= 0 {
		log.Fatal("users must be greater than 0")
	}
	if listingCount <= 0 {
		log.Fatal("listings must be greater than 0")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	pool, err := gen.ConnectToDatabase(ctx, dsn)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer pool.Close()

	categories := gen.DefaultCategories()
	categoryRecords := gen.BuildCategories(categories)
	profiles := gen.DefaultUserConfig()
	users := gen.GenerateUserData(seed, userCount, categories)
	sellers := gen.SelectSellerUsers(seed, users)
	listings := gen.GenerateListingData(seed, listingCount, categories, sellers, time.Now().UTC())
	editEvents := gen.GenerateListingEditEvents(seed+100, listings)
	activityEvents, deals, reviews := gen.GenerateUserEvents(seed, users, profiles, listings, categories)

	allActivityEvents := append(activityEvents, editEvents...)

	if err := gen.InsertCategories(ctx, pool, categoryRecords); err != nil {
		log.Fatalf("insert categories: %v", err)
	}
	if err := gen.InsertUsers(ctx, pool, users); err != nil {
		log.Fatalf("insert users: %v", err)
	}
	if err := gen.InsertListings(ctx, pool, listings); err != nil {
		log.Fatalf("insert listings: %v", err)
	}
	if err := gen.InsertActivityEvents(ctx, pool, allActivityEvents); err != nil {
		log.Fatalf("insert activity events: %v", err)
	}
	if err := gen.InsertDeals(ctx, pool, deals); err != nil {
		log.Fatalf("insert deals: %v", err)
	}
	if err := gen.InsertReviews(ctx, pool, reviews); err != nil {
		log.Fatalf("insert reviews: %v", err)
	}

	fmt.Printf("Inserted %d users (%d sellers), %d listings, %d activity events, %d deals, %d reviews\n",
		len(users), len(sellers), len(listings), len(allActivityEvents), len(deals), len(reviews))

	// year coverage: how many events land in each calendar year, and for
	// how many distinct users. recap request for year Y will 422 with
	// "insufficient activity" for any user with zero events in that row.
	fmt.Println("\n=== Event Coverage By Year ===")
	eventsByYear := make(map[int]int)
	usersByYear := make(map[int]map[string]struct{})
	for _, event := range allActivityEvents {
		year := event.OccurredAt.Year()
		eventsByYear[year]++
		if usersByYear[year] == nil {
			usersByYear[year] = make(map[string]struct{})
		}
		usersByYear[year][event.UserID.String()] = struct{}{}
	}
	for year := time.Now().UTC().Year() - 3; year <= time.Now().UTC().Year(); year++ {
		fmt.Printf("  %d: %d events across %d/%d users\n", year, eventsByYear[year], len(usersByYear[year]), len(users))
	}

	// print info about the first 3 users with default configs
	fmt.Println("\n=== First 3 Users (Default Configs) ===")
	for i := 0; i < 3 && i < len(users); i++ {
		user := users[i]
		profile := profiles[i]

		eventCounts := make(map[string]int)
		intentCount := 0
		for _, event := range allActivityEvents {
			if event.UserID == user.ID {
				eventType := string(event.Type)
				eventCounts[eventType]++
				if eventType == "search" {
					intentCount++
				}
			}
		}

		userDeals := 0
		userReviews := 0
		for _, deal := range deals {
			if deal.BuyerID == user.ID {
				userDeals++
			}
		}
		for _, review := range reviews {
			if review.AuthorID == user.ID {
				userReviews++
			}
		}

		fmt.Printf("\nUser %d:\n", i+1)
		fmt.Printf("  ID: %s\n", user.ID)
		fmt.Printf("  Username: %s\n", user.DisplayName)
		fmt.Printf("  Region: %s\n", user.Region)
		fmt.Printf("  Price Segment: %s\n", profile.PriceSegment)
		fmt.Printf("  Generated Data:\n")
		fmt.Printf("    Intents: %d\n", intentCount)
		fmt.Printf("    Total Events: %d\n", len(eventCounts))
		for eventType, count := range eventCounts {
			fmt.Printf("    %s: %d\n", eventType, count)
		}
		fmt.Printf("    Deals: %d\n", userDeals)
		fmt.Printf("    Reviews: %d\n", userReviews)
		fmt.Printf("  Probabilities:\n")
		fmt.Printf("    Search: %.2f\n", profile.Search)
		fmt.Printf("    View: %.2f\n", profile.View)
		fmt.Printf("    Favorite: %.2f\n", profile.Favorite)
		fmt.Printf("    Contact: %.2f\n", profile.Contact)
		fmt.Printf("    Deal: %.2f\n", profile.Deal)
		fmt.Printf("    Review: %.2f\n", profile.Review)
		fmt.Printf("    Notification Open: %.2f\n", profile.NotificationOpen)
		fmt.Printf("    Sell Frequency: %.2f\n", profile.SellFrequency)
		fmt.Printf("  Preferred Categories: %v\n", profile.PreferredCategories)
		fmt.Printf("  Unlikely Categories: %v\n", profile.UnlikelyCategories)
		fmt.Printf("  Preferred Sell Categories: %v\n", profile.PreferredSellCategories)
		fmt.Printf("  Unlikely Sell Categories: %v\n", profile.UnlikelySellCategories)
	}
}
