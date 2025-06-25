import { getEffectiveModelType } from "../../core/checker.js";
import { getDiscriminatedUnionFromInheritance, } from "../../core/helpers/discriminator-utils.js";
import { getDiscriminator } from "../../core/intrinsic-type-state.js";
import { createRekeyableMap } from "../../utils/misc.js";
import { createDiagnosable } from "../create-diagnosable.js";
import { defineKit } from "../define-kit.js";
import { copyMap, decoratorApplication } from "../utils.js";
const indexCache = new Map();
defineKit({
    model: {
        create(desc) {
            const properties = createRekeyableMap(Array.from(Object.entries(desc.properties)));
            const model = this.program.checker.createType({
                kind: "Model",
                name: desc.name ?? "",
                decorators: decoratorApplication(this, desc.decorators),
                properties: properties,
                derivedModels: desc.derivedModels ?? [],
                sourceModels: desc.sourceModels ?? [],
                indexer: desc.indexer,
            });
            this.program.checker.finishType(model);
            return model;
        },
        is(type) {
            return type.entityKind === "Type" && type.kind === "Model";
        },
        isExpresion(type) {
            return type.name === "";
        },
        getEffectiveModel(model, filter) {
            return getEffectiveModelType(this.program, model, filter);
        },
        getIndexType(model) {
            if (indexCache.has(model)) {
                return indexCache.get(model);
            }
            if (!model.indexer) {
                return undefined;
            }
            if (model.indexer.key.name === "string") {
                const record = this.record.create(model.indexer.value);
                indexCache.set(model, record);
                return record;
            }
            if (model.indexer.key.name === "integer") {
                const array = this.array.create(model.indexer.value);
                indexCache.set(model, array);
                return array;
            }
            return undefined;
        },
        getProperties(model, options = {}) {
            // Add explicitly defined properties
            const properties = copyMap(model.properties);
            // Add discriminator property if it exists
            const discriminator = this.model.getDiscriminatedUnion(model);
            if (discriminator) {
                const discriminatorName = discriminator.propertyName;
                properties.set(discriminatorName, this.modelProperty.create({ name: discriminatorName, type: this.builtin.string }));
            }
            if (options.includeExtended) {
                let base = model.baseModel;
                while (base) {
                    for (const [key, value] of base.properties) {
                        if (!properties.has(key)) {
                            properties.set(key, value);
                        }
                    }
                    base = base.baseModel;
                }
            }
            // TODO: Add Spread?
            return properties;
        },
        getAdditionalPropertiesRecord(model) {
            // model MyModel is Record<> {} should be model with additional properties
            if (this.model.is(model) && model.sourceModel && this.record.is(model.sourceModel)) {
                return model.sourceModel;
            }
            // model MyModel extends Record<> {} should be model with additional properties
            if (model.baseModel && this.record.is(model.baseModel)) {
                return model.baseModel;
            }
            // model MyModel { ...Record<>} should be model with additional properties
            const indexType = this.model.getIndexType(model);
            if (indexType && this.record.is(indexType)) {
                return indexType;
            }
            if (model.baseModel) {
                return this.model.getAdditionalPropertiesRecord(model.baseModel);
            }
            return undefined;
        },
        getDiscriminatedUnion: createDiagnosable(function (model) {
            const discriminator = getDiscriminator(this.program, model);
            if (!discriminator) {
                return [undefined, []];
            }
            return getDiscriminatedUnionFromInheritance(model, discriminator);
        }),
    },
});
//# sourceMappingURL=model.js.map