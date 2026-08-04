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

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pool, err := gen.ConnectToDatabase(ctx, dsn)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer pool.Close()

	categories := gen.DefaultCategories()
	profiles := gen.DefaultUserConfig()
	users := gen.GenerateUserData(seed, userCount, categories)
	sellers := gen.GenerateSellerUsers(seed, gen.SellerCountForListings(listingCount), int64(len(users)+1))
	listings := gen.GenerateListingData(seed, listingCount, categories, sellers, time.Now().UTC().AddDate(0, 0, 30))
	searches, views, favorites, contacts, deals, reviews, userEvents := gen.GenerateUserEvents(seed+2, users, profiles, listings, categories)

	allUsers := append(users, sellers...)
	if err := gen.InsertUsers(ctx, pool, allUsers); err != nil {
		log.Fatalf("insert users: %v", err)
	}
	if err := gen.InsertListings(ctx, pool, listings); err != nil {
		log.Fatalf("insert listings: %v", err)
	}
	if err := gen.InsertSearches(ctx, pool, searches); err != nil {
		log.Fatalf("insert searches: %v", err)
	}
	if err := gen.InsertViews(ctx, pool, views); err != nil {
		log.Fatalf("insert views: %v", err)
	}
	if err := gen.InsertFavorites(ctx, pool, favorites); err != nil {
		log.Fatalf("insert favorites: %v", err)
	}
	if err := gen.InsertContacts(ctx, pool, contacts); err != nil {
		log.Fatalf("insert contacts: %v", err)
	}
	if err := gen.InsertDeals(ctx, pool, deals); err != nil {
		log.Fatalf("insert deals: %v", err)
	}
	if err := gen.InsertReviews(ctx, pool, reviews); err != nil {
		log.Fatalf("insert reviews: %v", err)
	}
	if err := gen.InsertUserEvents(ctx, pool, userEvents); err != nil {
		log.Fatalf("insert user events: %v", err)
	}

	fmt.Printf("Inserted %d users (%d test + %d sellers), %d listings, %d searches, %d views, %d favorites, %d contacts, %d deals, %d reviews, %d user events\n",
		len(allUsers), len(users), len(sellers), len(listings), len(searches), len(views), len(favorites), len(contacts), len(deals), len(reviews), len(userEvents))
}
