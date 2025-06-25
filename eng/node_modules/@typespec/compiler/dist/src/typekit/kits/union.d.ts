import { DiscriminatedUnion } from "../../core/helpers/discriminator-utils.js";
import type { Entity, Enum, Type, Union, UnionVariant } from "../../core/types.js";
import { Diagnosable } from "../create-diagnosable.js";
import { DecoratorArgs } from "../utils.js";
/**
 * A descriptor for a union type.
 */
export interface UnionDescriptor {
    /**
     * The name of the union. If name is provided, it is a union declaration.
     * Otherwise, it is a union expression.
     */
    name?: string;
    /**
     * Decorators to apply to the union.
     */
    decorators?: DecoratorArgs[];
    /**
     * The variants of the union. If a variant is a string, number, or boolean, it
     * will be converted to a union variant with the same name and type.
     */
    variants?: Record<string | symbol, string | number> | UnionVariant[];
}
/**
 * Utilities for working with unions.
 * @typekit union
 */
export interface UnionKit {
    /**
     * Creates a union type with filtered variants.
     * @param filterFn Function to filter the union variants
     */
    filter(union: Union, filterFn: (variant: UnionVariant) => boolean): Union;
    /**
     * Create a union type.
     *
     * @param desc The descriptor of the union.
     */
    create(desc: UnionDescriptor): Union;
    /**
     * Create an anonymous union type from an array of types.
     *
     * @param children The types to create a union from.
     *
     * Any API documentation will be rendered and preserved in the resulting union.
     *
     * No other decorators are copied from the enum to the union.
     */
    create(children: Type[]): Union;
    /**
     * Creates a union type from an enum.
     *
     * @remarks
     *
     * @param type The enum to create a union from.
     *
     * For member without an explicit value, the member name is used as the value.
     *
     * Any API documentation will be rendered and preserved in the resulting union.
     *
     * No other decorators are copied from the enum to the union.
     */
    createFromEnum(type: Enum): Union;
    /**
     * Check if the given `type` is a union.
     *
     * @param type The type to check.
     */
    is(type: Entity): type is Union;
    /**
     * Check if the union is a valid enum. Specifically, this checks if the
     * union has a name (since there are no enum expressions), and whether each
     * of the variant types is a valid enum member value.
     *
     * @param type The union to check.
     */
    isValidEnum(type: Union): boolean;
    /**
     * Check if a union is extensible. Extensible unions are unions which contain a variant
     * that is a supertype of all the other types. This means that the subtypes of the common
     * supertype are known example values, but others may be present.
     * @param type The union to check.
     */
    isExtensible(type: Union): boolean;
    /**
     * Checks if an union is an expression (anonymous) or declared.
     * @param type Uniton to check if it is an expression
     */
    isExpression(type: Union): boolean;
    /**
     * Resolves a discriminated union for the given union.
     * @param type Union to resolve the discriminated union for.
     */
    getDiscriminatedUnion: Diagnosable<(type: Union) => DiscriminatedUnion | undefined>;
}
interface TypekitExtension {
    /**
     * Utilities for working with unions.
     */
    union: UnionKit;
}
declare module "../define-kit.js" {
    interface Typekit extends TypekitExtension {
    }
}
export declare const UnionKit: void;
export {};
//# sourceMappingURL=union.d.ts.map