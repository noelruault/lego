DROP TABLE colors;
CREATE TABLE colors (
	id serial primary key,
	legoid text,
	name text,
	hex text,
	r int,
	g int,
	b int
);

COPY colors FROM '/Users/noelruault/go/src/github.com/noelruault/lego-project/colors-db.csv' WITH (
	FORMAT csv,
	DELIMITER ',',
	HEADER true
);


ALTER TABLE colors
	ADD COLUMN id serial primary key,
	ADD COLUMN r int,
	ADD COLUMN g int,
	ADD COLUMN b int,
	ADD COLUMN hexmatch text;

DROP TABLE seen;
CREATE TABLE seen (
	hex            text,
	minlegoid int REFERENCES colors(id),
	hexmindistance text,
	maxlegoid int REFERENCES colors(id),
	hexmaxdistance text
);
CREATE UNIQUE INDEX seen_hex_key ON seen(hex text_ops);

DROP TABLE artboards;
CREATE TABLE artboards (
	id serial primary key,
	imagename text,
	pixels text
);
