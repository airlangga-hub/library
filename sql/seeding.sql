-- 1. Insert Categories
INSERT INTO categories (name) VALUES
('Science Fiction'), ('Biography'), ('History'), ('Technology');

-- 2. Insert Users
INSERT INTO users (full_name, email, password, admin, author) VALUES
('Alice Author', 'alice@example.com', 'password', false, true),
('Bob Writer', 'bob@example.com', 'password', false, true),
('Charlie Renter', 'charlie@example.com', 'password', false, false),
('David Library', 'david@example.com', 'password', true, false);

-- 3. Insert Books
INSERT INTO books (title, description, author_id, category_id, available) VALUES
('Go Microservices', 'A guide to Go', 6, 4, false),
('Postgres Mastery', 'Deep dive into SQL', 6, 4, false),
('The Great Beyond', 'Space exploration', 7, 1, false),
('Ancient Rome', 'History of an Empire', 7, 3, false),
('Mystery of Code', 'Debugging life', 6, 4, false),
('Life of Pi-thon', 'Biography of a dev', 7, 2, false),
('Future Tech', 'AI and more', 6, 4, false),
('The Last Page', 'Fiction thriller', 7, 1, false);

-- 4. Insert Rents
INSERT INTO rents (book_id, user_id, created_at, due_date, return_date, fine, active) VALUES
-- User 1 Rents
(9, 2, NOW(), NOW() + INTERVAL '7 days', NULL, 0, true),
(10, 2, NOW(), NOW() + INTERVAL '7 days', NULL, 0, true),
(11, 2, NOW(), NOW() + INTERVAL '7 days', NULL, 0, true),
(12, 2, NOW(), NOW() + INTERVAL '7 days', NULL, 0, true),

-- User 2 Rents
(13, 5, NOW(), NOW() + INTERVAL '7 days', NULL, 0, true),
(14, 5, NOW(), NOW() + INTERVAL '7 days', NULL, 0, true),
(15, 5, NOW(), NOW() + INTERVAL '7 days', NULL, 0, true),
(16, 5, NOW(), NOW() + INTERVAL '7 days', NULL, 0, true),

-- User 3 Rents
(9, 8, NOW() - INTERVAL '14 days', NOW() - INTERVAL '7 days', NOW() - INTERVAL '6 days', 2000, false),
(10, 8, NOW() - INTERVAL '14 days', NOW() - INTERVAL '7 days', NOW() - INTERVAL '5 days', 4000, false),
(11, 8, NOW(), NOW() + INTERVAL '7 days', NULL, 0, true),
(12, 8, NOW(), NOW() + INTERVAL '7 days', NULL, 0, true),

-- User 4 Rents
(13, 9, NOW(), NOW() + INTERVAL '7 days', NULL, 0, true),
(14, 9, NOW(), NOW() + INTERVAL '7 days', NULL, 0, true),
(15, 9, NOW(), NOW() + INTERVAL '7 days', NULL, 0, true),
(16, 9, NOW(), NOW() + INTERVAL '7 days', NULL, 0, true);
