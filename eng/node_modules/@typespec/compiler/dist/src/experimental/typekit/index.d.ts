import { type Typekit } from "../../typekit/define-kit.js";
import { Realm } from "../realm.js";
/**
 * Create a new Typekit that operates in the given realm.
 *
 * Ordinarily, you should use the default typekit `$` to manipulate types in the current program, or call `$` with a
 * Realm or Program as the first argument if you want to work in a specific realm or in the default typekit realm of
 * a specific program.
 *
 * @param realm - The realm to create the typekit in.
 *
 * @experimental
 */
export declare function createTypekit(realm: Realm): Typekit;
//# sourceMappingURL=index.d.ts.map