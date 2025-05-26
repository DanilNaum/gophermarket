-- Active: 1723551594755@@127.0.0.1@5432@db
ALTER  TABLE  withdrawal 
ADD COLUMN  id  UUID NOT NULL DEFAULT gen_random_uuid();