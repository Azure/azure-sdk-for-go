import { Entity, ModelProperty, Operation, Type } from "../../core/types.js";
import { PagingOperation } from "../../lib/paging.js";
import { Diagnosable } from "../create-diagnosable.js";
/**
 * A descriptor for an operation.
 */
export interface OperationDescriptor {
    /**
     * The name of the model property.
     */
    name: string;
    /**
     * The parameters to the model
     */
    parameters: ModelProperty[];
    /**
     * The return type of the model
     */
    returnType: Type;
}
/**
 * Utilities for working with operation properties.
 * @typekit operation
 */
export interface OperationKit {
    /**
     * Create an operation type.
     *
     * @param desc The descriptor of the operation.
     */
    create(desc: OperationDescriptor): Operation;
    /**
     * Check if the type is an operation.
     * @param type type to check
     */
    is(type: Entity): type is Operation;
    /**
     * Get the paging operation's metadata for an operation.
     * @param operation operation to get the paging operation for
     */
    getPagingMetadata: Diagnosable<(operation: Operation) => PagingOperation | undefined>;
}
interface TypekitExtension {
    /**
     * Utilities for working with operation properties.
     */
    operation: OperationKit;
}
declare module "../define-kit.js" {
    interface Typekit extends TypekitExtension {
    }
}
export {};
//# sourceMappingURL=operation.d.ts.map