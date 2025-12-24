#!/bin/bash
# Generate sample mixed logs (JSON and non-JSON)

echo "Starting application..."
echo '{"level":"info","message":"Server starting","port":8080,"timestamp":"2025-01-01T10:00:00Z"}'
echo "DEBUG: Loading configuration"
echo '{"level":"debug","message":"Config loaded","config":{"database":"postgres","cache":"redis"},"timestamp":"2025-01-01T10:00:01Z"}'
echo '{"level":"info","message":"Database connected","host":"localhost","port":5432,"timestamp":"2025-01-01T10:00:02Z"}'
echo "WARNING: This is a plain text warning"
echo '{"level":"warn","message":"High memory usage","memory_percent":85.5,"timestamp":"2025-01-01T10:00:03Z"}'
echo '{"level":"error","message":"Request failed","error":"connection timeout","request_id":"abc123","user":{"id":42,"name":"John"},"timestamp":"2025-01-01T10:00:04Z"}'
echo "Plain text log line"
echo '{"level":"info","message":"Request completed","method":"GET","path":"/api/users","status":200,"duration_ms":45,"timestamp":"2025-01-01T10:00:05Z"}'
echo '{"level":"debug","message":"Cache hit","key":"user:42","ttl":3600,"timestamp":"2025-01-01T10:00:06Z"}'
echo 'Not valid JSON {'
echo '{"level":"info","message":"Graceful shutdown initiated","active_connections":10,"timestamp":"2025-01-01T10:00:07Z"}'
