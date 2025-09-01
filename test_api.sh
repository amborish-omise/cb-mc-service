#!/bin/bash

# MasterCom Service API Test Script
# This script demonstrates how to test the MasterCom service API endpoints

BASE_URL="http://localhost:8080"
API_BASE="$BASE_URL/api/v6"

echo "🚀 Starting MasterCom Service API Tests"
echo "======================================"

# Test 1: Health Check
echo -e "\n1. Testing Health Check..."
curl -s -X GET "$BASE_URL/health" | jq '.'

# Test 2: Create a Case
echo -e "\n2. Creating a new case..."
CASE_RESPONSE=$(curl -s -X POST "$API_BASE/cases" \
  -H "Content-Type: application/json" \
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

# Test 3: Get the Created Case
echo -e "\n3. Getting the created case..."
curl -s -X GET "$API_BASE/cases/$CASE_ID" | jq '.'

# Test 4: List All Cases
echo -e "\n4. Listing all cases..."
curl -s -X GET "$API_BASE/cases" | jq '.'

# Test 5: List Cases with Pagination
echo -e "\n5. Listing cases with pagination..."
curl -s -X GET "$API_BASE/cases?page=1&limit=5" | jq '.'

# Test 6: Upload a Document
echo -e "\n6. Uploading a document..."
# Create a temporary test file
echo "This is a test document content for the case." > test_document.txt

DOC_RESPONSE=$(curl -s -X POST "$API_BASE/documents" \
  -F "file=@test_document.txt" \
  -F "caseId=$CASE_ID" \
  -F "description=Supporting documentation for the case" \
  -F "uploadedBy=test-user")

echo "$DOC_RESPONSE" | jq '.'

# Extract document ID
DOC_ID=$(echo "$DOC_RESPONSE" | jq -r '.id')
echo "Uploaded document ID: $DOC_ID"

# Test 7: Get the Uploaded Document
echo -e "\n7. Getting the uploaded document..."
curl -s -X GET "$API_BASE/documents/$DOC_ID" | jq '.'

# Test 8: Update the Case
echo -e "\n8. Updating the case..."
curl -s -X PUT "$API_BASE/cases/$CASE_ID" \
  -H "Content-Type: application/json" \
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
  }' | jq '.'

# Test 9: Get Updated Case
echo -e "\n9. Getting the updated case..."
curl -s -X GET "$API_BASE/cases/$CASE_ID" | jq '.'

# Test 10: Error Handling - Invalid Case ID
echo -e "\n10. Testing error handling - Invalid case ID..."
curl -s -X GET "$API_BASE/cases/invalid-id" | jq '.'

# Test 11: Error Handling - Invalid JSON
echo -e "\n11. Testing error handling - Invalid JSON..."
curl -s -X POST "$API_BASE/cases" \
  -H "Content-Type: application/json" \
  -d '{"invalid": json}' | jq '.'

# Test 12: Error Handling - Missing Required Fields
echo -e "\n12. Testing error handling - Missing required fields..."
curl -s -X POST "$API_BASE/cases" \
  -H "Content-Type: application/json" \
  -d '{"caseType": "PRE_ARBITRATION"}' | jq '.'

# Test 13: CORS Preflight Request
echo -e "\n13. Testing CORS preflight request..."
curl -s -X OPTIONS "$API_BASE/cases" \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type" \
  -v

# Test 14: Delete the Document
echo -e "\n14. Deleting the document..."
curl -s -X DELETE "$API_BASE/documents/$DOC_ID" | jq '.'

# Test 15: Delete the Case
echo -e "\n15. Deleting the case..."
curl -s -X DELETE "$API_BASE/cases/$CASE_ID" | jq '.'

# Test 16: Verify Deletion
echo -e "\n16. Verifying case deletion..."
curl -s -X GET "$API_BASE/cases/$CASE_ID" | jq '.'

# Cleanup
echo -e "\n�� Cleaning up test files..."
rm -f test_document.txt

echo -e "\n✅ API Tests Completed!"
echo "======================================"
echo "All tests have been executed successfully."
echo "Check the output above for any errors or unexpected responses."
