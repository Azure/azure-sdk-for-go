# Frequently Asked Questions for Azure Go Management SDK

### 1. **JSON Unmarshal Error**

When a user reports a bug with the description "JSON unmarshal error," follow these steps to resolve the issue:

- **Step 1: Verify the SDK version**
  - Check the reported SDK version. If it is not the latest, suggest the user try the latest version locally. If the latest version resolves the issue, recommend upgrading to the latest version.
  
  Example:  
  If the issue is fixed with the latest version, suggest upgrading to that version.

- **Step 2: Check the response body**
  - Suggest the user open the logger to view the response body. This will allow us to inspect the details returned by the API and determine if the issue is related to the SDK.
  
  Example:  
  If the SDK's unmarshaling result differs from the API response, the issue may lie with the SDK's handling of the data.

- **Step 3: Add 'service-attention' label**
  - If the issue persists, label the ticket with `service-attention`.

---

### 2. **Add 'Service-Attention' Label**

When a user reports issues related to product experience or functionality of the service, add the `service-attention` label. For example:

- Issues regarding service-specific behavior (e.g., an error with a particular Azure service feature).
  
Example:  
- A report about a feature in Azure's resource manager not working as expected could warrant this label.

---

### 3. **Some Cases Do Not Belong to Management SDK**

Some reports do not belong to the Azure Resource Management SDK. These include:

- **Example 1: Not related to resource management**  
  - If a report involves a different aspect of Azure, such as management groups (e.g., Azure#23895), and not resource management, the issue may fall outside the scope of the SDK.
  
- **Example 2: Namespace-related issues**  
  - Reports about namespaces starting with `az` (e.g., `Azure#23889` about `azcosmos`) should be considered as they fall outside the scope of `sdk/resourcemanager`. These should not be handled by the Azure Resource Management SDK team.

---
