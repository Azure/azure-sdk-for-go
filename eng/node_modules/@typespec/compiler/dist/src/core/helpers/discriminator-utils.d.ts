import { DiscriminatedOptions } from "../../../generated-defs/TypeSpec.js";
import { Discriminator } from "../intrinsic-type-state.js";
import type { Program } from "../program.js";
import { Diagnostic, Model, Type, Union } from "../types.js";
export interface DiscriminatedUnion {
    readonly options: Required<DiscriminatedOptions>;
    readonly variants: Map<string, Type>;
    readonly defaultVariant?: Type;
    readonly type: Union;
}
export interface DiscriminatedUnionLegacy {
    kind: "legacy";
    propertyName: string;
    variants: Map<string, Model>;
}
export declare function getDiscriminatedUnion(typeOrProgram: Program, typeOrDiscriminator: Union): [DiscriminatedUnion | undefined, readonly Diagnostic[]];
/**
 * Run the validation on all discriminated models to make sure the discriminator are valid.
 * This has to be done after the checker so we can have the full picture of all the dervied models.
 */
export declare function validateInheritanceDiscriminatedUnions(program: Program): void;
export declare function getDiscriminatedUnionFromInheritance(type: Model, discriminator: Discriminator): [DiscriminatedUnionLegacy, readonly Diagnostic[]];
//# sourceMappingURL=discriminator-utils.d.ts.map