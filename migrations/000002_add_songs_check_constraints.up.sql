ALTER TABLE songs ADD CONSTRAINT songs_length_check CHECK (length >= 0);
ALTER TABLE songs ADD CONSTRAINT songs_year_check CHECK (year BETWEEN 1888 AND date_part('year', now()));
ALTER TABLE songs ADD CONSTRAINT genres_length_check CHECK (array_length(genres, 1) BETWEEN 1 AND 5);
