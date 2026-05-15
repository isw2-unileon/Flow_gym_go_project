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

-- Table for storing user routines
CREATE TABLE IF NOT EXISTS routines (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Intermediate table for linking routines with exercises and maintaining the sequence of the routine
CREATE TABLE IF NOT EXISTS routine_exercises (
    id SERIAL PRIMARY KEY,
    routine_id INTEGER NOT NULL,
    exercise_id INTEGER NOT NULL,
    exercise_order INTEGER NOT NULL, -- Determine which exercise comes first, second, and so on.
    FOREIGN KEY (routine_id) REFERENCES routines(id) ON DELETE CASCADE,
    FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE,
    UNIQUE (routine_id, exercise_id) -- The same exercise should not be repeated in the same routine
);

-- (Optional) Sample data
INSERT INTO routines (user_id, name) VALUES (1, 'Full Body Workout');
INSERT INTO routine_exercises (routine_id, exercise_id, exercise_order) VALUES 
(1, 1, 1), -- Bench Press 1º
(1, 3, 2), -- Lat Pulldown 2º
(1, 5, 3); -- Leg Press 3º