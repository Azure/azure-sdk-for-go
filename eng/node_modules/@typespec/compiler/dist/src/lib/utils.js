/**
 * Filters the properties of a model by removing them from the model instance if
 * a given `filter` predicate is not satisfied.
 *
 * @param model - the model to filter properties on
 * @param filter - the predicate to filter properties with
 */
export function filterModelPropertiesInPlace(model, filter) {
    for (const [key, prop] of model.properties) {
        if (!filter(prop)) {
            model.properties.delete(key);
        }
    }
}
/**
 * Creates a unique symbol for storing state on objects
 * @param name The name/description of the state
 */
export function createStateSymbol(name) {
    return Symbol.for(`TypeSpec.${name}`);
}
/**
 * Instantiate a NameTemplate string with the properties of a source object.
 *
 * @param formatString - The template string to format. It should contain placeholders in the form of {propertyName}.
 * @param sourceObject - The object containing the properties to replace in the template string.
 * @returns The formatted string with the placeholders replaced by the corresponding property values from the source object.
 */
export function replaceTemplatedStringFromProperties(formatString, sourceObject) {
    // Template parameters are not valid source objects, just skip them
    if (sourceObject.kind === "TemplateParameter") {
        return formatString;
    }
    return formatString.replace(/{(\w+)}/g, (_, propName) => {
        return sourceObject[propName];
    });
}
//# sourceMappingURL=utils.js.map