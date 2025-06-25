import { defineCodeFix, getSourceLocation } from "../diagnostics.js";
/**
 * Quick fix that remove unused code.
 */
export function removeUnusedCodeCodeFix(node) {
    return defineCodeFix({
        id: "remove-unused-code",
        label: `Remove unused code`,
        fix: (context) => {
            const location = getSourceLocation(node);
            return context.replaceText(location, "");
        },
    });
}
//# sourceMappingURL=remove-unused-code.codefix.js.map