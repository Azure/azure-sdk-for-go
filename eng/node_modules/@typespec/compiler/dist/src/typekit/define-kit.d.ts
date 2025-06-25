import { type Program } from "../core/program.js";
import { Realm } from "../experimental/realm.js";
/**
 * A Typekit is a collection of utility functions and namespaces that allow you to work with TypeSpec types.
 */
export interface Typekit {
    readonly program: Program;
    readonly realm: Realm;
}
/**
 * contextual typing to type guards is annoying (often have to restate the signature),
 * so this helper will remove the type assertions from the interface you are currently defining.
 */
export type StripGuards<T> = {
    [K in keyof T]: T[K] extends (...args: infer P) => infer R ? (...args: P) => R : T[K] extends Record<string, any> ? StripGuards<T[K]> : T[K];
};
export declare const TypekitNamespaceSymbol: unique symbol;
/**
 * Defines an extension to the Typekit interface.
 *
 * All Typekit instances will inherit the functionality defined by calls to this function.
 */
export declare function defineKit<T extends Record<string, any>>(source: StripGuards<T> & ThisType<Typekit>): void;
//# sourceMappingURL=define-kit.d.ts.map