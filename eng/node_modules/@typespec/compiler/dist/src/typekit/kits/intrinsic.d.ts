import type { Entity, ErrorType, IntrinsicType, NeverType, NullType, UnknownType, VoidType } from "../../core/types.js";
/**
 * @typekit intrinsic
 */
export interface IntrinsicKit {
    /** The intrinsic 'any' type. */
    get any(): UnknownType;
    /** The intrinsic 'error' type. */
    get error(): ErrorType;
    /** The intrinsic 'never' type. */
    get never(): NeverType;
    /** The intrinsic 'null' type. */
    get null(): NullType;
    /** The intrinsic 'void' type. */
    get void(): VoidType;
    /**
     * Check if `entity` is an intrinsic type.
     * @param entity The `entity` to check.
     */
    is(entity: Entity): entity is IntrinsicType;
}
interface TypekitExtension {
    intrinsic: IntrinsicKit;
}
declare module "../define-kit.js" {
    interface Typekit extends TypekitExtension {
    }
}
export {};
//# sourceMappingURL=intrinsic.d.ts.map