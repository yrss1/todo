-- EXTENSIONS --
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- TABLES --
CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                     name VARCHAR(100) NOT NULL,
                                     email VARCHAR(100) UNIQUE NOT NULL,
                                     password VARCHAR(255) NOT NULL,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tasks (
                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                     user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                                     title VARCHAR(255) NOT NULL,
                                     description TEXT,
                                     status VARCHAR(10) DEFAULT 'active',
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- DATA --
-- INSERT INTO users (name, email, password) VALUES
--                                               ('Rick Sanchez', 'WubbaLubbadubdub@c137.com', 'password123'),
--                                               ('Morty Smith', 'theonetruemorty@c137.com', 'password456'),
--                                               ('Who u are', 'whouare@whoami.com', 'password789');
--
-- INSERT INTO tasks (user_id, title, description, status) VALUES
--                                                             ((SELECT id FROM users WHERE email = 'WubbaLubbadubdub@c137.com'), 'Buy groceries', 'Buy milk, bread, and eggs', 'active'),
--                                                             ((SELECT id FROM users WHERE email = 'theonetruemorty@c137.com'), 'Complete project', 'Finish the project report', 'active'),
--                                                             ((SELECT id FROM users WHERE email = 'theonetruemorty@c137.com'), 'Book flight tickets', 'Book tickets for the vacation', 'active');
