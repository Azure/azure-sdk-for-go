import { createDiagnosable } from "../create-diagnosable.js";
import { defineKit } from "../define-kit.js";
defineKit({
    entity: {
        isAssignableTo: createDiagnosable(function (source, target, diagnosticTarget) {
            return this.program.checker.isTypeAssignableTo(source, target, diagnosticTarget ?? source);
        }),
        resolve: createDiagnosable(function (reference) {
            return this.program.resolveTypeOrValueReference(reference);
        }),
    },
});
//# sourceMappingURL=entity.js.map