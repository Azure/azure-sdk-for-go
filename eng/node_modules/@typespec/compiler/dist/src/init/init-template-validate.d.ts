import { Diagnostic, NoTarget, SourceFile } from "../index.js";
export type ValidationResult = {
    valid: boolean;
    diagnostics: readonly Diagnostic[];
};
export declare function validateTemplateDefinitions(template: unknown, templateName: SourceFile | typeof NoTarget, strictValidation: boolean): ValidationResult;
//# sourceMappingURL=init-template-validate.d.ts.map