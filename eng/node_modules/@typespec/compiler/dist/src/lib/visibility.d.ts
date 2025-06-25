import type { DefaultVisibilityDecorator, InvisibleDecorator, ParameterVisibilityDecorator, RemoveVisibilityDecorator, ReturnTypeVisibilityDecorator, VisibilityDecorator, WithDefaultKeyVisibilityDecorator, WithLifecycleUpdateDecorator, WithUpdateablePropertiesDecorator, WithVisibilityDecorator, WithVisibilityFilterDecorator } from "../../generated-defs/TypeSpec.js";
import type { Program } from "../core/program.js";
import { Operation } from "../core/types.js";
import { VisibilityFilter } from "../core/visibility/core.js";
export declare const $withDefaultKeyVisibility: WithDefaultKeyVisibilityDecorator;
export declare const $parameterVisibility: ParameterVisibilityDecorator;
/**
 * A context-specific provider for visibility information that applies when parameter or return type visibility
 * constraints are not explicitly specified. Visibility providers are provided by libraries that define implied
 * visibility semantics, such as `@typespec/http`.
 *
 * If you are not working in a protocol that has specific visibility semantics, you can use the
 * {@link EmptyVisibilityProvider} from this package as a default provider. It will consider all properties visible by
 * default unless otherwise explicitly specified.
 */
export interface VisibilityProvider {
    parameters(program: Program, operation: Operation): VisibilityFilter;
    returnType(program: Program, operation: Operation): VisibilityFilter;
}
/**
 * An empty visibility provider. This provider returns an empty filter that considers all properties visible. This filter
 * is used when no context-specific visibility provider is available.
 *
 * When working with an HTTP specification, use the `HttpVisibilityProvider` from the `@typespec/http` library instead.
 */
export declare const EmptyVisibilityProvider: VisibilityProvider;
/**
 * Get the visibility filter that should apply to the parameters of the given operation, or `undefined` if no parameter
 * visibility is set.
 *
 * If you are not working in a protocol that has specific implicit visibility semantics, you can use the
 * {@link EmptyVisibilityProvider} as a default provider. If you working in a protocol or context where parameters have
 * implicit visibility transformations (like HTTP), you should use the visibility provider from that library (for HTTP,
 * use the `HttpVisibilityProvider` from the `@typespec/http` library).
 *
 * @param program - the Program in which the operation is defined
 * @param operation - the Operation to get the parameter visibility filter for
 * @param defaultProvider - a provider for visibility filters that apply when no visibility constraints are explicitly
 *                         set. Defaults to an empty provider that returns an empty filter if not provided.
 * @returns a visibility filter for the parameters of the operation, or `undefined` if no parameter visibility is set
 */
export declare function getParameterVisibilityFilter(program: Program, operation: Operation, defaultProvider: VisibilityProvider): VisibilityFilter;
export declare const $returnTypeVisibility: ReturnTypeVisibilityDecorator;
/**
 * Get the visibility filter that should apply to the return type of the given operation, or `undefined` if no return
 * type visibility is set.
 *
 * @param program - the Program in which the operation is defined
 * @param operation - the Operation to get the return type visibility filter for
 * @param defaultProvider - a provider for visibility filters that apply when no visibility constraints are explicitly
 *                          set. Defaults to an empty provider that returns an empty filter if not provided.
 * @returns a visibility filter for the return type of the operation, or `undefined` if no return type visibility is set
 */
export declare function getReturnTypeVisibilityFilter(program: Program, operation: Operation, defaultProvider: VisibilityProvider): VisibilityFilter;
export declare const $visibility: VisibilityDecorator;
export declare const $removeVisibility: RemoveVisibilityDecorator;
export declare const $invisible: InvisibleDecorator;
export declare const $defaultVisibility: DefaultVisibilityDecorator;
export declare const $withVisibility: WithVisibilityDecorator;
/**
 * Filters a model for properties that are updateable.
 *
 * @param context - the program context
 * @param target - Model to filter for updateable properties
 */
export declare const $withUpdateableProperties: WithUpdateablePropertiesDecorator;
export declare const $withVisibilityFilter: WithVisibilityFilterDecorator;
export declare const $withLifecycleUpdate: WithLifecycleUpdateDecorator;
//# sourceMappingURL=visibility.d.ts.map