import { defineKit } from "../define-kit.js";
defineKit({
    intrinsic: {
        get any() {
            return this.program.checker.anyType;
        },
        get error() {
            return this.program.checker.errorType;
        },
        get never() {
            return this.program.checker.neverType;
        },
        get null() {
            return this.program.checker.nullType;
        },
        get void() {
            return this.program.checker.voidType;
        },
        is(entity) {
            return entity.entityKind === "Type" && entity.kind === "Intrinsic";
        },
    },
});
//# sourceMappingURL=intrinsic.js.map