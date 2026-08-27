package database

import (
	"context"
	"fmt"
	"movies-api/internal/handlers"
	"movies-api/internal/service"
)

// resetDatabase deletes all data and reseeds database with dummy data. Used for testing.
func ResetDatabase(app *handlers.App) error {
	ctx := context.Background()

	query := `
	DROP TABLE genres_movies;
	DROP TABLE movies_actors;
	DROP TABLE genres;
	DROP TABLE actors;
	DROP TABLE movies;
	`

	// Delete all data
	_, err := app.Repo.DB.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	fmt.Println("Deleted all data from database")

	// Reset schema and recreate tables
	if err := InitDB(app.Repo.DB); err != nil {
		return err
	}
	fmt.Println("Reset database schema")

	// Repopulate tables with dummy data
	if err := seedDatabase(ctx, app); err != nil {
		return err
	}

	return nil
}

// seedDatabase populates database with dummy data
func seedDatabase(ctx context.Context, app *handlers.App) error {
	if err := seedGenres(ctx, app.GenreService); err != nil {
		return fmt.Errorf("seedGenres: %w", err)
	}
	if err := seedActors(ctx, app.ActorService); err != nil {
		return fmt.Errorf("seedActors: %w", err)
	}
	if err := seedMovies(ctx, app.MovieService); err != nil {
		return fmt.Errorf("seedMovies: %w", err)
	}
	return nil
}

