-- Создание таблицы ролей
CREATE TABLE roles (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(50) NOT NULL UNIQUE,
                       permissions TEXT[] NOT NULL DEFAULT '{}'
);

-- Создание таблицы файлов
CREATE TABLE files_s3 (
                          id SERIAL PRIMARY KEY,
                          file_url TEXT,
                          created_at TIMESTAMP NOT NULL DEFAULT NOW()
);


-- Создание таблицы пользователей
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       email VARCHAR(255) NOT NULL UNIQUE,
                       password_hash TEXT NOT NULL,
                       username TEXT NOT NULL UNIQUE,
                       avatar_id INTEGER NULL,
                       role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
                       created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
  ALTER TABLE users ADD CONSTRAINT fk_avatar_id FOREIGN KEY (avatar_id)
        REFERENCES files_s3(id)
        ON DELETE RESTRICT;

-- Создание таблицы видео-групп
CREATE TABLE video_groups (
                              id SERIAL PRIMARY KEY,
                              name TEXT NOT NULL,
                              type TEXT NOT NULL,
                              description TEXT,
                              cover_url TEXT
);

-- Создание таблицы видео
CREATE TABLE videos (
                        id SERIAL PRIMARY KEY,
                        group_id INTEGER NOT NULL REFERENCES video_groups(id) ON DELETE CASCADE,
                        title TEXT NOT NULL,
                        description TEXT,
                        duration INTEGER NOT NULL,
                        hls_url TEXT NOT NULL,
                        thumbnail_url TEXT NOT NULL
);

-- Создание таблицы комнат
CREATE TABLE rooms (
                       id SERIAL PRIMARY KEY,
                       owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                       current_video_id INTEGER REFERENCES videos(id) ON DELETE SET NULL,
                       name TEXT NOT NULL,
                       invite_code TEXT UNIQUE,
                       is_private BOOLEAN NOT NULL DEFAULT FALSE,
                       created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Создание таблицы состояний воспроизведения
CREATE TABLE playback_states (
                                 room_id INTEGER PRIMARY KEY REFERENCES rooms(id) ON DELETE CASCADE,
                                 status TEXT NOT NULL,
                                 current_times FLOAT NOT NULL DEFAULT 0,
                                 last_updated TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Создание таблицы участников комнат
CREATE TABLE room_members (
                              room_id INTEGER NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
                              user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                              joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
                              PRIMARY KEY (room_id, user_id)
);

-- Создание таблицы сообщений
CREATE TABLE messages (
                          id SERIAL PRIMARY KEY,
                          room_id INTEGER NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
                          user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                          text TEXT NOT NULL,
                          created_at TIMESTAMP NOT NULL DEFAULT NOW()
);


-- Индексы для ускорения поиска
CREATE INDEX idx_room_members_user ON room_members(user_id);
CREATE INDEX idx_messages_room ON messages(room_id);
CREATE INDEX idx_videos_group ON videos(group_id);
CREATE INDEX idx_file_url ON files_s3(file_url);