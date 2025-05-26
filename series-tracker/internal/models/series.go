package models

// Serie represents a series as stored in the database and as expected
// in JSON responses to the frontend.
type Serie struct {
	ID             int    `json:"id" example:"10"`                 // Unique identifier for the series
	Title          string `json:"title" example:"Breaking Bad"`    // Title of the series
	Ranking        int    `json:"ranking" example:"10"`            // Score of the series used for ranking
	Status         string `json:"status" example:"Watching"`       // Current status of the series; "Watching", "Plan to Watch", "Dropped", "Completed"
	CurrentEpisode int    `json:"lastEpisodeWatched" example:"15"` // Last episode watched of the series
	TotalEpisodes  int    `json:"totalEpisodes" example:"24"`      // Quantity of episodes in the series
}
