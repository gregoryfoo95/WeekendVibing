-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    level INTEGER DEFAULT 1,
    points INTEGER DEFAULT 0,
    character VARCHAR(255) DEFAULT 'Rookie Hero',
    job_title VARCHAR(255) DEFAULT 'Fitness Novice',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Tasks table
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    points INTEGER NOT NULL,
    category VARCHAR(100),
    difficulty VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Daily tasks table
CREATE TABLE daily_tasks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    task_id INTEGER REFERENCES tasks(id) ON DELETE CASCADE,
    completed BOOLEAN DEFAULT FALSE,
    date TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

-- Achievements table
CREATE TABLE achievements (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    icon VARCHAR(255),
    points_cost INTEGER NOT NULL,
    type VARCHAR(50) NOT NULL, -- character, upgrade, badge
    created_at TIMESTAMP DEFAULT NOW()
);

-- User achievements table
CREATE TABLE user_achievements (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    achievement_id INTEGER REFERENCES achievements(id) ON DELETE CASCADE,
    unlocked_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, achievement_id)
);

-- Insert sample tasks
INSERT INTO tasks (title, description, points, category, difficulty) VALUES
-- Easy tasks (10-20 points)
('Take a 5-minute walk', 'Step outside and take a short walk around your neighborhood or office building', 10, 'cardio', 'easy'),
('Do 10 jumping jacks', 'Simple jumping exercise to get your heart pumping', 15, 'cardio', 'easy'),
('Stretch for 3 minutes', 'Gentle stretching to improve flexibility and reduce tension', 10, 'flexibility', 'easy'),
('Drink a glass of water', 'Stay hydrated! Drink 8oz of water right now', 5, 'wellness', 'easy'),
('Take the stairs instead of elevator', 'Choose stairs over elevator for the rest of the day', 15, 'cardio', 'easy'),
('Do 5 push-ups (modified okay)', 'Standard or knee push-ups - your choice!', 20, 'strength', 'easy'),
('Stand up and sit down 10 times', 'Great exercise you can do right at your desk', 15, 'strength', 'easy'),
('Deep breathing for 2 minutes', 'Practice mindful breathing to reduce stress', 10, 'wellness', 'easy'),

-- Medium tasks (20-40 points)
('Walk for 15 minutes', 'Take a longer walk to boost your energy and mood', 25, 'cardio', 'medium'),
('Do 20 squats', 'Bodyweight squats to strengthen your legs and glutes', 30, 'strength', 'medium'),
('Plank for 1 minute', 'Core strengthening exercise - hold that plank!', 35, 'strength', 'medium'),
('Dance to 2 songs', 'Put on your favorite music and dance like nobody is watching', 25, 'cardio', 'medium'),
('Do yoga for 10 minutes', 'Follow a short yoga routine or app', 30, 'flexibility', 'medium'),
('Walk up 3 flights of stairs', 'Great cardio workout using stairs', 25, 'cardio', 'medium'),
('Do 15 lunges (each leg)', 'Forward or reverse lunges to work your legs', 30, 'strength', 'medium'),
('Wall sit for 45 seconds', 'Lean against a wall and hold the squat position', 35, 'strength', 'medium'),

-- Hard tasks (40-60 points)
('30-minute walk or jog', 'Longer cardio session to really get your heart rate up', 50, 'cardio', 'hard'),
('Do 50 jumping jacks', 'High-energy cardio exercise', 40, 'cardio', 'hard'),
('Hold plank for 2 minutes', 'Advanced core strengthening', 55, 'strength', 'hard'),
('Do 25 push-ups', 'Upper body strength challenge', 45, 'strength', 'hard'),
('100 bodyweight squats', 'Leg day challenge - pace yourself!', 60, 'strength', 'hard'),
('15-minute HIIT workout', 'High-intensity interval training session', 55, 'cardio', 'hard'),
('20-minute bike ride', 'Outdoor cycling or stationary bike', 50, 'cardio', 'hard'),
('Burpee challenge: 10 burpees', 'Full-body explosive exercise', 60, 'strength', 'hard');

-- Insert sample achievements
INSERT INTO achievements (title, description, icon, points_cost, type) VALUES
-- Character upgrades
('Fitness Apprentice', 'Upgrade to the next superhero level!', '🦸‍♂️', 100, 'character'),
('Health Guardian', 'Become a guardian of your own health', '🛡️', 250, 'character'),
('Wellness Warrior', 'Fight the good fight against sedentary lifestyle', '⚔️', 500, 'character'),
('Fitness Champion', 'Champion level hero with incredible strength', '🏆', 1000, 'character'),
('Ultimate Hero', 'The pinnacle of fitness achievement', '👑', 2000, 'character'),

-- Job titles
('Personal Trainer', 'Advance your career in fitness', '💪', 150, 'upgrade'),
('Fitness Coach', 'Help others on their fitness journey', '🎯', 300, 'upgrade'),
('Wellness Expert', 'Master of health and wellness', '🧠', 600, 'upgrade'),
('Fitness Director', 'Lead the fitness revolution', '👔', 1200, 'upgrade'),
('Health Guru', 'The ultimate fitness professional', '🌟', 2500, 'upgrade'),

-- Special badges
('Early Bird', 'Complete morning workouts consistently', '🌅', 200, 'badge'),
('Consistency Master', 'Complete daily tasks for 7 days straight', '📈', 300, 'badge'),
('Cardio King/Queen', 'Master of cardiovascular exercises', '❤️', 400, 'badge'),
('Strength Legend', 'Legend in strength training', '💪', 400, 'badge'),
('Flexibility Master', 'Master of stretching and flexibility', '🤸', 400, 'badge'),
('Hydration Hero', 'Stay hydrated like a true hero', '💧', 150, 'badge'),
('Stair Climber', 'Master of vertical challenges', '🪜', 250, 'badge'),
('Weekend Warrior', 'Active even on weekends', '⚡', 350, 'badge');

-- Create indexes for better performance
CREATE INDEX idx_daily_tasks_user_date ON daily_tasks(user_id, date);
CREATE INDEX idx_user_achievements_user ON user_achievements(user_id);
CREATE INDEX idx_users_points ON users(points DESC); 