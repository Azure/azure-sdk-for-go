import { defineKit } from "../define-kit.js";
defineKit({
    tuple: {
        is(type) {
            return type.entityKind === "Type" && type.kind === "Tuple";
        },
        create(values = []) {
            const tuple = this.program.checker.createType({
                kind: "Tuple",
                name: "Tuple",
                values,
            });
            this.program.checker.finishType(tuple);
            return tuple;
        },
    },
});
//# sourceMappingURL=tuple.js.map