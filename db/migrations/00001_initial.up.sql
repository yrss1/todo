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
INSERT INTO users (name, email, password) VALUES
                                              ('Rick Sanchez', 'WubbaLubbadubdub@c137.com', '$2a$10$W7e.RujTp0v6M8/iGwXdr.SX6GTEm1Rdjptg2vMGUiPTK96pp09q6'),
                                              ('Morty Smith', 'theonetruemorty@c137.com', '$2a$10$DmiAv1tMmT9ZXZNB2Gp7yea8eT2KvNyRhhA.TBEMrFzW5fg8XybCG'),
                                              ('admin', 'admin@admin.com', '$2a$10$WHOiYYS7SHDnDQ2nVFyDnuUCq69PX4iHH9Se9VO5wAVF2SM/0h0jK');
-- password123, password456, admin,

-- Вставка задач для пользователей
INSERT INTO tasks (user_id, title, description, status) VALUES
                                                            ((SELECT id FROM users WHERE email = 'WubbaLubbadubdub@c137.com'), 'Buy groceries', 'Buy milk, bread, and eggs', 'active'),
                                                            ((SELECT id FROM users WHERE email = 'WubbaLubbadubdub@c137.com'), 'Walk the dog', 'Walk the dog in the park', 'active'),
                                                            ((SELECT id FROM users WHERE email = 'WubbaLubbadubdub@c137.com'), 'Read a book', 'Read at least 30 pages of a book', 'done'),
                                                            ((SELECT id FROM users WHERE email = 'WubbaLubbadubdub@c137.com'), 'Call mom', 'Call mom and check on her', 'active'),
                                                            ((SELECT id FROM users WHERE email = 'WubbaLubbadubdub@c137.com'), 'Pay bills', 'Pay utility bills online', 'active'),

                                                            ((SELECT id FROM users WHERE email = 'theonetruemorty@c137.com'), 'Complete project', 'Finish the project report', 'active'),
                                                            ((SELECT id FROM users WHERE email = 'theonetruemorty@c137.com'), 'Buy flight tickets', 'Purchase tickets for upcoming trip', 'done'),
                                                            ((SELECT id FROM users WHERE email = 'theonetruemorty@c137.com'), 'Update resume', 'Revise and update the resume', 'active'),
                                                            ((SELECT id FROM users WHERE email = 'theonetruemorty@c137.com'), 'Book doctor appointment', 'Schedule an appointment with the doctor', 'active'),
                                                            ((SELECT id FROM users WHERE email = 'theonetruemorty@c137.com'), 'Clean the house', 'Clean the entire house', 'done'),

                                                            ((SELECT id FROM users WHERE email = 'admin@admin.com'), 'Organize files', 'Sort and organize files on computer', 'active'),
                                                            ((SELECT id FROM users WHERE email = 'admin@admin.com'), 'Exercise', 'Go for a run or workout at the gym', 'done'),
                                                            ((SELECT id FROM users WHERE email = 'admin@admin.com'), 'Buy a gift', 'Purchase a gift for a friend’s birthday', 'active'),
                                                            ((SELECT id FROM users WHERE email = 'admin@admin.com'), 'Learn a new skill', 'Spend time learning a new skill online', 'active'),
                                                            ((SELECT id FROM users WHERE email = 'admin@admin.com'), 'Cook dinner', 'Prepare a meal for the evening', 'done'),

                                                            ((SELECT id FROM users WHERE email = 'WubbaLubbadubdub@c137.com'), 'Finish reading the book', 'Complete reading the current book', 'active'),
                                                            ((SELECT id FROM users WHERE email = 'WubbaLubbadubdub@c137.com'), 'Visit the gym', 'Attend a fitness class at the gym', 'done'),
                                                            ((SELECT id FROM users WHERE email = 'WubbaLubbadubdub@c137.com'), 'Buy new shoes', 'Purchase new running shoes', 'active'),
                                                            ((SELECT id FROM users WHERE email = 'theonetruemorty@c137.com'), 'Prepare presentation', 'Work on the presentation for the meeting', 'done'),
                                                            ((SELECT id FROM users WHERE email = 'theonetruemorty@c137.com'), 'Call service provider', 'Contact the service provider for issues', 'active'),

                                                            ((SELECT id FROM users WHERE email = 'admin@admin.com'), 'Plan weekend trip', 'Organize details for a weekend trip', 'done'),
                                                            ((SELECT id FROM users WHERE email = 'admin@admin.com'), 'Write blog post', 'Draft and publish a new blog post', 'active'),
                                                            ((SELECT id FROM users WHERE email = 'WubbaLubbadubdub@c137.com'), 'Attend webinar', 'Participate in an online webinar', 'done'),
                                                            ((SELECT id FROM users WHERE email = 'theonetruemorty@c137.com'), 'Review reports', 'Go through the latest project reports', 'active'),
                                                            ((SELECT id FROM users WHERE email = 'admin@admin.com'), 'Check email', 'Respond to important emails', 'done'),

                                                            ((SELECT id FROM users WHERE email = 'WubbaLubbadubdub@c137.com'), 'Update social media', 'Post updates on social media accounts', 'active'),
                                                            ((SELECT id FROM users WHERE email = 'WubbaLubbadubdub@c137.com'), 'Organize garage', 'Clean and organize the garage', 'done'),
                                                            ((SELECT id FROM users WHERE email = 'theonetruemorty@c137.com'), 'Check bank account', 'Review the latest bank transactions', 'active'),
                                                            ((SELECT id FROM users WHERE email = 'admin@admin.com'), 'Study for exam', 'Prepare for the upcoming exam', 'done'),
                                                            ((SELECT id FROM users WHERE email = 'WubbaLubbadubdub@c137.com'), 'Prepare for meeting', 'Get ready for the upcoming work meeting', 'active'),

                                                            ((SELECT id FROM users WHERE email = 'theonetruemorty@c137.com'), 'Schedule car maintenance', 'Book a service appointment for the car', 'done'),
                                                            ((SELECT id FROM users WHERE email = 'admin@admin.com'), 'Decorate home', 'Put up decorations for the season', 'active'),
                                                            ((SELECT id FROM users WHERE email = 'WubbaLubbadubdub@c137.com'), 'Research new technology', 'Look into the latest tech trends', 'done'),
                                                            ((SELECT id FROM users WHERE email = 'theonetruemorty@c137.com'), 'Finish the article', 'Complete writing the article', 'active'),
                                                            ((SELECT id FROM users WHERE email = 'admin@admin.com'), 'Grocery shopping', 'Buy groceries for the week', 'done');
