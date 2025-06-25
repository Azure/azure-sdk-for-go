import { defineCodeFix, getSourceLocation } from "../diagnostics.js";
import { SyntaxKind, } from "../types.js";
/**
 * Quick fix that convert a tuple to an array value.
 */
export function createTupleToArrayValueCodeFix(node) {
    return defineCodeFix({
        id: "tuple-to-array-value",
        label: `Convert to an array value \`#[]\``,
        fix: (context) => {
            const result = [];
            addCreatedCodeFixResult(node);
            createChildTupleToArrValCodeFix(node, addCreatedCodeFixResult);
            return result;
            function addCreatedCodeFixResult(node) {
                const location = getSourceLocation(node);
                result.push(context.prependText(location, "#"));
            }
        },
    });
}
/**
 * Quick fix that convert a model expression to an object value.
 */
export function createModelToObjectValueCodeFix(node) {
    return defineCodeFix({
        id: "model-to-object-value",
        label: `Convert to an object value \`#{}\``,
        fix: (context) => {
            const result = [];
            addCreatedCodeFixResult(node);
            createChildModelToObjValCodeFix(node, addCreatedCodeFixResult);
            return result;
            function addCreatedCodeFixResult(node) {
                const location = getSourceLocation(node);
                result.push(context.prependText(location, "#"));
            }
        },
    });
}
function createChildTupleToArrValCodeFix(node, addCreatedCodeFixResult) {
    for (const childNode of node.values) {
        if (childNode.kind === SyntaxKind.ModelExpression) {
            addCreatedCodeFixResult(childNode);
            createChildModelToObjValCodeFix(childNode, addCreatedCodeFixResult);
        }
        else if (childNode.kind === SyntaxKind.TupleExpression) {
            addCreatedCodeFixResult(childNode);
            createChildTupleToArrValCodeFix(childNode, addCreatedCodeFixResult);
        }
    }
}
function createChildModelToObjValCodeFix(node, addCreatedCodeFixResult) {
    for (const prop of node.properties.values()) {
        if (prop.kind === SyntaxKind.ModelProperty) {
            const childNode = prop.value;
            if (childNode.kind === SyntaxKind.ModelExpression) {
                addCreatedCodeFixResult(childNode);
                createChildModelToObjValCodeFix(childNode, addCreatedCodeFixResult);
            }
            else if (childNode.kind === SyntaxKind.TupleExpression) {
                addCreatedCodeFixResult(childNode);
                createChildTupleToArrValCodeFix(childNode, addCreatedCodeFixResult);
            }
        }
    }
}
//# sourceMappingURL=convert-to-value.codefix.js.map