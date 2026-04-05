CREATE TABLE IF NOT EXISTS muscle_groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS machines (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    is_available BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE IF NOT EXISTS exercises (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    muscle_group_id INTEGER NOT NULL,
    FOREIGN KEY (muscle_group_id) REFERENCES muscle_groups(id)
);

CREATE TABLE IF NOT EXISTS exercise_machines (
    id SERIAL PRIMARY KEY,
    exercise_id INTEGER NOT NULL,
    machine_id INTEGER NOT NULL,
    FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE,
    FOREIGN KEY (machine_id) REFERENCES machines(id) ON DELETE CASCADE,
    UNIQUE (exercise_id, machine_id)
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE
);
INSERT INTO muscle_groups (name) VALUES
('Chest'),
('Back'),
('Legs'),
('Shoulders'),
('Arms');

INSERT INTO machines (name, is_available) VALUES
('Bench', false),
('Cable Machine', true),
('Lat Pulldown Machine', true),
('Leg Press Machine', true),
('Dumbbells', true),
('Shoulder Press Machine', true);

INSERT INTO exercises (name, muscle_group_id) VALUES
('Bench Press', 1),
('Cable Fly', 1),
('Lat Pulldown', 2),
('Seated Row', 2),
('Leg Press', 3),
('Dumbbell Lunge', 3),
('Shoulder Press', 4),
('Lateral Raise', 4),
('Bicep Curl', 5),
('Tricep Pushdown', 5);

INSERT INTO exercise_machines (exercise_id, machine_id) VALUES
(1, 1),
(2, 2),
(3, 3),
(4, 2),
(5, 4),
(6, 5),
(7, 6),
(8, 5),
(9, 5),
(10, 2);

INSERT INTO users (name, email) VALUES
('Test User', 'testuser@flowgym.com');
