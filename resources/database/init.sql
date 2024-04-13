CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY,
    date CHAR(8) NOT NULL,
    title VARCHAR (512) NOT NULL,
    comment VARCHAR (1024) NOT NULL DEFAULT "",
    repeat VARCHAR (128) NOT NULL
);

CREATE INDEX scheduler_date_idx ON scheduler(date);