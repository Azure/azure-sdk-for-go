/**
 * The prototype object for Typekit instances.
 *
 * @see {@link defineKit}
 *
 * @internal
 */
export const TypekitPrototype = {};
export const TypekitNamespaceSymbol = Symbol.for("TypekitNamespace");
/**
 * Defines an extension to the Typekit interface.
 *
 * All Typekit instances will inherit the functionality defined by calls to this function.
 */
export function defineKit(source) {
    for (const [name, fnOrNs] of Object.entries(source)) {
        let kits = fnOrNs;
        if (TypekitPrototype[name] !== undefined) {
            kits = { ...TypekitPrototype[name], ...fnOrNs };
        }
        // Tag top-level namespace objects with the symbol
        if (typeof kits === "object" && kits !== null) {
            Object.defineProperty(kits, TypekitNamespaceSymbol, {
                value: true,
                enumerable: false, // Keep the symbol non-enumerable
                configurable: false,
            });
        }
        TypekitPrototype[name] = kits;
    }
}
//# sourceMappingURL=define-kit.js.map