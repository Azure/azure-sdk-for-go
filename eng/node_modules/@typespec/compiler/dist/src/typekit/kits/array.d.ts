import { ArrayModelType, Entity, Model, Type } from "../../core/types.js";
/**
 * Typekits for working with array types(Model with number indexer).
 * @typekit array
 */
export interface ArrayKit {
    /**
     * Check if a type is an array.
     */
    is(type: Entity): type is ArrayModelType;
    /**
     * Get the element type of an array.
     */
    getElementType(type: Model): Type;
    /**
     * Create an array type.
     */
    create(elementType: Type): ArrayModelType;
}
interface TypekitExtension {
    array: ArrayKit;
}
declare module "../define-kit.js" {
    interface Typekit extends TypekitExtension {
    }
}
export {};
//# sourceMappingURL=array.d.ts.map