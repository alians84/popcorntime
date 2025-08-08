-- Создание таблицы файлов
CREATE TABLE files_s3 (
                          id SERIAL PRIMARY KEY,
                          file_url TEXT,
                          created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

ALTER TABLE users
    ADD COLUMN avatar_url INTEGER NULL,
    ADD CONSTRAINT fk_avatar_url FOREIGN KEY (avatar_url)
        REFERENCES files_s3(id)
        ON DELETE RESTRICT;