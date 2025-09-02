# Ethoca Webhook Integration

This document describes the Ethoca webhook integration for the MasterCom service, which allows merchants and partners to submit outcomes to fraud alerts and customer disputes.

## Overview

The Ethoca webhook integration follows the [Ethoca Alerts Merchant API specification](https://developer.mastercard.com/ethoca-alerts) and enables:

- Processing of fraud alert outcomes
- Handling of customer dispute resolutions
- Automatic refund processing
- Fraud prevention and chargeback reduction

## API Endpoints

### Webhook Endpoint
```
POST /api/v6/webhooks/ethoca
```

**Content-Type:** `application/json`

**Description:** Receives webhook payloads from Ethoca containing alert outcomes.

### Health Check
```
GET /api/v6/webhooks/ethoca/health
```

**Description:** Returns the health status and configuration of the webhook endpoint.

### Statistics
```
GET /api/v6/webhooks/ethoca/stats
```

**Description:** Returns processing statistics for webhook events.

## Webhook Payload Structure

The webhook expects a JSON payload with the following structure:

```json
{
  "outcomes": [
    {
      "alertId": "A4IM9K2MIYL9F2BPF9TWUIXTU",
      "outcome": "STOPPED",
      "refundStatus": "NOT_REFUNDED",
      "refund": {
        "amount": {
          "value": 361.56,
          "currencyCode": "USD"
        },
        "type": "REFUND",
        "timestamp": "2021-06-18T22:11:05+05:00",
        "transactionId": "23aer543245678984ew39awse0",
        "acquirerReferenceNumber": "98765432456789876345213"
      },
      "amountStopped": {
        "value": 361.56,
        "currencyCode": "USD"
      },
      "comments": "Order stopped due to suspected fraud",
      "actionTimestamp": "2021-06-18T22:11:05+05:00"
    }
  ]
}
```

### Supported Outcomes

#### Confirmed Fraud Outcomes
- `STOPPED` - The order was stopped
- `PARTIALLY_STOPPED` - Part of the order was stopped
- `PREVIOUSLY_CANCELLED` - The transaction was already canceled
- `MISSED` - Too late, the order has shipped / service consumed
- `NOT_FOUND` - The order could not be found
- `ACCOUNT_SUSPENDED` - The account has been suspended
- `OTHER` - Anything else not covered above

#### Customer Dispute Outcomes
- `RESOLVED` - Case resolved with the customer
- `RESOLVED_PREVIOUSLY_REFUNDED` - Refund already processed
- `UNRESOLVED_DISPUTE` - Merchant disagrees with reason for dispute
- `NOT_FOUND` - Alert could not be found in the system
- `OTHER` - Any other outcome as described in the comments

### Refund Status Values
- `REFUNDED` - Transaction was refunded
- `NOT_REFUNDED` - Transaction was not refunded
- `NOT_SETTLED` - Transaction did not go to settlement

## Response Format

### Success Response (200 OK)
```json
{
  "outcomeResponses": [
    {
      "alertId": "A4IM9K2MIYL9F2BPF9TWUIXTU",
      "status": "SUCCESS"
    }
  ]
}
```

### Error Response (400 Bad Request)
```json
{
  "error": "No outcomes provided in webhook payload",
  "code": "NO_OUTCOMES"
}
```

### Error Response (500 Internal Server Error)
```json
{
  "error": "Failed to process webhook",
  "code": "PROCESSING_ERROR"
}
```

## Configuration

The webhook can be configured using environment variables:

```bash
# Webhook endpoint configuration
ETHOCA_WEBHOOK_ENDPOINT=/api/v6/webhooks/ethoca
ETHOCA_WEBHOOK_SECRET_KEY=your-secret-key-here

# Processing configuration
ETHOCA_WEBHOOK_TIMEOUT=30
ETHOCA_WEBHOOK_MAX_RETRIES=3
ETHOCA_WEBHOOK_BATCH_SIZE=25
```

## Validation Rules

1. **Payload Structure**: Must contain at least one outcome
2. **Outcome Count**: Maximum of 25 outcomes per webhook
3. **Alert ID**: Must be exactly 25 characters
4. **Amount Validation**: 
   - Refund amount must be > 0 when refund status is `REFUNDED`
   - Amount stopped must be > 0 when outcome is `STOPPED`
5. **Timestamp Format**: ISO 8601 format (e.g., `2021-06-18T22:11:05+05:00`)

## Processing Flow

1. **Receive Webhook**: Validate HTTP method and content type
2. **Parse Payload**: Parse JSON and validate structure
3. **Process Outcomes**: Process each outcome based on type
4. **Business Logic**: Apply fraud/dispute processing rules
5. **Generate Response**: Return acknowledgment with status updates
6. **Logging**: Comprehensive logging with Datadog integration

## Security Features

- **Request ID Tracking**: Each webhook request gets a unique ID
- **Input Validation**: Comprehensive validation of all input fields
- **Error Handling**: Graceful error handling with detailed error messages
- **Rate Limiting**: Built-in support for rate limiting (configurable)
- **Audit Logging**: All webhook events are logged for audit purposes

## Testing

### Sample Payload
Use the sample payload in `docs/sample-webhook-payload.json` for testing:

```bash
curl -X POST http://localhost:8080/api/v6/webhooks/ethoca \
  -H "Content-Type: application/json" \
  -d @docs/sample-webhook-payload.json
```

### Health Check
```bash
curl http://localhost:8080/api/v6/webhooks/ethoca/health
```

### Statistics
```bash
curl http://localhost:8080/api/v6/webhooks/ethoca/stats
```

## Monitoring and Observability

### Datadog Integration
- **Tracing**: Each webhook request is traced with spans
- **Logging**: Structured logging with correlation IDs
- **Metrics**: Processing time and success rate tracking
- **Error Tracking**: Automatic error capture and reporting

### Key Metrics
- Webhook processing time
- Success/failure rates
- Outcome type distribution
- Error frequency by type

## Error Handling

The service handles various error scenarios gracefully:

1. **Invalid JSON**: Returns 400 with parsing error details
2. **Validation Errors**: Returns 400 with specific validation messages
3. **Processing Errors**: Returns 500 with error codes
4. **Partial Failures**: Individual outcomes can fail while others succeed

## Best Practices

1. **Response Time**: Respond to alerts within 24 hours for optimal chargeback prevention
2. **Batch Processing**: Use batch processing for multiple outcomes (max 25)
3. **Error Handling**: Implement retry logic for failed webhooks
4. **Monitoring**: Set up alerts for webhook processing failures
5. **Security**: Use HTTPS and validate webhook signatures in production

## Integration Examples

### Fraud Alert Processing
```json
{
  "outcomes": [{
    "alertId": "A4IM9K2MIYL9F2BPF9TWUIXTU",
    "outcome": "STOPPED",
    "refundStatus": "NOT_REFUNDED",
    "refund": {
      "amount": {"value": 100.00, "currencyCode": "USD"},
      "timestamp": "2021-06-18T22:11:05+05:00"
    },
    "amountStopped": {"value": 100.00, "currencyCode": "USD"}
  }]
}
```

### Dispute Resolution
```json
{
  "outcomes": [{
    "alertId": "B5JN0L3NJZM0G3CQG0UXVJYUV",
    "outcome": "RESOLVED",
    "refundStatus": "REFUNDED",
    "refund": {
      "amount": {"value": 50.00, "currencyCode": "USD"},
      "timestamp": "2021-06-18T22:15:30+05:00"
    },
    "amountStopped": {"value": 0.00, "currencyCode": "USD"},
    "comments": "Customer dispute resolved with partial refund"
  }]
}
```

## Support

For technical support or questions about the webhook integration:

- **API Documentation**: [Ethoca Alerts Merchant API](https://developer.mastercard.com/ethoca-alerts)
- **Mastercard Developer Support**: apisupport@mastercard.com
- **Service Logs**: Check Datadog logs for detailed processing information
