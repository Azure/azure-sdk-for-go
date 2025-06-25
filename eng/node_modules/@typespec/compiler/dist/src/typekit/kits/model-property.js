import { getVisibilityForClass } from "../../core/visibility/core.js";
import { getEncode, getFormat } from "../../lib/decorators.js";
import { defineKit } from "../define-kit.js";
defineKit({
    modelProperty: {
        is(type) {
            return type.entityKind === "Type" && type.kind === "ModelProperty";
        },
        getEncoding(type) {
            return getEncode(this.program, type) ?? getEncode(this.program, type.type);
        },
        getFormat(type) {
            return getFormat(this.program, type) ?? getFormat(this.program, type.type);
        },
        getVisibilityForClass(property, visibilityClass) {
            return getVisibilityForClass(this.program, property, visibilityClass);
        },
        create(desc) {
            return this.program.checker.createType({
                kind: "ModelProperty",
                name: desc.name,
                type: desc.type,
                optional: desc.optional ?? false,
                decorators: [],
                defaultValue: desc.defaultValue,
            });
        },
    },
});
//# sourceMappingURL=model-property.js.map