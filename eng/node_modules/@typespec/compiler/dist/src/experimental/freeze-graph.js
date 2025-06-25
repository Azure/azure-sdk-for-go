import { navigateProgram } from "../core/semantic-walker.js";
export function freezeGraph(program) {
    function freeze(type) {
        Object.freeze(type);
    }
    navigateProgram(program, {
        templateParameter: freeze,
        scalar: freeze,
        model: freeze,
        modelProperty: freeze,
        interface: freeze,
        enum: freeze,
        enumMember: freeze,
        namespace: freeze,
        operation: freeze,
        string: freeze,
        number: freeze,
        boolean: freeze,
        tuple: freeze,
        union: freeze,
        unionVariant: freeze,
        intrinsic: freeze,
    });
}
//# sourceMappingURL=freeze-graph.js.map