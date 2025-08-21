#!/bin/bash

# Test script for OTP Authentication Service
# Make sure the service is running on localhost:8080

BASE_URL="http://localhost:8080/api/v1"
PHONE_NUMBER="+1234567890"

echo "=== OTP Authentication Service Test ==="
echo

# Test 1: Request OTP
echo "1. Requesting OTP for $PHONE_NUMBER..."
RESPONSE=$(curl -s -X POST "$BASE_URL/auth/request-otp" \
  -H "Content-Type: application/json" \
  -d "{\"phone_number\": \"$PHONE_NUMBER\"}")

echo "Response: $RESPONSE"
echo

# Note: In a real scenario, you would get the OTP from SMS
# For this demo, check the console output of the running service
echo "2. Check the console output of the running service for the OTP code"
echo "   Then use that OTP in the next step"
echo

# Test 2: Verify OTP (replace 123456 with actual OTP from console)
echo "3. Verifying OTP (replace 123456 with actual OTP)..."
echo "   curl -X POST \"$BASE_URL/auth/verify-otp\" \\"
echo "     -H \"Content-Type: application/json\" \\"
echo "     -d '{\"phone_number\": \"$PHONE_NUMBER\", \"otp\": \"123456\"}'"
echo

# Test 3: Get Users (requires JWT token)
echo "4. To get users list (requires JWT token from step 3):"
echo "   curl -X GET \"$BASE_URL/users\" \\"
echo "     -H \"Authorization: Bearer YOUR_JWT_TOKEN\""
echo

# Test 4: Health check
echo "5. Testing health check..."
HEALTH=$(curl -s "$BASE_URL/../health")
echo "Health check response: $HEALTH"
echo

echo "=== Test completed ==="
echo "Access Swagger UI at: http://localhost:8080/swagger/index.html"
