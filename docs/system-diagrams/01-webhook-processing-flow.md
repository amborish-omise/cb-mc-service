# Webhook Processing Flow Diagram

This diagram shows the complete flow of how Ethoca webhook requests are processed through the system.

```mermaid
flowchart TD
    A[Ethoca Webhook Request] --> B{Validate Request}
    B -->|Invalid Method| C[Return 405 Method Not Allowed]
    B -->|Invalid Content-Type| D[Return 400 Bad Request]
    B -->|Valid Request| E[Parse JSON Payload]
    
    E -->|Parse Error| F[Return 400 Invalid JSON]
    E -->|Valid JSON| G{Validate Payload Structure}
    
    G -->|No Outcomes| H[Return 400 No Outcomes]
    G -->|Too Many Outcomes| I[Return 400 Too Many Outcomes]
    G -->|Valid Structure| J[Process Each Outcome]
    
    J --> K[Generate Request ID]
    K --> L[Create Webhook Event]
    L --> M{Validate Outcome}
    
    M -->|Invalid| N[Mark as Failed]
    M -->|Valid| O{Determine Outcome Type}
    
    O -->|Fraud Alert| P[Process Fraud Outcome]
    O -->|Customer Dispute| Q[Process Dispute Outcome]
    O -->|Other| R[Process Other Outcome]
    
    P --> S[Update Fraud Database]
    Q --> T[Process Refund/Resolution]
    R --> U[Route for Review]
    
    S --> V[Mark as Success]
    T --> V
    U --> V
    N --> W[Create Error Status]
    
    V --> X[Add to Status Updates]
    W --> X
    X --> Y{More Outcomes?}
    
    Y -->|Yes| J
    Y -->|No| Z[Generate Acknowledgment]
    Z --> AA[Log Processing Complete]
    AA --> BB[Return 200 Success Response]
    
    style A fill:#e1f5fe
    style BB fill:#c8e6c9
    style C fill:#ffcdd2
    style D fill:#ffcdd2
    style F fill:#ffcdd2
    style H fill:#ffcdd2
    style I fill:#ffcdd2
    style N fill:#ffcdd2
```

## Key Decision Points

1. **Request Validation**: HTTP method and content-type validation
2. **Payload Structure**: JSON parsing and basic structure validation
3. **Business Rules**: Outcome-specific validation and processing
4. **Error Handling**: Comprehensive error tracking and status updates
5. **Response Generation**: Acknowledgment with processing results

## Error Scenarios

- **405 Method Not Allowed**: Non-POST requests
- **400 Bad Request**: Invalid content-type or JSON
- **400 Validation Errors**: Business rule violations
- **500 Internal Server Error**: Processing failures
