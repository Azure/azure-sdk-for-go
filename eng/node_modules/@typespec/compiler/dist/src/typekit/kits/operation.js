import { getPagingOperation } from "../../lib/paging.js";
import { createDiagnosable } from "../create-diagnosable.js";
import { defineKit } from "../define-kit.js";
defineKit({
    operation: {
        is(type) {
            return type.entityKind === "Type" && type.kind === "Operation";
        },
        getPagingMetadata: createDiagnosable(function (operation) {
            return getPagingOperation(this.program, operation);
        }),
        create(desc) {
            const parametersModel = this.model.create({
                name: `${desc.name}Parameters`,
                properties: desc.parameters.reduce((acc, property) => {
                    acc[property.name] = property;
                    return acc;
                }, {}),
            });
            const operation = this.program.checker.createType({
                kind: "Operation",
                name: desc.name,
                decorators: [],
                parameters: parametersModel,
                returnType: desc.returnType,
            });
            this.program.checker.finishType(operation);
            return operation;
        },
    },
});
//# sourceMappingURL=operation.js.map