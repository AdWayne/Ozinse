-- ==========================
-- ROLES
-- ==========================

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    permissions JSONB NOT NULL DEFAULT '{}'
);

-- ==========================
-- USERS
-- ==========================

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(150) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    user_avatar_url VARCHAR(255),
    full_name VARCHAR(100),
    phone VARCHAR(20),
    birth_date VARCHAR(100),
    role_id INT REFERENCES roles (id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- ==========================
-- CATEGORIES
-- ==========================

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

-- ==========================
-- AGE RATINGS
-- ==========================

CREATE TABLE age_ratings (
    id SERIAL PRIMARY KEY,
    range VARCHAR(50) NOT NULL UNIQUE
);

-- ==========================
-- PROJECTS
-- ==========================

CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    release_year INT CHECK (
        release_year IS NULL
        OR release_year >= 1800
    ),
    director VARCHAR(255),
    producer VARCHAR(255),
    cover_image_url VARCHAR(255),
    banner_image_url VARCHAR(255),
    keywords TEXT,
    category_id INT REFERENCES categories (id) ON DELETE SET NULL,
    age_rating_id INT REFERENCES age_ratings (id) ON DELETE SET NULL,
    project_type VARCHAR(20) NOT NULL CHECK (
        project_type IN ('MOVIE', 'SERIES')
    ),
    duration_minutes INT CHECK (
        duration_minutes IS NULL
        OR duration_minutes > 0
    ),
    youtube_video_id VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- ==========================
-- GENRES
-- ==========================

CREATE TABLE genres (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE project_genres (
    project_id INT NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    genre_id INT NOT NULL REFERENCES genres (id) ON DELETE CASCADE,
    PRIMARY KEY (project_id, genre_id)
);

-- ==========================
-- SEASONS
-- ==========================

CREATE TABLE seasons (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    season_number INT NOT NULL CHECK (season_number > 0),
    UNIQUE (project_id, season_number)
);

-- ==========================
-- EPISODES
-- ==========================

CREATE TABLE episodes (
    id SERIAL PRIMARY KEY,
    season_id INT NOT NULL REFERENCES seasons (id) ON DELETE CASCADE,
    episode_number INT NOT NULL CHECK (episode_number > 0),
    title VARCHAR(255),
    youtube_video_id VARCHAR(50) NOT NULL,
    duration INT CHECK (
        duration IS NULL
        OR duration > 0
    ),
    UNIQUE (season_id, episode_number)
);

-- ==========================
-- FAVORITES
-- ==========================

CREATE TABLE favorites (
    user_id INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    project_id INT NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    added_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, project_id)
);

-- ==========================
-- FEATURED CONTENT
-- ==========================

CREATE TABLE featured_content (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    block_type VARCHAR(50) NOT NULL,
    sort_order INT NOT NULL CHECK (sort_order > 0),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (block_type, sort_order)
);

-- ==========================
-- REFRESH TOKENS
-- ==========================

CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- ==========================
-- WATCH HISTORY
-- ==========================

CREATE TABLE watch_history (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    project_id INT REFERENCES projects (id) ON DELETE CASCADE,
    episode_id INT REFERENCES episodes (id) ON DELETE CASCADE,
    watched_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    progress_seconds INT DEFAULT 0 CHECK (progress_seconds >= 0),
    CHECK (
        (
            project_id IS NOT NULL
            AND episode_id IS NULL
        )
        OR (
            project_id IS NULL
            AND episode_id IS NOT NULL
        )
    )
);

