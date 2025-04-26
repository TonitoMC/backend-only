CREATE TABLE IF NOT EXISTS series (
  id SERIAL PRIMARY KEY,
  title VARCHAR UNIQUE NOT NULL,
  ranking INTEGER NOT NULL CHECK (ranking >= 0),
  status VARCHAR NOT NULL CHECK (status IN ('Watching', 'Plan to Watch', 'Dropped', 'Completed')),
  current_episode INTEGER NOT NULL,
  total_episodes INTEGER NOT NULL
);


INSERT INTO series (title, ranking, status, current_episode, total_episodes)
VALUES 
('Fullmetal Alchemist: Brotherhood', 10, 'Completed', 64, 64),
('One Piece', 9, 'Watching', 1100, 1100),
('Death Note', 8, 'Completed', 37, 37),
('Vinland Saga', 8, 'Plan to Watch', 0, 48),
('Bleach: Thousand-Year Blood War', 7, 'Watching', 26, 52);

SELECT * FROM series;
