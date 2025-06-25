import { Realm } from "../experimental/realm.js";
export { createDiagnosable } from "./create-diagnosable.js";
export { defineKit } from "./define-kit.js";
export { UnionKit, } from "./kits/index.js";
const DEFAULT_REALM = Symbol.for("TypeSpec.Typekit.DEFAULT_TYPEKIT_REALM");
function _$(arg) {
    let realm;
    if (Object.hasOwn(arg, "projectRoot")) {
        // arg is a Program
        realm = arg[DEFAULT_REALM] ??= new Realm(arg, "default typekit realm");
    }
    else {
        // arg is a Realm
        realm = arg;
    }
    return realm.typekit;
}
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
export const $ = _$;
//# sourceMappingURL=index.js.map