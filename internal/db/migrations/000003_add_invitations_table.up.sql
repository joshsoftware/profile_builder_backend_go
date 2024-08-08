CREATE TABLE IF NOT EXISTS invitations (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	profile_id INT NOT NULL,
	is_profile_complete INT NOT NULL CHECK (is_profile_complete BETWEEN 0 AND 1),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by_id INT NOT NULL,
	updated_by_id INT NOT NULL,
	
	CONSTRAINT fk_profile_id5
		FOREIGN KEY(profile_id)
		REFERENCES profiles(id)
		ON DELETE CASCADE
);
-- check uniique constraint for profile_id and profile_complete
CREATE UNIQUE INDEX unique_profile_complete_zero ON invitations (profile_id) 
WHERE is_profile_complete = 0;