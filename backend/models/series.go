package models

// Represents a series as it's represented in the Database &
// expected JSON by the frontend
type Serie struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	Ranking        int    `json:"ranking"`
	Status         string `json:"status"`
	CurrentEpisode int    `json:"lastEpisodeWatched"`
	TotalEpisodes  int    `json:"totalEpisodes"`
}

// Represents the payload when updating series status
type Status struct {
	Status string `json:"status"`
}
