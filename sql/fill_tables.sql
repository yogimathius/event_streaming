-- Insert sample teams
INSERT INTO "teams" ("team_name", "routine")
VALUES 
  ('Security Team', 'Standard'),
  ('Clean Up Team', 'Intermittent'),
  ('Catering Team', 'Concentrated'),
  ('Officiant Team', 'Standard'),
  ('Waiters Team', 'Intermittent');

-- Insert workers into Security Team
INSERT INTO "workers" ("team_id", "current_status", "availability")
SELECT 1, 'Idle', true
FROM generate_series(1, 2);

-- Insert workers into Clean Up Team
INSERT INTO "workers" ("team_id", "current_status", "availability")
SELECT 2, 'Idle', true
FROM generate_series(1, 4);

-- Insert workers into Catering Team
INSERT INTO "workers" ("team_id", "current_status", "availability")
SELECT 3, 'Idle', true
FROM generate_series(1, 10); 

-- Insert workers into Officiant Team
INSERT INTO "workers" ("team_id", "current_status", "availability")
SELECT 4, 'Idle', true
FROM generate_series(1, 1);

-- Insert workers into Waiters Team
INSERT INTO "workers" ("team_id", "current_status", "availability")
SELECT 5, 'Idle', true
FROM generate_series(1, 8);
