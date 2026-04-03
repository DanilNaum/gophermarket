-- Active: 1723551594755@@127.0.0.1@5432@db
CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    password_salt TEXT NOT NULL,


    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);