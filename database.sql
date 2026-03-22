CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE appointments (
    id VARCHAR(50) PRIMARY KEY, 
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    task_date DATE NOT NULL, 
    name VARCHAR(80) NOT NULL, 
    description TEXT, 
    time_start TIME, 
    time_end TIME, 
    alert_type VARCHAR(20) DEFAULT 'none', 
    CONSTRAINT chk_alert_type CHECK (alert_type IN ('none', '1_day_before', 'custom_date')),
    alert_dates DATE[], 
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_tasks_task_date ON tasks(task_date);
CREATE INDEX idx_tasks_user_id ON tasks(user_id);