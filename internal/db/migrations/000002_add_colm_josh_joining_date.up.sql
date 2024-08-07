-- 000002_add_age_to_users.up.sql
ALTER TABLE profiles ADD COLUMN josh_joining_date varchar(50);
ALTER TABLE projects ALTER COLUMN responsibilities TYPE TEXT;
