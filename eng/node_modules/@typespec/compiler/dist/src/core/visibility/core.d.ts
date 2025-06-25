import type { Program } from "../program.js";
import type { DecoratorContext, Enum, EnumMember, ModelProperty } from "../types.js";
import type { VisibilityFilter as GeneratedVisibilityFilter } from "../../../generated-defs/TypeSpec.js";
export { GeneratedVisibilityFilter };
/**
 * Set the default visibility modifier set for a visibility class.
 *
 * This function may only be called ONCE per visibility class and must be called
 * before the default modifier set is used by any operation.
 */
export declare function setDefaultModifierSetForVisibilityClass(program: Program, visibilityClass: Enum, defaultSet: Set<EnumMember>): void;
/**
 * Check if a property has had its visibility modifiers sealed.
 *
 * If the property has been sealed globally, this function will return true. If the property has been sealed for the
 * given visibility class, this function will return true.
 *
 * Otherwise, this function returns false.
 *
 * @param property - the property to check
 * @param visibilityClass - the optional visibility class to check
 * @returns true if the property is sealed for the given visibility class, false otherwise
 */
export declare function isSealed(program: Program, property: ModelProperty, visibilityClass?: Enum): boolean;
/**
 * Seals a property's visibility modifiers.
 *
 * If the `visibilityClass` is provided, the property's visibility modifiers will be sealed for that visibility class
 * only. Otherwise, the property's visibility modifiers will be sealed for all visibility classes (globally).
 *
 * @param property - the property to seal
 * @param visibilityClass - the optional visibility class to seal the property for
 */
export declare function sealVisibilityModifiers(program: Program, property: ModelProperty, visibilityClass?: Enum): void;
/**
 * Seals a program's visibility modifiers.
 *
 * This affects all properties in the program and prevents any further modifications to visibility modifiers within the
 * program.
 *
 * Once the modifiers for a program are sealed, they cannot be unsealed.
 *
 * @param program - the program to seal
 */
export declare function sealVisibilityModifiersForProgram(program: Program): void;
/**
 * Add visibility modifiers to a property.
 *
 * This function will add all the `modifiers` to the active set of visibility modifiers for the given `property`.
 *
 * If no set of active modifiers exists for the given `property`, an empty set will be created for the property.
 *
 * If the visibility modifiers for `property` in the given modifier's visibility class have been sealed, this function
 * will issue a diagnostic and ignore that modifier, but it will still add the rest of the modifiers whose classes have
 * not been sealed.
 *
 * @param program - the program in which the ModelProperty occurs
 * @param property - the property to add visibility modifiers to
 * @param modifiers - the visibility modifiers to add
 * @param context - the optional decorator context to use for displaying diagnostics
 */
export declare function addVisibilityModifiers(program: Program, property: ModelProperty, modifiers: EnumMember[], context?: DecoratorContext): void;
/**
 * Remove visibility modifiers from a property.
 *
 * This function will remove all the `modifiers` from the active set of visibility modifiers for the given `property`.
 *
 * If no set of active modifiers exists for the given `property`, the default set for the modifier's visibility class
 * will be used.
 *
 * If the visibility modifiers for `property` in the given modifier's visibility class have been sealed, this function
 * will issue a diagnostic and ignore that modifier, but it will still remove the rest of the modifiers whose classes
 * have not been sealed.
 *
 * @param program - the program in which the ModelProperty occurs
 * @param property - the property to remove visibility modifiers from
 * @param modifiers - the visibility modifiers to remove
 * @param context - the optional decorator context to use for displaying diagnostics
 */
export declare function removeVisibilityModifiers(program: Program, property: ModelProperty, modifiers: EnumMember[], context?: DecoratorContext): void;
/**
 * Clears the visibility modifiers for a property in a given visibility class.
 *
 * If the visibility modifiers for the given class are sealed, this function will issue a diagnostic and leave the
 * visibility modifiers unchanged.
 *
 * @param program - the program in which the ModelProperty occurs
 * @param property - the property to clear visibility modifiers for
 * @param visibilityClass - the visibility class to clear visibility modifiers for
 * @param context - the optional decorator context to use for displaying diagnostics
 */
