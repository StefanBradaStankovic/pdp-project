/*
Add soft delete property to actors table
*/

ALTER TABLE actors
ADD COLUMN "is_visible" BOOLEAN NOT NULL DEFAULT TRUE;

/*
Add soft delete property to directors table
*/

ALTER TABLE directors
ADD COLUMN "is_visible" BOOLEAN NOT NULL DEFAULT TRUE;

/*
Add soft delete property to movies table
*/

ALTER TABLE movies
ADD COLUMN "is_visible" BOOLEAN NOT NULL DEFAULT TRUE;

/*
Add soft delete property to movie_revenues table
*/

ALTER TABLE movie_revenues
ADD COLUMN "is_visible" BOOLEAN NOT NULL DEFAULT TRUE;