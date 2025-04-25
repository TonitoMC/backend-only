package models

// Represents a series as it's represented in the Database
type Serie struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	Ranking        int    `json:"ranking"`
	Status         string `json:"status"`
	CurrentEpisode string `json:"current_episode"`
	TotalEpisodes  string `json:"total_episodes"`
}
