\c event_streaming

CREATE TABLE IF NOT EXISTS "events" (
  "event_id" SERIAL PRIMARY KEY,
  "guest_satisfaction" BOOLEAN, -- TRUE if happy, FALSE if upset
  "stress_marks" INT DEFAULT 0, -- Number of stress marks recorded
  "timestamp" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "event_messages" (
  "event_message_id" SERIAL PRIMARY KEY,
  "event_id" INT,
  "event_type" VARCHAR(100) NOT NULL,
  "priority" VARCHAR(10) NOT NULL,
  "description" TEXT,
  "timestamp" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("event_id") REFERENCES "events" ("event_id")
);

CREATE TABLE IF NOT EXISTS "teams" (
  "team_id" SERIAL PRIMARY KEY,
  "team_name" VARCHAR(50) NOT NULL,
  "routine" VARCHAR(20) NOT NULL
);

CREATE TABLE IF NOT EXISTS "workers" (
  "worker_id" SERIAL PRIMARY KEY,
  "team_id" INT,
  "current_status" VARCHAR(10) NOT NULL,
  "availability" BOOLEAN DEFAULT true,
  FOREIGN KEY ("team_id") REFERENCES "teams" ("team_id")
);

CREATE TABLE IF NOT EXISTS "event_logs" (
  "log_id" SERIAL PRIMARY KEY,
  "event_message_id" INT,
  "event_type" VARCHAR(100) NOT NULL,
  "status" VARCHAR(10) NOT NULL,
  "stress_level" INT,
  "log_timestamp" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY ("event_message_id") REFERENCES "event_messages" ("event_message_id")
);

-- Add indexes

CREATE INDEX IF NOT EXISTS idx_event_id ON "events" ("event_id");
CREATE INDEX IF NOT EXISTS idx_event_message_id ON "events" ("event_message_id");
CREATE INDEX IF NOT EXISTS idx_event_message_message_id ON "event_messages" ("event_message_id");
CREATE INDEX IF NOT EXISTS idx_team_id ON "workers" ("team_id");
CREATE INDEX IF NOT EXISTS idx_event_message_log_id ON "event_logs" ("event_message_id");
