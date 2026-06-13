-- ==========================
-- ROLES
-- ==========================

INSERT INTO roles (name, permissions) VALUES
('admin', '{"all": true}'),
('user', '{"read": true}')
ON CONFLICT (name) DO NOTHING;

-- ==========================
-- ADMIN USER
-- ==========================

INSERT INTO users (email, password_hash, full_name, user_avatar_url, role_id)
VALUES (
    'admin@ozinse.com',
    '$2b$12$KZZjyWFPQv5JWo0Uj0RzUe.9Nf0IYK/7cEa8Abn2U6y/QqCqaz3dq',
    'Admin User',
    'http://localhost:8080/static/avatars/Avatar.svg',
    1
)
ON CONFLICT (email) DO NOTHING;

-- ==========================
-- CATEGORIES
-- ==========================

INSERT INTO categories (name) VALUES
('Телехикая'),
('Мультфильм'),
('Көркем фильм'),
('Деректі фильм'),
('Тв-бағдарлама және реалити-шоу'),
('Ситком'),
('Аниме'),
('Шетел фильмдері')
ON CONFLICT (name) DO NOTHING;

-- ==========================
-- GENRES
-- ==========================

INSERT INTO genres (name) VALUES
('Комедиялар'),
('Отбасымен көретіндер'),
('Ғылыми-танымдық'),
('Ойын-сауық'),
('Ғылыми фантастика және фэнтези'),
('Шытырман оқиғалы'),
('Қысқаметрлі'),
('Музыкалық'),
('Спорттық')
ON CONFLICT (name) DO NOTHING;

-- ==========================
-- AGE RATINGS
-- ==========================

INSERT INTO age_ratings (range) VALUES
('8-10 жас'),
('10-12 жас'),
('12-14 жас'),
('14-16 жас'),
('16-18 жас')
ON CONFLICT (range) DO NOTHING;

-- ==========================
-- MOVIE (id=1)
-- ==========================

INSERT INTO projects (
    id, title, description, release_year, director, producer,
    cover_image_url, banner_image_url, keywords,
    category_id, age_rating_id, project_type,
    duration_minutes, youtube_video_id
) VALUES (
    1,
    'Ғарышқа саяхат',
    'Адамзаттың ғарышқа алғашқы сапары және жаңа планеталарды зерттеуі туралы қызықты көркем фильм.',
    2023,
    'Асқар Үсенов',
    'Берік Қалиев',
    'https://example.com/images/space_movie_cover.jpg',
    'https://example.com/images/space_movie_banner.jpg',
    'ғарыш, саяхат, фантастика',
    3, 3, 'MOVIE',
    140, 'dQw4w9WgXcQ'
)
ON CONFLICT (id) DO NOTHING;

-- ==========================
-- SERIES (id=2)
-- ==========================

INSERT INTO projects (
    id, title, description, release_year, director, producer,
    cover_image_url, banner_image_url, keywords,
    category_id, age_rating_id, project_type,
    duration_minutes
) VALUES (
    2,
    'Ауылдастар',
    'Қаладан ауылға көшіп келген жастардың қызықты да күлкілі оқиғалары туралы ситком.',
    2022,
    'Мақсат Оспанов',
    'Айжан Серікқызы',
    'https://example.com/images/auyldastar_cover.jpg',
    'https://example.com/images/auyldastar_banner.jpg',
    'ауыл, комедия, ситком',
    6, 4, 'SERIES',
    30
)
ON CONFLICT (id) DO NOTHING;

-- ==========================
-- GENRES
-- ==========================

INSERT INTO project_genres (project_id, genre_id) VALUES
(1, 5),
(2, 1),
(2, 2)
ON CONFLICT (project_id, genre_id) DO NOTHING;

-- ==========================
-- SEASON 1
-- ==========================

INSERT INTO seasons (id, project_id, season_number) VALUES
(1, 2, 1)
ON CONFLICT (project_id, season_number) DO NOTHING;

-- ==========================
-- EPISODES
-- ==========================

INSERT INTO episodes (id, season_id, episode_number, title, youtube_video_id, duration) VALUES
(1, 1, 1, '1-бөлім', 'episode001', 1500),
(2, 1, 2, '2-бөлім', 'episode002', 1450)
ON CONFLICT (season_id, episode_number) DO NOTHING;

-- ==========================
-- FEATURED CONTENT
-- ==========================

INSERT INTO featured_content (project_id, block_type, sort_order) VALUES
(1, 'HERO_BANNER', 1),
(2, 'TRENDING', 1)
ON CONFLICT (block_type, sort_order) DO NOTHING;