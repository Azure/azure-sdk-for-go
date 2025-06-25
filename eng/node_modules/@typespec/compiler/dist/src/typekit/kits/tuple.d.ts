import { Entity, Tuple, Type } from "../../core/types.js";
/**
 * @typekit tuple
 */
export interface TupleKit {
    /**
     * Check if a type is a tuple.
     */
    is(type: Entity): type is Tuple;
    /**
     * Creates a tuple type.
     *
     * @param values The tuple values, if any.
     */
    create(values?: Type[]): Tuple;
}
interface TypekitExtension {
    tuple: TupleKit;
}
declare module "../define-kit.js" {
    interface Typekit extends TypekitExtension {
    }
}
export {};
//# sourceMappingURL=tuple.d.ts.map