import { defineKit } from "../define-kit.js";
import { decoratorApplication } from "../utils.js";
defineKit({
    enumMember: {
        create(desc) {
            const member = this.program.checker.createType({
                kind: "EnumMember",
                name: desc.name,
                value: desc.value,
                decorators: decoratorApplication(this, desc.decorators),
                enum: desc.enum, // initialized in enum.build if not provided here
            });
            this.program.checker.finishType(member);
            return member;
        },
        is(type) {
            return type.entityKind === "Type" && type.kind === "EnumMember";
        },
    },
});
//# sourceMappingURL=enum-member.js.map