# Data Flow Sequence Diagram

This diagram shows the detailed sequence of interactions between components during webhook processing.

```mermaid
sequenceDiagram
    participant Client as Ethoca Client
    participant LB as Load Balancer
    participant Server as HTTP Server
    participant Handler as Webhook Handler
    participant Service as Webhook Service
    participant Logger as Datadog Logger
    participant Tracer as Datadog Tracer
    
    Client->>LB: POST /api/v6/webhooks/ethoca
    LB->>Server: Route Request
    Server->>Handler: HandleEthocaWebhook()
    
    Handler->>Tracer: Start Span
    Handler->>Handler: Generate Request ID
    Handler->>Handler: Validate HTTP Method
    Handler->>Handler: Validate Content-Type
    
    alt Invalid Request
        Handler->>Client: Return Error Response
    else Valid Request
        Handler->>Handler: Parse JSON Payload
        Handler->>Handler: Validate Payload Structure
        
        alt Validation Failed
            Handler->>Client: Return Validation Error
        else Validation Success
            Handler->>Service: ProcessWebhook()
            
            loop For Each Outcome
                Service->>Service: Validate Outcome
                Service->>Service: Process Based on Type
                Service->>Logger: Log Processing
                Service->>Service: Create Status Update
            end
            
            Service->>Handler: Return Acknowledgment
            Handler->>Logger: Log Success
            Handler->>Tracer: Finish Span
            Handler->>Client: Return 200 OK
        end
    end
```

## Sequence Flow Explanation

### **1. Request Initiation**
- Client sends POST request to webhook endpoint
- Load balancer routes to appropriate server instance
- Server delegates to webhook handler

### **2. Request Validation**
- Handler starts Datadog tracing span
- Generates unique request ID for tracking
- Validates HTTP method (must be POST)
- Validates content-type (must be application/json)

### **3. Payload Processing**
- Parses JSON payload into structured data
- Validates payload structure and business rules
- Returns appropriate error responses for invalid data

### **4. Business Logic Processing**
- Delegates processing to webhook service
- Service iterates through each outcome
- Applies business-specific validation rules
- Processes based on outcome type (fraud, dispute, other)

### **5. Response Generation**
- Service returns processing acknowledgment
- Handler logs success/failure information
- Tracer finishes span for observability
- Returns HTTP 200 with processing results

## Key Interactions

- **Handler ↔ Service**: Business logic delegation
- **Service ↔ Logger**: Processing status logging
- **Handler ↔ Tracer**: Distributed tracing
- **All ↔ Client**: HTTP response generation

## Error Handling

- **Invalid Requests**: Immediate error responses
- **Validation Failures**: Detailed error messages
- **Processing Errors**: Graceful degradation with status updates
- **System Failures**: Proper error logging and tracing