// seedGenres seeds database with 10 dummy genres
func seedGenres(ctx context.Context, gs *service.GenreService) error {
	genres := []service.GenreSubmission{
		{Name: "Action"},
		{Name: "Adventure"},
		{Name: "Animation"},
		{Name: "Comedy"},
		{Name: "Drama"},
		{Name: "Horror"},
		{Name: "Romance"},
		{Name: "Science Fiction"},
		{Name: "Thriller"},
		{Name: "Fantasy"},
	}

	for _, g := range genres {
		_, err := gs.AddGenre(ctx, g)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Database seeded with %d dummy genres\n", len(genres))
	return nil
}

// seedActors seeds database with 25 dummy actors
func seedActors(ctx context.Context, as *service.ActorService) error {
	actors := []service.ActorSubmission{
		{
			Name:      "James Carter",
			BirthDate: "1985-03-14",
		},
		{
			Name:      "Emma Wilson",
			BirthDate: "1990-07-22",
		},
		{
			Name:      "Michael Brooks",
			BirthDate: "1978-11-05",
		},
		{
			Name:      "Sophia Bennett",
			BirthDate: "1993-02-18",
		},
		{
			Name:      "Daniel Foster",
			BirthDate: "1982-09-30",
		},
		{
			Name:      "Olivia Hayes",
			BirthDate: "1988-12-11",
		},
		{
			Name:      "Robert Mitchell",
			BirthDate: "1975-06-26",
		},
		{
			Name:      "Isabella Turner",
			BirthDate: "1995-04-09",
		},
		{
			Name:      "William Parker",
			BirthDate: "1980-01-17",
		},
		{
			Name:      "Ava Richardson",
			BirthDate: "1991-10-03",
		},
		{
			Name:      "Christopher Morgan",
			BirthDate: "1986-05-28",
		},
		{
			Name:      "Mia Cooper",
			BirthDate: "1994-08-15",
		},
		{
			Name:      "Matthew Edwards",
			BirthDate: "1979-03-07",
		},
		{
			Name:      "Charlotte Reed",
			BirthDate: "1987-11-21",
		},
		{
			Name:      "Anthony Collins",
			BirthDate: "1983-06-12",
		},
		{
			Name:      "Amelia Stewart",
			BirthDate: "1992-01-29",
		},
		{
			Name:      "Joseph Ward",
			BirthDate: "1976-10-16",
		},
		{
			Name:      "Harper Murphy",
			BirthDate: "1996-05-04",
		},
		{
			Name:      "Andrew Bailey",
			BirthDate: "1981-08-23",
		},
		{
			Name:      "Evelyn Rivera",
			BirthDate: "1989-12-06",
		},
		{
			Name:      "Joshua Cooper",
			BirthDate: "1977-04-19",
		},
		{
			Name:      "Grace Peterson",
			BirthDate: "1993-09-27",
		},
		{
			Name:      "Ryan Hughes",
			BirthDate: "1984-02-02",
		},
		{
			Name:      "Lily Simmons",
			BirthDate: "1997-07-13",
		},
		{
			Name:      "Nathan Price",
			BirthDate: "1980-11-24",
		},
	}

	for _, a := range actors {
		_, err := as.AddActor(ctx, a)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Database seeded with %d dummy actors\n", len(actors))
	return nil
}

// seedMovies seeds database with 100 dummy movies
func seedMovies(ctx context.Context, ms *service.MovieService) error {
	movies := []service.MovieSubmission{
		{Title: "The Last Horizon", ReleaseYear: 1995, Duration: 118, GenreIDs: []int64{1}, ActorIDs: []int64{1}},
		{Title: "Midnight Runaway", ReleaseYear: 1996, Duration: 104, GenreIDs: []int64{2, 5, 9}, ActorIDs: []int64{2, 14}},
		{Title: "Echoes of Tomorrow", ReleaseYear: 1997, Duration: 126, GenreIDs: []int64{1, 4, 7, 9}, ActorIDs: []int64{3, 9, 18}},
		{Title: "Summer at Dawn", ReleaseYear: 1998, Duration: 97, GenreIDs: []int64{3}, ActorIDs: []int64{4, 11}},
		{Title: "Shadow Protocol", ReleaseYear: 1999, Duration: 132, GenreIDs: []int64{2, 7, 10}, ActorIDs: []int64{5, 16, 21, 6}},
		{Title: "The Forgotten Road", ReleaseYear: 2000, Duration: 111, GenreIDs: []int64{4, 8, 10}, ActorIDs: []int64{6}},
		{Title: "City of Glass", ReleaseYear: 2001, Duration: 121, GenreIDs: []int64{1, 5, 6}, ActorIDs: []int64{7, 20, 3, 12}},
		{Title: "Beyond the Storm", ReleaseYear: 2002, Duration: 109, GenreIDs: []int64{3, 7}, ActorIDs: []int64{8, 15}},
		{Title: "A Quiet Place in Time", ReleaseYear: 2003, Duration: 103, GenreIDs: []int64{4}, ActorIDs: []int64{9, 17, 22, 1, 13}},
		{Title: "Crimson Skies", ReleaseYear: 2004, Duration: 128, GenreIDs: []int64{2, 8}, ActorIDs: []int64{10}},
		{Title: "The Hidden Valley", ReleaseYear: 2005, Duration: 115, GenreIDs: []int64{1, 3}, ActorIDs: []int64{11, 23, 7}},
		{Title: "Broken Compass", ReleaseYear: 2006, Duration: 107, GenreIDs: []int64{5, 7, 8}, ActorIDs: []int64{12, 16}},
		{Title: "Letters from Paris", ReleaseYear: 2007, Duration: 101, GenreIDs: []int64{3}, ActorIDs: []int64{13, 18, 2, 20}},
		{Title: "Edge of Darkness", ReleaseYear: 2008, Duration: 137, GenreIDs: []int64{2, 4}, ActorIDs: []int64{14, 21, 24}},
		{Title: "The Silver Key", ReleaseYear: 2009, Duration: 113, GenreIDs: []int64{1, 8, 5}, ActorIDs: []int64{15, 20}},
		{Title: "Winter's End", ReleaseYear: 2010, Duration: 119, GenreIDs: []int64{3, 5}, ActorIDs: []int64{16, 22, 25}},
		{Title: "Neon Boulevard", ReleaseYear: 2011, Duration: 105, GenreIDs: []int64{2}, ActorIDs: []int64{17, 19, 4}},
		{Title: "The Long Way Home", ReleaseYear: 2012, Duration: 124, GenreIDs: []int64{4, 7, 8}, ActorIDs: []int64{18, 23, 9, 1, 5}},
		{Title: "Whispers in the Dark", ReleaseYear: 2013, Duration: 98, GenreIDs: []int64{5}, ActorIDs: []int64{19, 24}},
		{Title: "Fire on the Mountain", ReleaseYear: 2014, Duration: 130, GenreIDs: []int64{1, 7}, ActorIDs: []int64{20, 21, 6}},
		{Title: "The Painted House", ReleaseYear: 2015, Duration: 116, GenreIDs: []int64{3, 4, 6}, ActorIDs: []int64{1, 13}},
		{Title: "Ocean of Stars", ReleaseYear: 2016, Duration: 122, GenreIDs: []int64{1, 6}, ActorIDs: []int64{2, 8, 17, 20}},
		{Title: "Dead Reckoning", ReleaseYear: 2017, Duration: 134, GenreIDs: []int64{2, 7, 8}, ActorIDs: []int64{3, 15, 20, 21, 7}},
		{Title: "The Empty Room", ReleaseYear: 2018, Duration: 96, GenreIDs: []int64{4}, ActorIDs: []int64{4}},
		{Title: "Under Northern Lights", ReleaseYear: 2019, Duration: 110, GenreIDs: []int64{3, 8}, ActorIDs: []int64{5, 18, 23}},
		{Title: "The Final Witness", ReleaseYear: 2020, Duration: 129, GenreIDs: []int64{2, 4}, ActorIDs: []int64{6, 14, 16, 22}},
		{Title: "After the Rain", ReleaseYear: 2021, Duration: 102, GenreIDs: []int64{5, 6, 7}, ActorIDs: []int64{7, 19}},
		{Title: "Parallel Lives", ReleaseYear: 2022, Duration: 127, GenreIDs: []int64{1}, ActorIDs: []int64{8, 12, 18, 3, 10}},
		{Title: "The Midnight Garden", ReleaseYear: 2023, Duration: 108, GenreIDs: []int64{3, 7}, ActorIDs: []int64{9, 16, 22}},
		{Title: "Kingdom of Ash", ReleaseYear: 2024, Duration: 141, GenreIDs: []int64{2, 8, 4}, ActorIDs: []int64{10, 21}},
		{Title: "The Secret Passenger", ReleaseYear: 2025, Duration: 114, GenreIDs: []int64{1, 5}, ActorIDs: []int64{11, 15, 24, 2}},
		{Title: "Road to Yesterday", ReleaseYear: 2026, Duration: 99, GenreIDs: []int64{3}, ActorIDs: []int64{12, 17}},
		{Title: "Red River", ReleaseYear: 1995, Duration: 123, GenreIDs: []int64{2, 7}, ActorIDs: []int64{13, 20, 5}},
		{Title: "The Glass Tower", ReleaseYear: 1996, Duration: 117, GenreIDs: []int64{4, 8, 6}, ActorIDs: []int64{14, 22, 1}},
		{Title: "One Last Summer", ReleaseYear: 1997, Duration: 95, GenreIDs: []int64{3, 5}, ActorIDs: []int64{15}},
		{Title: "The Silent Witness", ReleaseYear: 1998, Duration: 131, GenreIDs: []int64{2}, ActorIDs: []int64{16, 21, 23, 9}},
		{Title: "Lost in Translation", ReleaseYear: 1999, Duration: 106, GenreIDs: []int64{1, 6, 8}, ActorIDs: []int64{17, 24}},
		{Title: "House of Shadows", ReleaseYear: 2000, Duration: 120, GenreIDs: []int64{5, 8}, ActorIDs: []int64{18, 20, 4, 11}},
		{Title: "Across the Divide", ReleaseYear: 2001, Duration: 112, GenreIDs: []int64{3, 7}, ActorIDs: []int64{19, 22}},
		{Title: "The Burning Sky", ReleaseYear: 2002, Duration: 136, GenreIDs: []int64{1, 2, 4}, ActorIDs: []int64{20, 23}},
		{Title: "A Place to Remember", ReleaseYear: 2003, Duration: 103, GenreIDs: []int64{4, 6}, ActorIDs: []int64{1, 9, 14, 18, 7}},
		{Title: "The Last Letter", ReleaseYear: 2004, Duration: 109, GenreIDs: []int64{3}, ActorIDs: []int64{2, 16}},
		{Title: "Into the Wild Blue", ReleaseYear: 2005, Duration: 128, GenreIDs: []int64{1, 7}, ActorIDs: []int64{3, 18, 21}},
		{Title: "The Darkest Hour", ReleaseYear: 2006, Duration: 142, GenreIDs: []int64{2, 8}, ActorIDs: []int64{4, 20, 6, 13}},
		{Title: "Garden of Secrets", ReleaseYear: 2007, Duration: 100, GenreIDs: []int64{4, 5, 6}, ActorIDs: []int64{5, 12, 24}},
		{Title: "The Open Road", ReleaseYear: 2008, Duration: 118, GenreIDs: []int64{3, 6}, ActorIDs: []int64{6, 17}},
		{Title: "City After Midnight", ReleaseYear: 2009, Duration: 125, GenreIDs: []int64{2, 7}, ActorIDs: []int64{7, 19, 23, 15}},
		{Title: "The Lost Kingdom", ReleaseYear: 2010, Duration: 138, GenreIDs: []int64{1, 8, 4}, ActorIDs: []int64{8, 21}},
		{Title: "Three Days in June", ReleaseYear: 2011, Duration: 97, GenreIDs: []int64{3, 4}, ActorIDs: []int64{9, 15, 2}},
		{Title: "The Iron Gate", ReleaseYear: 2012, Duration: 133, GenreIDs: []int64{2, 5}, ActorIDs: []int64{10, 18, 22}},
		{Title: "Blue Horizon", ReleaseYear: 2013, Duration: 115, GenreIDs: []int64{1, 6, 7}, ActorIDs: []int64{11, 20, 22, 3}},
		{Title: "The Last Frontier", ReleaseYear: 2014, Duration: 129, GenreIDs: []int64{4}, ActorIDs: []int64{12, 16}},
		{Title: "Falling Stars", ReleaseYear: 2015, Duration: 104, GenreIDs: []int64{3, 8}, ActorIDs: []int64{13, 19, 24, 5, 10}},
		{Title: "The Forgotten City", ReleaseYear: 2016, Duration: 121, GenreIDs: []int64{2, 4, 8}, ActorIDs: []int64{14, 23}},
		{Title: "Beyond the Wall", ReleaseYear: 2017, Duration: 126, GenreIDs: []int64{1, 5}, ActorIDs: []int64{15, 21}},
		{Title: "The River Runs North", ReleaseYear: 2018, Duration: 110, GenreIDs: []int64{3, 6}, ActorIDs: []int64{16, 20, 8, 12}},
		{Title: "Midnight Expressway", ReleaseYear: 2019, Duration: 119, GenreIDs: []int64{2}, ActorIDs: []int64{17, 22}},
		{Title: "The Hidden Signal", ReleaseYear: 2020, Duration: 107, GenreIDs: []int64{4, 8, 5}, ActorIDs: []int64{18, 24, 6}},
		{Title: "When We Were Young", ReleaseYear: 2021, Duration: 101, GenreIDs: []int64{3, 5}, ActorIDs: []int64{19, 21, 4}},
		{Title: "The Final Journey", ReleaseYear: 2022, Duration: 135, GenreIDs: []int64{1, 7}, ActorIDs: []int64{20, 23, 9, 14}},
		{Title: "The Crimson Door", ReleaseYear: 2023, Duration: 112, GenreIDs: []int64{2, 8}, ActorIDs: []int64{1, 6, 18}},
		{Title: "Autumn Memories", ReleaseYear: 2024, Duration: 98, GenreIDs: []int64{4, 6, 3}, ActorIDs: []int64{2, 13}},
		{Title: "The Great Escape", ReleaseYear: 2025, Duration: 124, GenreIDs: []int64{1, 5}, ActorIDs: []int64{3, 17, 22, 11, 19}},
		{Title: "Nightfall", ReleaseYear: 2026, Duration: 105, GenreIDs: []int64{2, 7}, ActorIDs: []int64{4, 20}},
		{Title: "The Endless Sea", ReleaseYear: 1995, Duration: 130, GenreIDs: []int64{3, 8, 6}, ActorIDs: []int64{5, 15, 24}},
		{Title: "A World Apart", ReleaseYear: 1996, Duration: 117, GenreIDs: []int64{1, 4}, ActorIDs: []int64{6, 19, 21, 2}},
		{Title: "The Broken Crown", ReleaseYear: 1997, Duration: 139, GenreIDs: []int64{2, 6}, ActorIDs: []int64{7, 21}},
		{Title: "Hidden in Plain Sight", ReleaseYear: 1998, Duration: 108, GenreIDs: []int64{5}, ActorIDs: []int64{8, 14, 16, 23}},
		{Title: "The Last Train Home", ReleaseYear: 1999, Duration: 116, GenreIDs: []int64{3, 5, 7}, ActorIDs: []int64{9, 16, 23}},
		{Title: "Golden Hour", ReleaseYear: 2000, Duration: 102, GenreIDs: []int64{1, 6}, ActorIDs: []int64{10, 18, 20, 5}},
		{Title: "The Black Forest", ReleaseYear: 2001, Duration: 127, GenreIDs: []int64{4, 8}, ActorIDs: []int64{11, 20, 24}},
		{Title: "Crossing Borders", ReleaseYear: 2002, Duration: 111, GenreIDs: []int64{2, 7, 1}, ActorIDs: []int64{12, 22}},
		{Title: "The Old Photograph", ReleaseYear: 2003, Duration: 94, GenreIDs: []int64{3, 5}, ActorIDs: []int64{13, 17}},
		{Title: "Storm Chaser", ReleaseYear: 2004, Duration: 123, GenreIDs: []int64{1, 7}, ActorIDs: []int64{14, 19, 6, 9}},
		{Title: "The Missing Piece", ReleaseYear: 2005, Duration: 106, GenreIDs: []int64{4, 6, 8}, ActorIDs: []int64{15, 21, 23}},
		{Title: "A Distant Memory", ReleaseYear: 2006, Duration: 114, GenreIDs: []int64{3}, ActorIDs: []int64{16, 24, 10}},
		{Title: "The Northern Star", ReleaseYear: 2007, Duration: 128, GenreIDs: []int64{2, 5}, ActorIDs: []int64{17, 20}},
		{Title: "Roads Less Traveled", ReleaseYear: 2008, Duration: 120, GenreIDs: []int64{1, 4, 7}, ActorIDs: []int64{18, 22, 3, 12}},
		{Title: "The Silent Lake", ReleaseYear: 2009, Duration: 100, GenreIDs: []int64{3, 6}, ActorIDs: []int64{19, 23}},
		{Title: "Empire of Dust", ReleaseYear: 2010, Duration: 137, GenreIDs: []int64{2, 8}, ActorIDs: []int64{20, 21, 24, 7}},
		{Title: "The Secret Garden", ReleaseYear: 2011, Duration: 109, GenreIDs: []int64{4, 5}, ActorIDs: []int64{1, 11}},
		{Title: "Chasing Shadows", ReleaseYear: 2012, Duration: 122, GenreIDs: []int64{2, 7, 8}, ActorIDs: []int64{2, 15, 19, 6}},
		{Title: "The Final Chapter", ReleaseYear: 2013, Duration: 131, GenreIDs: []int64{1, 8}, ActorIDs: []int64{3, 18}},
		{Title: "Wild Hearts", ReleaseYear: 2014, Duration: 103, GenreIDs: []int64{3, 6}, ActorIDs: []int64{4, 22}},
		{Title: "The Dark River", ReleaseYear: 2015, Duration: 126, GenreIDs: []int64{2, 4}, ActorIDs: []int64{5, 13, 20}},
		{Title: "A New Beginning", ReleaseYear: 2016, Duration: 99, GenreIDs: []int64{3, 5, 6}, ActorIDs: []int64{6, 17, 21, 10}},
		{Title: "The Last Defender", ReleaseYear: 2017, Duration: 140, GenreIDs: []int64{1, 7}, ActorIDs: []int64{7, 20}},
		{Title: "Beyond the Horizon", ReleaseYear: 2018, Duration: 118, GenreIDs: []int64{4, 8, 1}, ActorIDs: []int64{8, 16, 23}},
		{Title: "The Memory Keeper", ReleaseYear: 2019, Duration: 107, GenreIDs: []int64{3}, ActorIDs: []int64{9, 19, 14, 2}},
		{Title: "The Last Kingdom", ReleaseYear: 2020, Duration: 143, GenreIDs: []int64{2, 7}, ActorIDs: []int64{10, 21, 24, 18, 5}},
		{Title: "The Hidden Door", ReleaseYear: 2021, Duration: 101, GenreIDs: []int64{1, 5}, ActorIDs: []int64{11, 18}},
		{Title: "Fireside Stories", ReleaseYear: 2022, Duration: 96, GenreIDs: []int64{3, 4}, ActorIDs: []int64{12, 22, 16}},
		{Title: "The Broken Bridge", ReleaseYear: 2023, Duration: 115, GenreIDs: []int64{2, 6}, ActorIDs: []int64{13, 16, 7}},
		{Title: "Across the Ocean", ReleaseYear: 2024, Duration: 129, GenreIDs: []int64{1, 8, 5}, ActorIDs: []int64{14, 19}},
		{Title: "The Final Hour", ReleaseYear: 2025, Duration: 133, GenreIDs: []int64{2, 7}, ActorIDs: []int64{15, 20, 23, 9}},
		{Title: "Dreams of Yesterday", ReleaseYear: 2026, Duration: 104, GenreIDs: []int64{3, 5}, ActorIDs: []int64{16, 21}},
		{Title: "The Midnight Bell", ReleaseYear: 1995, Duration: 117, GenreIDs: []int64{4, 8}, ActorIDs: []int64{17, 24, 11, 3}},
		{Title: "The Forgotten Hero", ReleaseYear: 1996, Duration: 125, GenreIDs: []int64{1, 6, 7}, ActorIDs: []int64{18, 20}},
		{Title: "Into the Unknown", ReleaseYear: 1997, Duration: 138, GenreIDs: []int64{2}, ActorIDs: []int64{19, 22, 8}},
		{Title: "The Last Sunrise", ReleaseYear: 1998, Duration: 110, GenreIDs: []int64{3, 5, 8}, ActorIDs: []int64{20, 23, 24, 1, 15}},
	}

	for _, m := range movies {
		_, err := ms.AddMovie(ctx, m)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Database seeded with %d dummy movies\n", len(movies))
	return nil
}
