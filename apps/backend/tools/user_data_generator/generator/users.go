package generator

import (
	"fmt"
	"math/rand"

	"avito-recap/internal/model"
)

func GenerateUserData(seed int64, numUsers int, categories []CategoryConfig) []model.User {
	if numUsers <= 0 {
		return nil
	}

	source := rand.NewSource(seed)
	rnd := rand.New(source)

	profiles := DefaultUserConfig()
	regions := DefaultRegions()
	users := make([]model.User, 0, numUsers)
	for i := 0; i < numUsers; i++ {
		profile := profiles[i%len(profiles)]
		preferred := profile.PreferredCategories
		if len(preferred) == 0 {
			preferred, _ = buildCategoryPreferences(rnd, categories)
		}
		unlikely := profile.UnlikelyCategories
		if len(unlikely) == 0 {
			_, unlikely = buildCategoryPreferences(rnd, categories)
		}

		users = append(users, model.User{
			ID:                  int64(i + 1),
			Username:            fmt.Sprintf("user_%d", rnd.Int()),
			RegisterDate:        randomDateBetweenYearsAgo(rnd, 2, 1),
			Region:              regions[rnd.Intn(len(regions))],
			PriceSegment:        profile.PriceSegment,
			PreferredCategories: preferred,
			UnlikelyCategories:  unlikely,
		})
	}
	return users
}

func GenerateSellerUsers(seed int64, numSellers int, startID int64) []model.User {
	if numSellers <= 0 {
		return nil
	}

	source := rand.NewSource(seed)
	rnd := rand.New(source)

	sellers := make([]model.User, 0, numSellers)
	regions := DefaultRegions()
	for i := 0; i < numSellers; i++ {
		sellers = append(sellers, model.User{
			ID:           startID + int64(i),
			Username:     fmt.Sprintf("seller_%d", i+1),
			RegisterDate: randomDateBetweenYearsAgo(rnd, 3, 1),
			Region:       regions[rnd.Intn(len(regions))],
		})
	}

	return sellers
}
