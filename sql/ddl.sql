CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    admin BOOLEAN NOT NULL DEFAULT false,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    author BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    author_id INT REFERENCES users(id) ON DELETE SET NULL,
    category_id INT REFERENCES categories(id) ON DELETE SET NULL,
    available BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE IF NOT EXISTS rents (
    id SERIAL PRIMARY KEY,
    book_id INT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    due_date TIMESTAMPTZ NOT NULL,
    return_date TIMESTAMPTZ,
    fine INT NOT NULL DEFAULT 0,
    active BOOLEAN DEFAULT true
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_books_author_id ON books(author_id);
CREATE INDEX idx_books_category_id ON books(category_id);
CREATE INDEX idx_rents_book_id ON rents(book_id);
CREATE INDEX idx_rents_user_id ON rents(user_id);