# Security & Validation Flow Diagram

This diagram shows the comprehensive security and validation flow for incoming webhook requests.

```mermaid
flowchart TD
    A[Incoming Webhook] --> B[Request Validation]
    
    subgraph "Request Validation"
        B --> C{HTTP Method = POST?}
        C -->|No| D[Return 405]
        C -->|Yes| E{Content-Type = application/json?}
        E -->|No| F[Return 400]
        E -->|Yes| G[Generate Request ID]
    end
    
    subgraph "Payload Validation"
        G --> H[Parse JSON]
        H --> I{Valid JSON?}
        I -->|No| J[Return 400 Invalid JSON]
        I -->|Yes| K{Has Outcomes?}
        K -->|No| L[Return 400 No Outcomes]
        K -->|Yes| M{Outcome Count <= 25?}
        M -->|No| N[Return 400 Too Many Outcomes]
        M -->|Yes| O[Validate Each Outcome]
    end
    
    subgraph "Business Validation"
        O --> P{Alert ID Length = 25?}
        P -->|No| Q[Mark as Failed]
        P -->|Yes| R{Outcome Length Valid?}
        R -->|No| Q
        R -->|Yes| S{Refund Status Valid?}
        S -->|No| Q
        S -->|Yes| T{Amount Validation}
        T -->|Invalid| Q
        T -->|Valid| U[Process Outcome]
    end
    
    subgraph "Processing & Response"
        U --> V[Generate Status Update]
        V --> W[Create Acknowledgment]
        W --> X[Return 200 Success]
        Q --> Y[Create Error Status]
        Y --> W
    end
    
    style A fill:#e1f5fe
    style X fill:#c8e6c9
    style D fill:#ffcdd2
    style F fill:#ffcdd2
    style J fill:#ffcdd2
    style L fill:#ffcdd2
    style N fill:#ffcdd2
    style Q fill:#ffcdd2
```

## Validation Layers

### **1. Request Validation**
- **HTTP Method**: Only POST requests allowed
- **Content-Type**: Must be `application/json`
- **Request ID**: Unique identifier for tracking

### **2. Payload Validation**
- **JSON Syntax**: Valid JSON structure required
- **Outcome Presence**: Must contain at least one outcome
- **Outcome Count**: Maximum 25 outcomes per request

### **3. Business Validation**
- **Alert ID**: Must be exactly 25 characters
- **Outcome**: Must be 5-30 characters
- **Refund Status**: Must be 8-12 characters
- **Amount Validation**: Currency and value validation

## Security Measures

### **Input Sanitization**
- JSON payload parsing with error handling
- String length validation to prevent buffer overflow
- Content-type validation to prevent MIME confusion

### **Rate Limiting**
- Maximum 25 outcomes per request
- Request size limits enforced
- Processing timeouts implemented

### **Error Handling**
- Detailed error messages for debugging
- Proper HTTP status codes
- Error logging for security monitoring

### **Request Tracking**
- Unique request ID generation
- Comprehensive logging and tracing
- Audit trail for compliance

## Response Codes

- **200 OK**: Successful processing
- **400 Bad Request**: Validation failures
- **405 Method Not Allowed**: Invalid HTTP method
- **500 Internal Server Error**: Processing failures

## Validation Rules

### **Alert ID**
- Required field
- Exact length: 25 characters
- Alphanumeric validation

### **Outcome**
- Required field
- Length: 5-30 characters
- Predefined values validation

### **Refund Status**
- Required field
- Length: 8-12 characters
- Status enum validation

### **Amount Fields**
- Required numeric values
- Currency code validation
- Decimal precision handling