export declare function clearVisibilityModifiersForClass(program: Program, property: ModelProperty, visibilityClass: Enum, context?: DecoratorContext): void;
/**
 * Resets the visibility modifiers for a property in a given visibility class.
 *
 * This does not clear the modifiers. It resets them to the _uninitialized_ state.
 *
 * This is useful when cloning properties and you want to reset the visibility modifiers on the clone.
 *
 * If the visibility modifiers for this property and given visibility class are sealed, this function will issue a
 * diagnostic and leave the visibility modifiers unchanged.
 *
 * @param program - the program in which the property occurs
 * @param property - the property to reset visibility modifiers for
 * @param visibilityClass - the visibility class to reset visibility modifiers for
 * @param context - the optional decorator context to use for displaying diagnostics
 */
export declare function resetVisibilityModifiersForClass(program: Program, property: ModelProperty, visibilityClass: Enum, context?: DecoratorContext): void;
/**
 * Returns the active visibility modifiers for a property in a given visibility class.
 *
 * This function is infallible. If the visibility modifiers for the given class have not been set explicitly, it will
 * return the default visibility modifiers for the class.
 *
 * @param program - the program in which the property occurs
 * @param property - the property to get visibility modifiers for
 * @param visibilityClass - the visibility class to get visibility modifiers for
 * @returns the set of active modifiers (enum members) for the property and visibility class
 */
export declare function getVisibilityForClass(program: Program, property: ModelProperty, visibilityClass: Enum): Set<EnumMember>;
/**
 * Determines if a property has a specified visibility modifier.
 *
 * If no visibility modifiers have been set for the visibility class of the modifier, the visibility class's default
 * modifier set is used.
 *
 * @param program - the program in which the property occurs
 * @param property - the property to check
 * @param modifier - the visibility modifier to check for
 * @returns true if the property has the visibility modifier, false otherwise
 */
export declare function hasVisibility(program: Program, property: ModelProperty, modifier: EnumMember): boolean;
/**
 * A visibility filter that can be used to determine if a property is visible.
 *
 * The filter is defined by three sets of visibility modifiers. The filter is satisfied if the property has:
 *
 * - ALL of the visibilities in the `all` set.
 *
 * AND
 *
 * - ANY of the visibilities in the `any` set.
 *
 * AND
 *
 * - NONE of the visibilities in the `none` set.
 *
 * Note: The constraints behave similarly to the `every` and `some` methods of the Array prototype in JavaScript. If the
 * `any` constraint is set to an empty set, it will _NEVER_ be satisfied (similarly, `Array.prototype.some` will always
 * return `false` for an empty array). If the `none` constraint is set to an empty set, it will _ALWAYS_ be satisfied.
 * If the `all` constraint is set to an empty set, it will be satisfied (similarly, `Array.prototype.every` will always
 * return `true` for an empty array).
 *
 */
export interface VisibilityFilter {
    /**
     * If set, the filter considers a property visible if it has ALL of these visibility modifiers.
     *
     * If this set is empty, the filter will be satisfied if the other constraints are satisfied.
     */
    all?: Set<EnumMember>;
    /**
     * If set, the filter considers a property visible if it has ANY of these visibility modifiers.
     *
     * If this set is empty, the filter will _NEVER_ be satisfied.
     */
    any?: Set<EnumMember>;
    /**
     * If set, the filter considers a property visible if it has NONE of these visibility modifiers.
     *
     * If this set is empty, the filter will be satisfied if the other constraints are satisfied.
     */
    none?: Set<EnumMember>;
}
export declare const VisibilityFilter: {
    /**
     * Convert a TypeSpec `GeneratedVisibilityFilter` value to a `VisibilityFilter`.
     *
     * @param filter - the decorator argument filter to convert
     * @returns a `VisibilityFilter` object that can be consumed by the visibility APIs
     */
    fromDecoratorArgument(filter: GeneratedVisibilityFilter): VisibilityFilter;
    /**
     * Extracts the unique visibility classes referred to by the modifiers in a
     * visibility filter.
     *
     * @param filter - the visibility filter to extract visibility classes from
     * @returns a set of visibility classes referred to by the filter
     */
    getVisibilityClasses(filter: VisibilityFilter): Set<Enum>;
    /**
     * Converts a visibility filter into a stable string representation.
     *
     * This can be used as a cache key for the filter that will be stable for filters that are not object-identical but
     * are semantically identical.
     */
    toCacheKey(program: Program, filter: VisibilityFilter): string;
};
/**
 * Determines if a property is visible according to the given visibility filter.
 *
 * @see VisibilityFilter
 *
 * @param program - the program in which the property occurs
 * @param property - the property to check
 * @param filter - the visibility filter to use
 * @returns true if the property is visible according to the filter, false otherwise
 */
export declare function isVisible(program: Program, property: ModelProperty, filter: VisibilityFilter): boolean;
//# sourceMappingURL=core.d.ts.map