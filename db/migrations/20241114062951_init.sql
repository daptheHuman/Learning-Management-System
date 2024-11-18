-- migrate:up
CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE curriculums (
    id SERIAL PRIMARY KEY,
    course_id INT REFERENCES courses(id) ON DELETE CASCADE,
    section_name VARCHAR(255) NOT NULL,
    section_order INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE materials (
    id SERIAL PRIMARY KEY,
    curriculum_id INT REFERENCES curriculums(id) ON DELETE CASCADE,
    material_type VARCHAR(50), -- e.g., "text", "video", "quiz"
    content TEXT, -- Content can be JSON for quizzes, or text for plain material
    "order" INT NOT NULL DEFAULT 1, -- Position of the material within the curriculum
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TYPE user_role AS ENUM ('student', 'instructor', 'admin');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL, -- Renamed for clarity and security
    role user_role NOT NULL DEFAULT 'student',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE enrollments (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE, 
    course_id INT REFERENCES courses(id) ON DELETE CASCADE,
    enrolled_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE progress_tracking (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE, 
    curriculum_id INT REFERENCES curriculums(id) ON DELETE CASCADE, -- Associated curriculum
    material_id INT REFERENCES materials(id) ON DELETE CASCADE, -- Specific material
    status VARCHAR(50), -- e.g., "in-progress", "completed"
    progress_percentage INT DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE student_paths (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE, 
    curriculum_id INT REFERENCES curriculums(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'not started', -- e.g., "in-progress", "completed"
    progress_percentage INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE assessments (
    id SERIAL PRIMARY KEY,
    course_id INT REFERENCES courses(id) ON DELETE CASCADE,
    type VARCHAR(50), -- e.g., "multiple-choice", "essay"
    question TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE submissions (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE, 
    assessment_id INT REFERENCES assessments(id) ON DELETE CASCADE,
    answer TEXT,
    grade INT,
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE achievements (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE, 
    course_id INT REFERENCES courses(id) ON DELETE CASCADE,
    achievement_type VARCHAR(255), -- e.g., "course completion"
    description TEXT,
    awarded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- migrate:down
DROP TABLE achievements;
DROP TABLE submissions;
DROP TABLE assessments;
DROP TABLE progress_tracking;
DROP TABLE enrollments;
DROP TABLE users;
DROP TABLE materials;
DROP TABLE curriculum;
DROP TABLE courses;

