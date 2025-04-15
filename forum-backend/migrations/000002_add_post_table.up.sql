CREATE TABLE posts(
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            author_id INTEGER NOT NULL,
            title TEXT NOT NULL,
            content TEXT NOT NULL,
            creation_date TEXT NOT NULL,
            FOREIGN KEY (author_id) REFERENCES users(id)
);