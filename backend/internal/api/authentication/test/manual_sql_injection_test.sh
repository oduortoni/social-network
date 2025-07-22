#!/bin/bash

# Manual SQL Injection and XSS Test Script for Login and Signup
# This script tests various SQL injection and XSS payloads against the login and signup endpoints
# Note: The automated tests are now located in internal/api/authentication/test/

echo "=== SQL Injection and XSS Testing for Login and Signup ==="
echo "Starting server in background..."

# Start the server in background (assuming it runs on port 8080)
# Uncomment the next line if you want to start the server automatically
# go run ../cmd/main.go &
# SERVER_PID=$!
# sleep 3

BASE_URL="http://localhost:8080"
LOGIN_ENDPOINT="$BASE_URL/login"
SIGNUP_ENDPOINT="$BASE_URL/signup"

echo "Testing endpoints: $LOGIN_ENDPOINT, $SIGNUP_ENDPOINT"
echo ""

# --- Login SQL Injection Tests ---

# Test 1: Basic SQL injection in email field
echo "=== Login SQL Injection Test 1: Basic SQL injection in email field ==="
echo "Payload: ' OR '1'='1"
curl -s -X POST "$LOGIN_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{"email": "'\'' OR '\''1'\''='\''1", "password": "TestPassword123!"}' \
  | jq '.' || echo "Response not JSON"
echo ""

# Test 2: Union-based injection
echo "=== Login SQL Injection Test 2: Union-based injection ==="
echo "Payload: ' UNION SELECT 1,'admin','password123'--"
curl -s -X POST "$LOGIN_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com'\'' UNION SELECT 1,'\''admin'\'','\''password123'\''--", "password": "TestPassword123!"}' \
  | jq '.' || echo "Response not JSON"
echo ""

# Test 3: Stacked queries (dangerous)
echo "=== Login SQL Injection Test 3: Stacked queries (dangerous) ==="
echo "Payload: test@example.com'; DROP TABLE Users;--"
curl -s -X POST "$LOGIN_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com'\''; DROP TABLE Users;--", "password": "TestPassword123!"}' \
  | jq '.' || echo "Response not JSON"
echo ""

# Test 4: SQL injection in password field
echo "=== Login SQL Injection Test 4: SQL injection in password field ==="
echo "Payload: ' OR 1=1--"
curl -s -X POST "$LOGIN_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "'\'' OR 1=1--"}' \
  | jq '.' || echo "Response not JSON"
echo ""

# Test 5: Form data injection
echo "=== Login SQL Injection Test 5: Form data injection ==="
echo "Payload: ' OR '1'='1 (via form data)"
curl -s -X POST "$LOGIN_ENDPOINT" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "email=' OR '1'='1&password=TestPassword123!" \
  | jq '.' || echo "Response not JSON"
echo ""

# Test 6: Valid login (should work)
echo "=== Login SQL Injection Test 6: Valid login (should work if user exists) ==="
echo "Credentials: test@example.com / TestPassword123!"
curl -s -X POST "$LOGIN_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "TestPassword123!"}' \
  | jq '.' || echo "Response not JSON"
echo ""

# --- Signup SQL Injection and XSS Tests ---

# Test 7: SQL injection in signup email field
echo "=== Signup SQL Injection Test 1: SQL injection in email field ==="
echo "Payload: ' OR '1'='1"
curl -s -X POST "$SIGNUP_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{"email": "'\'' OR '\''1'\''='\''1", "password": "TestPassword123!"}' \
  | jq '.' || echo "Response not JSON"
echo ""

# Test 8: SQL injection in signup password field
echo "=== Signup SQL Injection Test 2: SQL injection in password field ==="
echo "Payload: ' OR '1'='1"
curl -s -X POST "$SIGNUP_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "'\'' OR '\''1'\''='\''1"}' \
  | jq '.' || echo "Response not JSON"
echo ""

# Test 9: XSS in signup email field
echo "=== Signup XSS Test 1: XSS in email field ==="
echo "Payload: <script>alert('XSS')</script>"
curl -s -X POST "$SIGNUP_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{"email": "<script>alert('\''XSS'\'')</script>", "password": "TestPassword123!"}' \
  | jq '.' || echo "Response not JSON"
echo ""

# Test 10: XSS in signup password field
echo "=== Signup XSS Test 2: XSS in password field ==="
echo "Payload: <script>alert('XSS')</script>"
curl -s -X POST "$SIGNUP_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "<script>alert('\''XSS'\'')</script>"}' \
  | jq '.' || echo "Response not JSON"
echo ""

# Test 11: Valid signup (should work)
echo "=== Signup SQL Injection and XSS Test 3: Valid signup (should work if email is unique) ==="
echo "Credentials: test@example.com / TestPassword123!"
curl -s -X POST "$SIGNUP_ENDPOINT" \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "TestPassword123!"}' \
  | jq '.' || echo "Response not JSON"
echo ""

echo "=== Testing Complete ==="
echo ""
echo "Expected Results:"
echo "- Tests 1-10 should FAIL (return 401 Unauthorized or similar)"
echo "- Test 11 should SUCCEED only if the test user does not already exist"
echo "- No SQL injection or XSS should succeed"
echo "- Database should remain intact"

# Cleanup
# if [ ! -z "$SERVER_PID" ]; then
#   echo "Stopping server..."
#   kill $SERVER_PID
# fi
