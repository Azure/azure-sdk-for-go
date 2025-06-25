import { $doc, getDoc } from "../../lib/decorators.js";
import { createRekeyableMap } from "../../utils/misc.js";
import { defineKit } from "../define-kit.js";
import { decoratorApplication } from "../utils.js";
defineKit({
    enum: {
        create(desc) {
            const en = this.program.checker.createType({
                kind: "Enum",
                name: desc.name,
                decorators: decoratorApplication(this, desc.decorators),
                members: createRekeyableMap(),
            });
            if (Array.isArray(desc.members)) {
                for (const member of desc.members) {
                    member.enum = en;
                    en.members.set(member.name, member);
                }
            }
            else {
                for (const [name, member] of Object.entries(desc.members ?? {})) {
                    en.members.set(name, this.enumMember.create({ name, value: member, enum: en }));
                }
            }
            this.program.checker.finishType(en);
            return en;
        },
        is(type) {
            return type.entityKind === "Type" && type.kind === "Enum";
        },
        createFromUnion(type) {
            if (!type.name) {
                throw new Error("Cannot create an enum from an anonymous union.");
            }
            const enumMembers = [];
            for (const variant of type.variants.values()) {
                if ((variant.name && typeof variant.name === "symbol") ||
                    (!this.literal.isString(variant.type) && !this.literal.isNumeric(variant.type))) {
                    continue;
                }
                const variantDoc = getDoc(this.program, variant);
                enumMembers.push(this.enumMember.create({
                    name: variant.name,
                    value: variant.type.value,
                    decorators: variantDoc ? [[$doc, variantDoc]] : undefined,
                }));
            }
            const unionDoc = getDoc(this.program, type);
            return this.enum.create({
                name: type.name,
                members: enumMembers,
                decorators: unionDoc ? [[$doc, unionDoc]] : undefined,
            });
        },
    },
});
//# sourceMappingURL=enum.js.map