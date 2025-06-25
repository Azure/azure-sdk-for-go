import type { Model, ModelProperty, Type } from "../core/types.js";
/**
 * Filters the properties of a model by removing them from the model instance if
 * a given `filter` predicate is not satisfied.
 *
 * @param model - the model to filter properties on
 * @param filter - the predicate to filter properties with
 */
export declare function filterModelPropertiesInPlace(model: Model, filter: (prop: ModelProperty) => boolean): void;
/**
 * Creates a unique symbol for storing state on objects
 * @param name The name/description of the state
 */
export declare function createStateSymbol(name: string): symbol;
/**
 * Instantiate a NameTemplate string with the properties of a source object.
 *
 * @param formatString - The template string to format. It should contain placeholders in the form of {propertyName}.
 * @param sourceObject - The object containing the properties to replace in the template string.
 * @returns The formatted string with the placeholders replaced by the corresponding property values from the source object.
 */
export declare function replaceTemplatedStringFromProperties(formatString: string, sourceObject: Type): string;
//# sourceMappingURL=utils.d.ts.map