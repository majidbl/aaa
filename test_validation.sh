#!/bin/bash

# Test script for Input Validation
# Make sure the service is running on localhost:8080

BASE_URL="http://localhost:8080/api/v1"

echo "=== Input Validation Test Suite ==="
echo

# Test 1: Invalid phone number format
echo "1. Testing invalid phone number format..."
echo "   Request: POST /auth/request-otp with invalid phone number"
RESPONSE=$(curl -s -X POST "$BASE_URL/auth/request-otp" \
  -H "Content-Type: application/json" \
  -d '{"phone_number": "invalid-phone"}')

echo "   Response: $RESPONSE"
echo

# Test 2: Missing phone number
echo "2. Testing missing phone number..."
echo "   Request: POST /auth/request-otp with missing phone_number"
RESPONSE=$(curl -s -X POST "$BASE_URL/auth/request-otp" \
  -H "Content-Type: application/json" \
  -d '{}')

echo "   Response: $RESPONSE"
echo

# Test 3: Invalid OTP format
echo "3. Testing invalid OTP format..."
echo "   Request: POST /auth/verify-otp with invalid OTP"
RESPONSE=$(curl -s -X POST "$BASE_URL/auth/verify-otp" \
  -H "Content-Type: application/json" \
  -d '{"phone_number": "+1234567890", "otp": "123"}')

echo "   Response: $RESPONSE"
echo

# Test 4: Invalid pagination parameters
echo "4. Testing invalid pagination parameters..."
echo "   Request: GET /users with invalid page/limit"
RESPONSE=$(curl -s -X GET "$BASE_URL/users?page=0&limit=200")

echo "   Response: $RESPONSE"
echo

# Test 5: Invalid search query
echo "5. Testing invalid search query..."
echo "   Request: GET /users with invalid search query"
RESPONSE=$(curl -s -X GET "$BASE_URL/users?search=ab")

echo "   Response: $RESPONSE"
echo

# Test 6: Invalid UUID format
echo "6. Testing invalid UUID format..."
echo "   Request: GET /users with invalid user ID"
RESPONSE=$(curl -s -X GET "$BASE_URL/users/invalid-uuid")

echo "   Response: $RESPONSE"
echo

# Test 7: Valid phone number format
echo "7. Testing valid phone number format..."
echo "   Request: POST /auth/request-otp with valid phone number"
RESPONSE=$(curl -s -X POST "$BASE_URL/auth/request-otp" \
  -H "Content-Type: application/json" \
  -d '{"phone_number": "+1234567890"}')

echo "   Response: $RESPONSE"
echo

echo "=== Validation Test Completed ==="
echo
echo "Expected Results:"
echo "- Tests 1-6 should return 400 Bad Request with validation errors"
echo "- Test 7 should return 200 OK (if Redis is running)"
echo
echo "Validation Features:"
echo "✓ Phone number format validation (international format required)"
echo "✓ OTP format validation (exactly 6 digits)"
echo "✓ UUID format validation"
echo "✓ Pagination parameter validation"
echo "✓ Search query validation"
echo "✓ Required field validation"
echo "✓ Input sanitization (trim whitespace)"
