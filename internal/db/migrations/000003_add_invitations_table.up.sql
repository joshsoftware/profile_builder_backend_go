CREATE TABLE IF NOT EXISTS invitations (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	profile_id INT NOT NULL,
	profile_complete INT NOT NULL CHECK (profile_complete BETWEEN 0 AND 1),
	created_at DATE NOT NULL DEFAULT CURRENT_DATE,
	updated_at DATE NOT NULL DEFAULT CURRENT_DATE,
	created_by_id INT NOT NULL,
	updated_by_id INT NOT NULL,
	
	CONSTRAINT fk_profile_id5
		FOREIGN KEY(profile_id)
		REFERENCES profiles(id)
		ON DELETE CASCADE
);
-- check uniique constraint for profile_id and profile_complete
CREATE UNIQUE INDEX unique_profile_complete_zero ON invitations (profile_id) 
WHERE profile_complete = 0;