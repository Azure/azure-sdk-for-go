import type { Entity, Scalar } from "../../core/types.js";
import { type EncodeData } from "../../lib/decorators.js";
/**
 * Operations for scalar types like strings, numerics, booleans, dates, etc.
 * @typekit scalar
 */
export interface ScalarKit {
    /**
     * Check if `type` is any scalar type.
     *
     * @param type The type to check.
     */
    is(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard boolean type.
     *
     * @param type The type to check.
     */
    isBoolean(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard bytes type.
     *
     * @param type The type to check.
     */
    isBytes(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard decimal type.
     *
     * @param type The type to check.
     */
    isDecimal(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard decimal128 type.
     *
     * @param type The type to check.
     */
    isDecimal128(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard duration type.
     *
     * @param type The type to check.
     */
    isDuration(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard float type.
     *
     * @param type The type to check.
     */
    isFloat(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard float32 type.
     *
     * @param type The type to check.
     */
    isFloat32(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard float64 type.
     *
     * @param type The type to check.
     */
    isFloat64(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard int8 type.
     *
     * @param type The type to check.
     */
    isInt8(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard int16 type.
     *
     * @param type The type to check.
     */
    isInt16(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard int32 type.
     *
     * @param type The type to check.
     */
    isInt32(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard int64 type.
     *
     * @param type The type to check.
     */
    isInt64(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard integer type.
     *
     * @param type The type to check.
     */
    isInteger(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard offsetDateTime type.
     *
     * @param type The type to check.
     */
    isOffsetDateTime(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard plainDate type.
     *
     * @param type The type to check.
     */
    isPlainDate(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard plainTime type.
     *
     * @param type The type to check.
     */
    isPlainTime(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard safeint type.
     *
     * @param type The type to check.
     */
    isSafeint(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard uint8 type.
     *
     * @param type The type to check.
     */
    isUint8(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard uint16 type.
     *
     * @param type The type to check.
     */
    isUint16(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard uint32 type.
     *
     * @param type The type to check.
     */
    isUint32(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard uint64 type.
     *
     * @param type The type to check.
     */
    isUint64(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard url type.
     *
     * @param type The type to check.
     */
    isUrl(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard utcDateTime type.
     *
     * @param type The type to check.
     */
    isUtcDateTime(type: Entity): type is Scalar;
    /**
     *
     * @param type The type to check.
     */
    isNumeric(type: Entity): type is Scalar;
    /**
     * Check if `type` is exactly the standard string type.
     *
     * @param type The type to check.
     */
    isString(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard boolean type.
     *
     * @param type The type to check.
     */
    extendsBoolean(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard string type.
     *
     * @param type The type to check.
     */
    extendsString(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard numeric type.
     *
     * @param type The type to check.
     */
    extendsNumeric(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard bytes type.
     *
     * @param type The type to check.
     */
    extendsBytes(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard decimal type.
     *
     * @param type The type to check.
     */
    extendsDecimal(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard decimal128 type.
     *
     * @param type The type to check.
     */
    extendsDecimal128(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard duration type.
     *
     * @param type The type to check.
     */
    extendsDuration(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard float type.
     *
     * @param type The type to check.
     */
    extendsFloat(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard float32 type.
     *
     * @param type The type to check.
     */
    extendsFloat32(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard float64 type.
     *
     * @param type The type to check.
     */
    extendsFloat64(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard int8 type.
     *
     * @param type The type to check.
     */
    extendsInt8(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard int16 type.
     *
     * @param type The type to check.
     */
    extendsInt16(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard int32 type.
     *
     * @param type The type to check.
     */
    extendsInt32(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard int64 type.
     *
     * @param type The type to check.
     */
    extendsInt64(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard integer type.
     *
     * @param type The type to check.
     */
    extendsInteger(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard offsetDateTime type.
     *
     * @param type The type to check.
     */
    extendsOffsetDateTime(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard plainDate type.
     *
     * @param type The type to check.
     */
    extendsPlainDate(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard plainTime type.
     *
     * @param type The type to check.
     */
    extendsPlainTime(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard safeint type.
     *
     * @param type The type to check.
     */
    extendsSafeint(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard uint8 type.
     *
     * @param type The type to check.
     */
    extendsUint8(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard uint16 type.
     *
     * @param type The type to check.
     */
    extendsUint16(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard uint32 type.
     *
     * @param type The type to check.
     */
    extendsUint32(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard uint64 type.
     *
     * @param type The type to check.
     */
    extendsUint64(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard url type.
     *
     * @param type The type to check.
     */
    extendsUrl(type: Entity): type is Scalar;
    /**
     * Check if `type` extends the standard utcDateTime type.
     *
     * @param type The type to check.
     */
    extendsUtcDateTime(type: Entity): type is Scalar;
    /**
     * Get the standard built-in base type of a scalar. For all built-in scalar
     * types (numeric, string, int32, etc.) this will just return the scalar
     * type. For user-defined scalars, this will return the first base scalar
     * that is built-in. For user-defined scalars without a standard base type,
     * this will return null.
     *
     * @param type The scalar to check.
     */
    getStdBase(type: Scalar): Scalar | null;
    /**
     * Get the encoding information for a scalar type. Returns undefined if no
     * encoding data is specified.
     *
     * Note: This will return the encoding data for the scalar type itself, not
     * the model property that uses the scalar type. If this scalar might be
     * referenced from a model property, use {@link modelProperty.getEncoding}
     * instead.
     *
     * @param scalar The scalar to get the encoding data for.
     */
    getEncoding(scalar: Scalar): EncodeData | undefined;
    /**
     * Get the well-known format for a string scalar. Returns undefined if no
     * format is specified.
     *
     * Note: This will return the format data for the scalar type itself, not
     * the model property that uses the scalar type. If this scalar might be
     * referenced from a model property, use {@link ModelPropertyKit.getEncoding}
     * instead.
     *
     * @param scalar The scalar to get the format for.
     */
    getFormat(scalar: Scalar): string | undefined;
}
interface TypekitExtension {
    /**
     * Operations for scalar types like strings, numerics, booleans, dates, etc.
     */
    scalar: ScalarKit;
}
declare module "../define-kit.js" {
    interface Typekit extends TypekitExtension {
    }
}
export {};
//# sourceMappingURL=scalar.d.ts.map