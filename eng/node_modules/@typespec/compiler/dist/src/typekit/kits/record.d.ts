import { Entity, Model, RecordModelType, Type } from "../../core/types.js";
/**
 * RecordKit provides utilities for working with Record Model types.
 * @typekit record
 */
export interface RecordKit {
    /**
     * Check if the given `type` is a Record.
     *
     * @param type The type to check.
     */
    is(type: Entity): type is RecordModelType;
    /**
     *  Get the element type of a Record
     * @param type a Record Model type
     */
    getElementType(type: Model): Type;
    /**
     * Create a Record Model type
     * @param elementType The type of the elements in the record
     */
    create(elementType: Type): RecordModelType;
}
interface TypekitExtension {
    record: RecordKit;
}
declare module "../define-kit.js" {
    interface Typekit extends TypekitExtension {
    }
}
export {};
//# sourceMappingURL=record.d.ts.map