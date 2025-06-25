import type { Entity, Enum, EnumMember, ModelProperty, Type, Value } from "../../core/types.js";
import { EncodeData } from "../../lib/decorators.js";
/**
 * A descriptor for a model property.
 */
export interface ModelPropertyDescriptor {
    /**
     * The name of the model property.
     */
    name: string;
    /**
     * The type of the model property.
     */
    type: Type;
    /**
     * Whether the model property is optional.
     */
    optional?: boolean;
    /**
     * Default value
     */
    defaultValue?: Value | undefined;
}
/**
 * Utilities for working with model properties.
 *
 * For many reflection operations, the metadata being asked for may be found
 * on the model property or the type of the model property. In such cases,
 * these operations will return the metadata from the model property if it
 * exists, or the type of the model property if it exists.
 * @typekit modelProperty
 */
export interface ModelPropertyKit {
    /**
     * Creates a modelProperty type.
     * @param desc The descriptor of the model property.
     */
    create(desc: ModelPropertyDescriptor): ModelProperty;
    /**
     * Check if the given `type` is a model property.
     *
     * @param type The type to check.
     */
    is(type: Entity): type is ModelProperty;
    /**
     * Get the encoding of the model property or its type. The property's type
     * must be a scalar.
     *
     * @param property The model property to get the encoding for.
     */
    getEncoding(property: ModelProperty): EncodeData | undefined;
    /**
     * Get the format of the model property or its type. The property's type must
     * be a string.
     *
     * @param property The model property to get the format for.
     */
    getFormat(property: ModelProperty): string | undefined;
    /**
     * Get the visibility of the model property.
     */
    getVisibilityForClass(property: ModelProperty, visibilityClass: Enum): Set<EnumMember>;
}
interface TypekitExtension {
    /**
     * Utilities for working with model properties.
     *
     * For many reflection operations, the metadata being asked for may be found
     * on the model property or the type of the model property. In such cases,
     * these operations will return the metadata from the model property if it
     * exists, or the type of the model property if it exists.
     */
    modelProperty: ModelPropertyKit;
}
declare module "../define-kit.js" {
    interface Typekit extends TypekitExtension {
    }
}
export {};
//# sourceMappingURL=model-property.d.ts.map