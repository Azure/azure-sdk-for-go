import type { Program } from "../core/program.js";
import { Realm } from "../experimental/realm.js";
import { type Typekit } from "./define-kit.js";
export { createDiagnosable, Diagnosable } from "./create-diagnosable.js";
export { defineKit, type Typekit } from "./define-kit.js";
export { ArrayKit, BuiltinKit, EntityKit, EnumKit, EnumMemberDescriptor, EnumMemberKit, IntrinsicKit, LiteralKit, ModelDescriptor, ModelKit, ModelPropertyDescriptor, ModelPropertyKit, OperationDescriptor, OperationKit, RecordKit, ScalarKit, TupleKit, TypeTypekit, UnionDescriptor, UnionKit, UnionVariantDescriptor, UnionVariantKit, ValueKit, } from "./kits/index.js";
declare function _$(realm: Realm): Typekit;
declare function _$(program: Program): Typekit;
/**
 * Typekit - Utilities for working with TypeSpec types.
 *
 * Each typekit is associated with a Realm in which it operates.
 *
 * You can get the typekit associated with that realm by calling
 * `$` with the realm as an argument, or by calling `$` with a program
 * as an argument (in this case, it will use that program's default
 * typekit realm or create one if it does not already exist).
 *
 * @example
 * ```ts
 * import { Realm } from "@typespec/compiler/experimental";
 * import { $ } from "@typespec/compiler/typekit";
 *
 * const realm = new Realm(program, "my custom realm");
 *
 * const clone = $(realm).type.clone(inputType);
 * ```
 *
 * @example
 * ```ts
 * import { $ } from "@typespec/compiler/typekit";
 *
 * const clone = $(program).type.clone(inputType);
 * ```
 *
 * @see {@link Realm}
 */
export declare const $: typeof _$;
//# sourceMappingURL=index.d.ts.map