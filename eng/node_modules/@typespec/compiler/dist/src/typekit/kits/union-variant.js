import { defineKit } from "../define-kit.js";
import { decoratorApplication } from "../utils.js";
defineKit({
    unionVariant: {
        create(desc) {
            const variant = this.program.checker.createType({
                kind: "UnionVariant",
                name: desc.name ?? Symbol("name"),
                decorators: decoratorApplication(this, desc.decorators),
                type: desc.type,
                union: desc.union,
            });
            this.program.checker.finishType(variant);
            return variant;
        },
        is(type) {
            return type.entityKind === "Type" && type.kind === "UnionVariant";
        },
    },
});
//# sourceMappingURL=union-variant.js.map