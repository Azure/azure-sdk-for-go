import type { Entity, Enum, EnumMember, Union } from "../../core/types.js";
import { DecoratorArgs } from "../utils.js";
/**
 * Describes an enum type for creation.
 */
interface EnumDescriptor {
    /**
     * The name of the enum declaration.
     */
    name: string;
    /**
     * Decorators to apply to the enum.
     */
    decorators?: DecoratorArgs[];
    /**
     * The members of the enum. If a member is an object, each property will be
     * converted to an EnumMember with the same name and value.
     */
    members?: Record<string, string | number> | EnumMember[];
}
/**
 * A kit for working with enum types.
 * @typekit enum
 */
export interface EnumKit {
    /**
     * Build an enum type. The enum type will be finished (i.e. decorators are
     * run).
     */
    create(desc: EnumDescriptor): Enum;
    /**
     * Build an equivalent enum from the given union.
     *
     *
     * @remarks
     *
     * Union variants which are
     * not valid enum members are skipped. You can check if a union is a valid
     * enum with {@link UnionKit.union}'s `isEnumValue`.
     *
     * Any API documentation will be rendered and preserved in the resulting enum.
     * - No other decorators are copied from the union to the enum
     *
     */
    createFromUnion(type: Union): Enum;
    /**
     * Check if `type` is an enum type.
     *
     * @param type the type to check.
     */
    is(type: Entity): type is Enum;
}
interface TypekitExtension {
    enum: EnumKit;
}
declare module "../define-kit.js" {
    interface Typekit extends TypekitExtension {
    }
}
export {};
//# sourceMappingURL=enum.d.ts.map