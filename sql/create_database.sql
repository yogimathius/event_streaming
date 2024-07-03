SELECT 'CREATE DATABASE event_streaming' 
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'event_streaming')\gexec;