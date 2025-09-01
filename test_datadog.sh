#!/bin/bash

# MasterCom Service with Datadog Integration Test Script
# This script demonstrates the Datadog tracing and logging features

BASE_URL="http://localhost:8080"
API_BASE="$BASE_URL/api/v6"

echo "🚀 Testing MasterCom Service with Datadog Integration"
echo "=================================================="

# Test 1: Health Check with Trace Info
echo -e "\n1. Testing Health Check with Trace Information..."
HEALTH_RESPONSE=$(curl -s -X GET "$BASE_URL/health")
echo "$HEALTH_RESPONSE" | jq '.'

# Extract trace information
TRACE_ID=$(echo "$HEALTH_RESPONSE" | jq -r '.dd_trace_id')
SPAN_ID=$(echo "$HEALTH_RESPONSE" | jq -r '.dd_span_id')

echo "Trace ID: $TRACE_ID"
echo "Span ID: $SPAN_ID"

# Test 2: Create a Case with Request ID for tracing
echo -e "\n2. Creating a case with request tracing..."
REQUEST_ID=$(uuidgen)
CASE_RESPONSE=$(curl -s -X POST "$API_BASE/cases" \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: $REQUEST_ID" \
  -d '{
    "caseType": "PRE_ARBITRATION",
    "primaryAccountNumber": "4111111111111111",
    "transactionAmount": 100.00,
    "transactionCurrency": "USD",
    "transactionDate": "2024-01-01T00:00:00Z",
    "transactionId": "123456789",
    "merchantName": "Test Merchant",
    "merchantCategoryCode": "5411",
    "reasonCode": "10.1",
    "disputeAmount": 100.00,
    "disputeCurrency": "USD",
    "filingAs": "ISSUER",
    "filingIca": "123456",
    "filedAgainstIca": "654321",
    "filedBy": "Test User",
    "filedByContactName": "John Doe",
    "filedByContactPhone": "+1234567890",
    "filedByContactEmail": "john.doe@example.com"
  }')

echo "$CASE_RESPONSE" | jq '.'

# Extract case ID for subsequent tests
CASE_ID=$(echo "$CASE_RESPONSE" | jq -r '.id')
echo "Created case ID: $CASE_ID"

# Test 3: Upload a Document with tracing
echo -e "\n3. Uploading a document with tracing..."
# Create a temporary test file
echo "This is a test document for Datadog tracing demonstration." > test_document.txt

DOC_RESPONSE=$(curl -s -X POST "$API_BASE/documents" \
  -H "X-Request-ID: $REQUEST_ID" \
  -F "file=@test_document.txt" \
  -F "caseId=$CASE_ID" \
  -F "description=Test document for Datadog tracing" \
  -F "uploadedBy=test-user")

echo "$DOC_RESPONSE" | jq '.'

# Extract document ID
DOC_ID=$(echo "$DOC_RESPONSE" | jq -r '.id')
echo "Uploaded document ID: $DOC_ID"

# Test 4: Get the case and verify trace headers
echo -e "\n4. Getting case with trace headers..."
CASE_GET_RESPONSE=$(curl -s -X GET "$API_BASE/cases/$CASE_ID" \
  -H "X-Request-ID: $REQUEST_ID" \
  -w "\n%{http_code}\n%{header_x-datadog-trace-id}\n%{header_x-datadog-span-id}")

echo "Response:"
echo "$CASE_GET_RESPONSE" | head -n -3 | jq '.'

echo "Trace Headers:"
echo "$CASE_GET_RESPONSE" | tail -n 3

# Test 5: List cases with pagination
echo -e "\n5. Listing cases with pagination..."
LIST_RESPONSE=$(curl -s -X GET "$API_BASE/cases?page=1&limit=5" \
  -H "X-Request-ID: $REQUEST_ID")

echo "$LIST_RESPONSE" | jq '.'

# Test 6: Update the case
echo -e "\n6. Updating the case..."
UPDATE_RESPONSE=$(curl -s -X PUT "$API_BASE/cases/$CASE_ID" \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: $REQUEST_ID" \
  -d '{
    "caseType": "PRE_ARBITRATION",
    "primaryAccountNumber": "4111111111111111",
    "transactionAmount": 150.00,
    "transactionCurrency": "USD",
    "transactionDate": "2024-01-01T00:00:00Z",
    "transactionId": "123456789",
    "merchantName": "Updated Test Merchant",
    "merchantCategoryCode": "5411",
    "reasonCode": "10.1",
    "disputeAmount": 150.00,
    "disputeCurrency": "USD",
    "filingAs": "ISSUER",
    "filingIca": "123456",
    "filedAgainstIca": "654321",
    "filedBy": "Test User",
    "filedByContactName": "John Doe",
    "filedByContactPhone": "+1234567890",
    "filedByContactEmail": "john.doe@example.com"
  }')

echo "$UPDATE_RESPONSE" | jq '.'

# Test 7: Error handling with tracing
echo -e "\n7. Testing error handling with tracing..."
ERROR_RESPONSE=$(curl -s -X GET "$API_BASE/cases/nonexistent-id" \
  -H "X-Request-ID: $REQUEST_ID" \
  -w "\n%{http_code}")

echo "Error Response:"
echo "$ERROR_RESPONSE" | head -n -1 | jq '.'

# Test 8: CORS preflight with tracing
echo -e "\n8. Testing CORS preflight with tracing..."
CORS_RESPONSE=$(curl -s -X OPTIONS "$API_BASE/cases" \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type" \
  -H "X-Request-ID: $REQUEST_ID" \
  -w "\n%{http_code}")

echo "CORS Response Headers:"
echo "$CORS_RESPONSE" | tail -n 1

# Cleanup
echo -e "\n🧹 Cleaning up test files..."
rm -f test_document.txt

echo -e "\n✅ Datadog Integration Tests Completed!"
echo "=================================================="
echo "Check your Datadog dashboard for:"
echo "- APM traces for all requests"
echo "- Structured logs with trace correlation"
echo "- Performance metrics"
echo "- Error tracking"
echo ""
echo "Request ID used for correlation: $REQUEST_ID"
echo "Case ID: $CASE_ID"
echo "Document ID: $DOC_ID"
