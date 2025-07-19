#!/bin/bash

# API Test Script for Document Reader Chatbot
# Make sure the server is running on localhost:8080

BASE_URL="http://localhost:8080/api/v1"

echo "🚀 Testing Document Reader Chatbot API"
echo "========================================="

# Test Registration
echo "📝 Testing User Registration..."
REGISTER_RESPONSE=$(curl -s -X POST \
  ${BASE_URL}/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john.doe@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }')

echo "Registration Response:"
echo $REGISTER_RESPONSE | jq '.'
echo ""

# Test Login with Email
echo "🔐 Testing Login with Email..."
LOGIN_EMAIL_RESPONSE=$(curl -s -X POST \
  ${BASE_URL}/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email_or_username": "john.doe@example.com",
    "password": "password123"
  }')

echo "Login with Email Response:"
echo $LOGIN_EMAIL_RESPONSE | jq '.'

# Extract token for further requests
TOKEN=$(echo $LOGIN_EMAIL_RESPONSE | jq -r '.data.token')
echo "JWT Token: $TOKEN"
echo ""

# Test Login with Username  
echo "🔐 Testing Login with Username..."
LOGIN_USERNAME_RESPONSE=$(curl -s -X POST \
  ${BASE_URL}/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email_or_username": "johndoe",
    "password": "password123"
  }')

echo "Login with Username Response:"
echo $LOGIN_USERNAME_RESPONSE | jq '.'
echo ""

# Test Profile Access
echo "👤 Testing Profile Access..."
PROFILE_RESPONSE=$(curl -s -X GET \
  ${BASE_URL}/profile \
  -H "Authorization: Bearer $TOKEN")

echo "Profile Response:"
echo $PROFILE_RESPONSE | jq '.'
echo ""

# Test Invalid Login
echo "❌ Testing Invalid Login..."
INVALID_LOGIN_RESPONSE=$(curl -s -X POST \
  ${BASE_URL}/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email_or_username": "nonexistent@example.com",
    "password": "wrongpassword"
  }')

echo "Invalid Login Response:"
echo $INVALID_LOGIN_RESPONSE | jq '.'
echo ""

echo "✅ API Testing Complete!" 