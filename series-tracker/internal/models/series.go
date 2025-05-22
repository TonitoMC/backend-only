package models

// Serie represents a series as stored in the database and as expected
// in JSON responses to the frontend.
type Serie struct {
	ID             int    `json:"id"`                 // Unique identifier for the series
	Title          string `json:"title"`              // Title of the series
	Ranking        int    `json:"ranking"`            // Score of the series used for ranking
	Status         string `json:"status"`             // Current status of the series; "Watching", "Plan to Watch", "Dropped", "Completed"
	CurrentEpisode int    `json:"lastEpisodeWatched"` // Last episode watched of the series
	TotalEpisodes  int    `json:"totalEpisodes"`      // Quantity of episodes in the series
}

// Status represents the payload for updating a series' status.
type Status struct {
	Status string `json:"status"` // Status of the series; "Watching", "Plan to Watch", "Dropped", "Completed"
}
