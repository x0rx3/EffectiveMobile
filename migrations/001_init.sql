-- Включаем расширение для генерации UUID
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Таблица подписок с UUID
CREATE TABLE IF NOT EXISTS subscribes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL,
    price INT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL
);

-- Добавим три тестовые подписки с фиксированными UUID для id и user_id
INSERT INTO subscribes (id, name, user_id, price, start_date, end_date) VALUES
('11111111-1111-1111-1111-111111111111', 'Netflix', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 999, '2025-01-01', '2025-02-01'),
('22222222-2222-2222-2222-222222222222', 'Spotify', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 499, '2025-01-15', '2025-03-15'),
('33333333-3333-3333-3333-333333333333', 'YouTube Premium', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 699, '2025-02-01', '2025-03-01');
