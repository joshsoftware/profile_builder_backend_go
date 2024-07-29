-- CREATE DATABASE profile_builder
--   TEMPLATE template0
--   ENCODING 'UTF8'
--   LC_COLLATE = 'en_US.UTF-8'
--   LC_CTYPE = 'en_US.UTF-8'
--   CONNECTION LIMIT = -1;


CREATE TABLE IF NOT EXISTS users(
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	email VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS profiles (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    gender VARCHAR(50),
    mobile char(10) NOT NULL UNIQUE,
    designation VARCHAR(100),
    description TEXT,
    title VARCHAR(150) NOT NULL,
    years_of_experience FLOAT NOT NULL,
    primary_skills TEXT[],
    secondary_skills TEXT[],
    github_link VARCHAR(100),
    linkedin_link VARCHAR(100),
    career_objectives TEXT,
    is_active INT NOT NULL,
    is_current_employee INT NOT NULL,
    created_at DATE NOT NULL DEFAULT CURRENT_DATE,
	updated_at DATE NOT NULL DEFAULT CURRENT_DATE,
    created_by_id INT NOT NULL,
    updated_by_id INT NOT NULL
);

CREATE TABLE IF NOT EXISTS educations (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	degree VARCHAR(100) NOT NULL,
    university_name VARCHAR(100),
    place VARCHAR(100),
    percent_or_cgpa VARCHAR(40),
    passing_year VARCHAR(50),
    priorities int NOT NULL,
	created_at DATE NOT NULL DEFAULT CURRENT_DATE,
	updated_at DATE NOT NULL DEFAULT CURRENT_DATE,
    created_by_id INT NOT NULL,
    updated_by_id INT NOT NULL,
    profile_id INT NOT NULL,
        
    CONSTRAINT fk_profile_id
		FOREIGN KEY(profile_id)
		REFERENCES profiles(id)
		ON DELETE CASCADE,
		
	CONSTRAINT unique_degree_profile UNIQUE (degree, profile_id)
);

CREATE TABLE IF NOT EXISTS certificates (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
    organization_name VARCHAR(100),
    description TEXT,
    issued_date VARCHAR(50),
    from_date VARCHAR(50),
    to_date VARCHAR(50),
    priorities int NOT NULL,
	created_at DATE NOT NULL DEFAULT CURRENT_DATE,
	updated_at DATE NOT NULL DEFAULT CURRENT_DATE,
    created_by_id INT NOT NULL,
    updated_by_id INT NOT NULL,
    profile_id INT NOT NULL,
        
    CONSTRAINT fk_profile_id1
		FOREIGN KEY(profile_id)
		REFERENCES profiles(id)
		ON DELETE CASCADE,
		
	CONSTRAINT unique_name_profile_cert UNIQUE (name, profile_id)
);

CREATE TABLE IF NOT EXISTS projects (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
    description TEXT,
    role VARCHAR(50),
    responsibilities VARCHAR(500),
    technologies TEXT[],
    tech_worked_on TEXT[] NOT NULL,    
    working_start_date VARCHAR(50),
    working_end_date VARCHAR(50),
    duration VARCHAR(50),
    priorities int NOT NULL,
	created_at DATE NOT NULL DEFAULT CURRENT_DATE,
	updated_at DATE NOT NULL DEFAULT CURRENT_DATE,
    created_by_id INT NOT NULL,
    updated_by_id INT NOT NULL,
    profile_id INT NOT NULL,
        
    CONSTRAINT fk_profile_id2
		FOREIGN KEY(profile_id)
		REFERENCES profiles(id)
		ON DELETE CASCADE,
		
	CONSTRAINT unique_name_profile_proj UNIQUE (name, profile_id)
);

CREATE TABLE IF NOT EXISTS experiences (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	designation TEXT NOT NULL,
    company_name VARCHAR(100) NOT NULL,
    from_date VARCHAR(50) NOT NULL,
    to_date VARCHAR(50),
    priorities int NOT NULL,
	created_at DATE NOT NULL DEFAULT CURRENT_DATE,
	updated_at DATE NOT NULL DEFAULT CURRENT_DATE,
    created_by_id INT NOT NULL,
    updated_by_id INT NOT NULL,
    profile_id INT NOT NULL,
        
    CONSTRAINT fk_profile_id3
		FOREIGN KEY(profile_id)
		REFERENCES profiles(id)
		ON DELETE CASCADE,
		
	CONSTRAINT unique_designation_company_profile UNIQUE (designation, company_name, profile_id)
);

CREATE TABLE IF NOT EXISTS achievements (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
	description TEXT,
    priorities int NOT NULL,
	created_at DATE NOT NULL DEFAULT CURRENT_DATE,
	updated_at DATE NOT NULL DEFAULT CURRENT_DATE,
    created_by_id INT NOT NULL,
    updated_by_id INT NOT NULL,
    profile_id INT NOT NULL,
        
    CONSTRAINT fk_profile_id4
		FOREIGN KEY(profile_id)
		REFERENCES profiles(id)
		ON DELETE CASCADE,
		
	CONSTRAINT unique_name_profile UNIQUE (name, profile_id)
);
