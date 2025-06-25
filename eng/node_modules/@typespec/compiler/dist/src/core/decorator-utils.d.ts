import type { Program } from "./program.js";
import { DecoratorContext, DecoratorFunction, Diagnostic, DiagnosticTarget, IntrinsicScalarName, ModelProperty, Scalar, Type } from "./types.js";
export type TypeSpecValue = Type | string | number | boolean;
export type InferredTypeSpecValue<K extends TypeKind> = K extends "Any" ? TypeSpecValue : K extends (infer T extends Type["kind"])[] ? InferredTypeSpecValue<T> : K extends "String" ? string : K extends "Number" ? number : K extends "Boolean" ? boolean : Type & {
    kind: K;
};
/**
 * Validate the decorator target is matching the expected value.
 * @param context
 * @param target
 * @param decoratorName
 * @param expectedType
 * @returns
 */
export declare function validateDecoratorTarget<K extends TypeKind>(context: DecoratorContext, target: Type, decoratorName: string, expectedType: K | readonly K[]): target is K extends "Any" ? Type : Type & {
    kind: K;
};
export declare function isIntrinsicType(program: Program, type: Scalar, kind: IntrinsicScalarName): boolean;
/**
 * Check if the given target is of any of the TypeSpec types.
 * @param target Target to validate.
 * @param expectedType One or multiple allowed TypeSpec types.
 * @returns boolean if the target is of one of the allowed types.
 */
export declare function isTypeSpecValueTypeOf<K extends TypeKind>(target: TypeSpecValue, expectedType: K | readonly K[]): target is InferredTypeSpecValue<K>;
export interface DecoratorDefinition<T extends TypeKind, P extends readonly DecoratorParamDefinition<TypeKind>[], S extends DecoratorParamDefinition<TypeKind> | undefined = undefined> {
    /**
     * Name of the decorator.
     */
    readonly name: string;
    /**
     * Decorator target.
     */
    readonly target: T | readonly T[];
    /**
     * List of positional arguments in the function.
     */
    readonly args: P;
    /**
     * Optional Type of the spread args at the end of the function if applicable.
     */
    readonly spreadArgs?: S;
}
export interface DecoratorParamDefinition<K extends TypeKind> {
    /**
     * Kind of the parameter
     */
    readonly kind: K | readonly K[];
    /**
     * Is the parameter optional.
     */
    readonly optional?: boolean;
}
type InferParameters<P extends readonly DecoratorParamDefinition<TypeKind>[], S extends DecoratorParamDefinition<TypeKind> | undefined> = S extends undefined ? InferPosParameters<P> : [...InferPosParameters<P>, ...InferSpreadParameter<S>];
type InferSpreadParameter<S extends DecoratorParamDefinition<TypeKind> | undefined> = S extends DecoratorParamDefinition<Type["kind"]> ? InferParameter<S>[] : never;
type InferPosParameters<P extends readonly DecoratorParamDefinition<TypeKind>[]> = {
    [K in keyof P]: InferParameter<P[K]>;
};
type InferParameter<P extends DecoratorParamDefinition<TypeKind>> = P["optional"] extends true ? InferParameterKind<P["kind"]> | undefined : InferParameterKind<P["kind"]>;
type InferParameterKind<P extends TypeKind | readonly TypeKind[]> = P extends readonly (infer T extends TypeKind)[] ? InferredTypeSpecValue<T> : P extends TypeKind ? InferredTypeSpecValue<P> : never;
export interface DecoratorValidator<T extends TypeKind, P extends readonly DecoratorParamDefinition<TypeKind>[], S extends DecoratorParamDefinition<TypeKind> | undefined = undefined> {
    validate(context: DecoratorContext, target: InferredTypeSpecValue<T>, parameters: InferParameters<P, S>): boolean;
}
export type TypeKind = Type["kind"] | "Any";
export declare function validateDecoratorParamCount(context: DecoratorContext, min: number, max: number | undefined, parameters: unknown[]): boolean;
/**
 * Convert a TypeSpec type to a serializable Json object.
 * Emits diagnostics if the given type is invalid
 * @param typespecType The type to convert to Json data
 * @param target The diagnostic target in case of errors.
 */
export declare function typespecTypeToJson<T>(typespecType: TypeSpecValue, target: DiagnosticTarget): [T | undefined, Diagnostic[]];
export declare function validateDecoratorUniqueOnNode(context: DecoratorContext, type: Type, decorator: DecoratorFunction): boolean;
/**
 * Validate that a given decorator is not on a type or any of its base types.
 * Useful to check for decorator usage that conflicts with another decorator.
 * @param context Decorator context
 * @param type The type to check
 * @param badDecorator The decorator we don't want present
 * @param givenDecorator The decorator that is the reason why we don't want the bad decorator present
 * @returns Whether the decorator application is valid
 */
export declare function validateDecoratorNotOnType(context: DecoratorContext, type: Type, badDecorator: DecoratorFunction, givenDecorator: DecoratorFunction): boolean;
/**
 * Return the type of the property or the model itself.
 */
export declare function getPropertyType(target: Scalar | ModelProperty): Type;
export {};
//# sourceMappingURL=decorator-utils.d.ts.map