package services

// createLocation will create a location in the database.
// func createLocation(ctx context.Context, dto Request) (models.Location, error, int) {
//
// 	c, err := maps.NewClient(maps.WithAPIKey(helpers.GetEnvVar("GOOGLE_KEY")))
// 	if err != nil {
// 		log.Fatalf("fatal error: %s", err)
// 	}
//
// 	q := fmt.Sprintf("%d %s %s %s %s USA", dto.StreetNumber, dto.Street, dto.State, dto.City, dto.PostalCode)
// 	r := &maps.TextSearchRequest{
// 		Query: q,
// 	}
//
// 	result, err := c.TextSearch(ctx, r)
// 	if err != nil {
// 		log.Printf("error obtaining lat long: %v", err)
// 		return models.Location{}, err, http.StatusInternalServerError
// 	}
//
// 	now := time.Time{}
// 	loc := models.Location{
// 		Latitude:         result.Results[0].Geometry.Location.Lat,
// 		Longitude:        result.Results[0].Geometry.Location.Lng,
// 		Info:             "",
// 		Street:           dto.Street,
// 		Number:           dto.StreetNumber,
// 		City:             dto.City,
// 		State:            dto.State,
// 		PostalCode:       dto.PostalCode,
// 		FormattedAddress: result.Results[0].FormattedAddress,
// 		CreatedAt:        now,
// 		UpdatedAt:        &now,
// 	}
//
// 	location, err := db.PostLocation(ctx, loc)
// 	if err != nil {
// 		return models.Location{}, err, http.StatusInternalServerError
// 	}
//
// 	return location, nil, http.StatusOK
// }
