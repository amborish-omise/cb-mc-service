# System Diagrams Overview

This directory contains comprehensive system diagrams that provide visual understanding of the `cb-mc-service` project architecture, workflows, and deployment processes.

## 📊 Available Diagrams

### **1. [Webhook Processing Flow](./01-webhook-processing-flow.md)**
- **Purpose**: Shows how Ethoca webhook requests flow through the system
- **Key Elements**: Request validation, payload processing, business logic, error handling
- **Use Case**: Understanding webhook processing workflow and decision points

### **2. [System Architecture](./02-system-architecture.md)**
- **Purpose**: Overview of the complete system architecture and component relationships
- **Key Elements**: Application layers, external systems, middleware, infrastructure
- **Use Case**: High-level system understanding and component interaction

### **3. [Data Flow Sequence](./03-data-flow-sequence.md)**
- **Purpose**: Detailed sequence of interactions between components during webhook processing
- **Key Elements**: Component interactions, timing, error scenarios, response generation
- **Use Case**: Understanding detailed component communication and data flow

### **4. [Project Structure](./04-project-structure.md)**
- **Purpose**: Complete project directory structure and organization
- **Key Elements**: File organization, package relationships, configuration files
- **Use Case**: Navigating the codebase and understanding project layout

### **5. [Security & Validation Flow](./05-security-validation-flow.md)**
- **Purpose**: Security measures and validation flow for incoming requests
- **Key Elements**: Input validation, security checks, error handling, response codes
- **Use Case**: Understanding security implementation and validation rules

### **6. [Monitoring & Observability](./06-monitoring-observability.md)**
- **Purpose**: How monitoring, logging, and observability are implemented
- **Key Elements**: Datadog integration, metrics collection, alerting, dashboards
- **Use Case**: Understanding monitoring setup and observability features

### **7. [Deployment & CI/CD Flow](./07-deployment-cicd-flow.md)**
- **Purpose**: Complete deployment pipeline and CI/CD workflow
- **Key Elements**: Buildkite pipeline, environment deployment, Kubernetes, monitoring
- **Use Case**: Understanding deployment process and infrastructure management

## 🎯 How to Use These Diagrams

### **For Developers**
- **Onboarding**: Start with System Architecture and Project Structure
- **Feature Development**: Reference Webhook Processing Flow and Data Flow Sequence
- **Debugging**: Use Security & Validation Flow and Monitoring & Observability

### **For DevOps Engineers**
- **Infrastructure**: Focus on Deployment & CI/CD Flow and System Architecture
- **Monitoring**: Use Monitoring & Observability and Security & Validation Flow
- **Troubleshooting**: Reference all diagrams for comprehensive system understanding

### **For Product Managers**
- **System Overview**: Start with System Architecture and Webhook Processing Flow
- **Process Understanding**: Use Data Flow Sequence and Security & Validation Flow
- **Feature Planning**: Reference Project Structure for implementation scope

### **For QA Engineers**
- **Test Planning**: Use Webhook Processing Flow and Security & Validation Flow
- **Test Coverage**: Reference Data Flow Sequence for comprehensive testing
- **Environment Setup**: Use Deployment & CI/CD Flow for test environment configuration

## 🔧 Diagram Rendering

All diagrams are written in **Mermaid** syntax and can be rendered in:

- **GitHub**: Automatically renders in markdown files
- **GitLab**: Supports Mermaid rendering
- **VS Code**: With Mermaid extension
- **Online**: [Mermaid Live Editor](https://mermaid.live/)
- **Documentation**: Many documentation platforms support Mermaid

## 📝 Updating Diagrams

When making changes to the system:

1. **Update the relevant diagram** to reflect the new architecture
2. **Maintain consistency** across all related diagrams
3. **Add new diagrams** for new major components or workflows
4. **Update this README** to include new diagrams

## 🚀 Quick Start

1. **Begin with System Architecture** for high-level understanding
2. **Review Webhook Processing Flow** for core functionality
3. **Examine Project Structure** for codebase navigation
4. **Reference specific diagrams** based on your current needs

## 📚 Related Documentation

- **Main README**: [../README.md](../README.md)
- **Webhook Guide**: [../ETHOCA_WEBHOOK.md](../ETHOCA_WEBHOOK.md)
- **Sample Payload**: [../sample-webhook-payload.json](../sample-webhook-payload.json)
- **API Specification**: [../../ethoca-alerts-merchant-api-swagger.yaml](../../ethoca-alerts-merchant-api-swagger.yaml)

---

*These diagrams are living documents that should be updated as the system evolves. For questions or suggestions about the diagrams, please refer to the project documentation or create an issue.*
