package generator

import (
	"encoding/binary"
	"fmt"
	"math/rand"

	"github.com/google/uuid"

	"avito-recap/internal/model"
)

// seededUUID generates a deterministic UUID based on the random seed
func seededUUID(rnd *rand.Rand) uuid.UUID {
	var bytes [16]byte
	binary.BigEndian.PutUint64(bytes[:8], rnd.Uint64())
	binary.BigEndian.PutUint64(bytes[8:], rnd.Uint64())
	
	// Set version (4) and variant bits
	bytes[6] = (bytes[6] & 0x0f) | 0x40 // Version 4
	bytes[8] = (bytes[8] & 0x3f) | 0x80 // Variant 1
	
	return uuid.UUID(bytes)
}

func GenerateUserData(seed int64, numUsers int, categories []CategoryConfig) []model.User {
	if numUsers <= 0 {
		return nil
	}

	source := rand.NewSource(seed)
	rnd := rand.New(source)

	config := DefaultGeneratorConfig()
	regions := DefaultRegions()
	users := make([]model.User, 0, numUsers)
	for i := 0; i < numUsers; i++ {
		users = append(users, model.User{
			ID:            seededUUID(rnd),
			DisplayName:   fmt.Sprintf("user_%d", i+1),
			RegisteredAt:  randomDateBetweenYearsAgo(rnd, config.UserRegistrationYearsAgo, 1),
			Region:        regions[rnd.Intn(len(regions))],
			Timezone:      "UTC",
			IsTestProfile: true,
		})
	}
	return users
}

func SelectSellerUsers(seed int64, users []model.User) []model.User {
	source := rand.NewSource(seed)
	rnd := rand.New(source)
	config := DefaultGeneratorConfig()

	sellers := make([]model.User, 0, len(users))
	for _, user := range users {
		if rnd.Float64() < config.SellerSelectionRate {
			sellers = append(sellers, user)
		}
	}

	if len(sellers) == 0 && len(users) > 0 {
		sellers = append(sellers, users[rnd.Intn(len(users))])
	}

	return sellers
}

func GenerateSellerUsers(seed int64, numSellers int, startID int64) []model.User {
	if numSellers <= 0 {
		return nil
	}

	source := rand.NewSource(seed)
	rnd := rand.New(source)
	config := DefaultGeneratorConfig()

	sellers := make([]model.User, 0, numSellers)
	regions := DefaultRegions()
	for i := 0; i < numSellers; i++ {
		sellers = append(sellers, model.User{
			ID:            seededUUID(rnd),
			DisplayName:   fmt.Sprintf("seller_%d", i+1),
			RegisteredAt:  randomDateBetweenYearsAgo(rnd, config.SellerRegistrationYearsAgo, 1),
			Region:        regions[rnd.Intn(len(regions))],
			Timezone:      "UTC",
			IsTestProfile: true,
		})
	}

	return sellers
}
