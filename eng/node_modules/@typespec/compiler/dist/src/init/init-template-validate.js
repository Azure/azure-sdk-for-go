import { createJSONSchemaValidator } from "../core/schema-validator.js";
import { InitTemplateSchema } from "./init-template.js";
export function validateTemplateDefinitions(template, templateName, strictValidation) {
    const validator = createJSONSchemaValidator(InitTemplateSchema, {
        strict: strictValidation,
    });
    const diagnostics = validator.validate(template, templateName);
    return { valid: diagnostics.length === 0, diagnostics };
}
//# sourceMappingURL=init-template-validate.js.map