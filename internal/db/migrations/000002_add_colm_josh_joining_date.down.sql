-- 000002_add_age_to_users.down.sql
ALTER TABLE profiles DROP COLUMN josh_joining_date;
ALTER TABLE projects ALTER COLUMN responsibilities TYPE varchar(500);
