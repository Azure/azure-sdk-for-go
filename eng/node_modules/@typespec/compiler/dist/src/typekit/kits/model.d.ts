import { DiscriminatedUnionLegacy } from "../../core/helpers/discriminator-utils.js";
import type { Entity, Model, ModelIndexer, ModelProperty, RekeyableMap, SourceModel } from "../../core/types.js";
import { Diagnosable } from "../create-diagnosable.js";
import { DecoratorArgs } from "../utils.js";
/**
 * A descriptor for creating a model.
 */
export interface ModelDescriptor {
    /**
     * The name of the Model. If name is provided, it is a Model  declaration.
     * Otherwise, it is a Model expression.
     */
    name?: string;
    /**
     * Decorators to apply to the Model.
     */
    decorators?: DecoratorArgs[];
    /**
     * Properties of the model.
     */
    properties: Record<string, ModelProperty>;
    /**
     * Models that extend this model.
     */
    derivedModels?: Model[];
    /**
     * Models that this model extends.
     */
    sourceModels?: SourceModel[];
    /**
     * The indexer property of the model.
     */
    indexer?: ModelIndexer;
}
/**
 * Utilities for working with models.
 * @typekit model
 */
export interface ModelKit {
    /**
     * Create a model type.
     *
     * @param desc The descriptor of the model.
     */
    create(desc: ModelDescriptor): Model;
    /**
     * Check if the given `type` is a model..
     *
     * @param type The type to check.
     */
    is(type: Entity): type is Model;
    /**
     * Check this is an anonyous model. Specifically, this checks if the
     * model has a name.
     *
     * @param type The model to check.
     */
    isExpresion(type: Model): boolean;
    /**
     * If the input is anonymous (or the provided filter removes properties)
     * and there exists a named model with the same set of properties
     * (ignoring filtered properties), then return that named model.
     * Otherwise, return the input unchanged.
     *
     * This can be used by emitters to find a better name for a set of
     * properties after filtering. For example, given `{ @metadata prop:
     * string} & SomeName`, and an emitter that wishes to discard properties
     * marked with `@metadata`, the emitter can use this to recover that the
     * best name for the remaining properties is `SomeName`.
     *
     * @param model The input model
     * @param filter An optional filter to apply to the input model's
     * properties.
     */
    getEffectiveModel(model: Model, filter?: (property: ModelProperty) => boolean): Model;
    /**
     * Given a model, return the index type if one exists.
     * For example, given the model: `model Foo { ...Record<string>; ...Record<int8>; }`,
     * the index type is `Record<string | int8>`.
     * @returns the index type of the model, or undefined if no index type exists.
     */
    getIndexType: (model: Model) => Model | undefined;
    /**
     * Gets all properties from a model, explicitly defined and implicitly defined.
     * @param model model to get the properties from
     */
    getProperties(model: Model, options?: {
        includeExtended?: boolean;
    }): RekeyableMap<string, ModelProperty>;
    /**
     * Get the record representing additional properties, if there are additional properties.
     * This method checks for additional properties in the following cases:
     * 1. If the model is a Record type.
     * 2. If the model extends a Record type.
     * 3. If the model spreads a Record type.
     *
     * @param model The model to get the additional properties type of.
     * @returns The record representing additional properties, or undefined if there are none.
     */
    getAdditionalPropertiesRecord(model: Model): Model | undefined;
    /**
     * Resolves a discriminated union for the given model from inheritance.
     * @param type Model to resolve the discriminated union for.
     */
    getDiscriminatedUnion: Diagnosable<(model: Model) => DiscriminatedUnionLegacy | undefined>;
}
interface TypekitExtension {
    /**
     * Utilities for working with models.
     */
    model: ModelKit;
}
declare module "../define-kit.js" {
    interface Typekit extends TypekitExtension {
    }
}
export {};
//# sourceMappingURL=model.d.ts.map