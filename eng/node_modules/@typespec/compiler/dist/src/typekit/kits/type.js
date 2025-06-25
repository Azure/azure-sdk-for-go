import { getLocationContext } from "../../core/helpers/location-context.js";
import { getMaxItems, getMaxLength, getMaxValue, getMaxValueExclusive, getMinItems, getMinLength, getMinValue, getMinValueExclusive, } from "../../core/intrinsic-type-state.js";
import { isNeverType } from "../../core/type-utils.js";
import { getDoc, getSummary, isErrorModel } from "../../lib/decorators.js";
import { resolveEncodedName } from "../../lib/encoded-names.js";
import { createDiagnosable } from "../create-diagnosable.js";
import { defineKit } from "../define-kit.js";
import { copyMap } from "../utils.js";
import { getPlausibleName } from "../utils/get-plausible-name.js";
defineKit({
    type: {
        is(entity) {
            return entity.entityKind === "Type";
        },
        finishType(type) {
            this.program.checker.finishType(type);
        },
        clone(type) {
            let clone;
            switch (type.kind) {
                case "Model":
                    clone = this.program.checker.createType({
                        ...type,
                        decorators: [...type.decorators],
                        derivedModels: [...type.derivedModels],
                        sourceModels: type.sourceModels.map((x) => ({ ...x })),
                        properties: copyMap(type.properties),
                        indexer: type.indexer ? { ...type.indexer } : undefined,
                    });
                    break;
                case "Union":
                    clone = this.program.checker.createType({
                        ...type,
                        decorators: [...type.decorators],
                        variants: copyMap(type.variants),
                        get options() {
                            return Array.from(this.variants.values()).map((v) => v.type);
                        },
                    });
                    break;
                case "Interface":
                    clone = this.program.checker.createType({
                        ...type,
                        decorators: [...type.decorators],
                        operations: copyMap(type.operations),
                    });
                    break;
                case "Enum":
                    clone = this.program.checker.createType({
                        ...type,
                        decorators: [...type.decorators],
                        members: copyMap(type.members),
                    });
                    break;
                case "Namespace":
                    clone = this.program.checker.createType({
                        ...type,
                        decorators: [...type.decorators],
                        instantiationParameters: type.instantiationParameters
                            ? [...type.instantiationParameters]
                            : undefined,
                        models: copyMap(type.models),
                        decoratorDeclarations: copyMap(type.decoratorDeclarations),
                        enums: copyMap(type.enums),
                        unions: copyMap(type.unions),
                        operations: copyMap(type.operations),
                        interfaces: copyMap(type.interfaces),
                        namespaces: copyMap(type.namespaces),
                        scalars: copyMap(type.scalars),
                    });
                    break;
                case "Scalar":
                    clone = this.program.checker.createType({
                        ...type,
                        decorators: [...type.decorators],
                        derivedScalars: [...type.derivedScalars],
                        constructors: copyMap(type.constructors),
                    });
                    break;
                case "Tuple":
                    clone = this.program.checker.createType({
                        ...type,
                        values: [...type.values],
                    });
                    break;
                default:
                    clone = this.program.checker.createType({
                        ...type,
                        ...("decorators" in type ? { decorators: [...type.decorators] } : {}),
                    });
                    break;
            }
            this.realm.addType(clone);
            return clone;
        },
        isError(type) {
            return isErrorModel(this.program, type);
        },
        getEncodedName(type, encoding) {
            return resolveEncodedName(this.program, type, encoding);
        },
        getSummary(type) {
            return getSummary(this.program, type);
        },
        getDoc(type) {
            return getDoc(this.program, type);
        },
        getPlausibleName(type) {
            return getPlausibleName(type);
        },
        maxValue(type) {
            return getMaxValue(this.program, type);
        },
        minValue(type) {
            return getMinValue(this.program, type);
        },
        maxLength(type) {
            return getMaxLength(this.program, type);
        },
        minLength(type) {
            return getMinLength(this.program, type);
        },
        maxItems(type) {
            return getMaxItems(this.program, type);
        },
        maxValueExclusive(type) {
            return getMaxValueExclusive(this.program, type);
        },
        minValueExclusive(type) {
            return getMinValueExclusive(this.program, type);
        },
        minItems(type) {
            return getMinItems(this.program, type);
        },
        isNever(type) {
            return isNeverType(type);
        },
        isUserDefined(type) {
            return getLocationContext(this.program, type).type === "project";
        },
        inNamespace(type, namespace) {
            // A namespace is always in itself
            if (type === namespace) {
                return true;
            }
            // Handle types with known containers
            switch (type.kind) {
                case "ModelProperty":
                    if (type.model) {
                        return this.type.inNamespace(type.model, namespace);
                    }
                    break;
                case "EnumMember":
                    return this.type.inNamespace(type.enum, namespace);
                case "UnionVariant":
                    return this.type.inNamespace(type.union, namespace);
                case "Operation":
                    if (type.interface) {
                        return this.type.inNamespace(type.interface, namespace);
                    }
                    // Operations that belong to a namespace directly will be handled in the generic case
                    break;
            }
            // Generic case handles all other types
            if ("namespace" in type && type.namespace) {
                return this.type.inNamespace(type.namespace, namespace);
            }
            // If we got this far, the type does not belong to the namespace
            return false;
        },
        isAssignableTo: createDiagnosable(function (source, target, diagnosticTarget) {
            return this.program.checker.isTypeAssignableTo(source, target, diagnosticTarget ?? source);
        }),
        resolve: createDiagnosable(function (reference, kind) {
            const [type, diagnostics] = this.program.resolveTypeReference(reference);
            if (type && kind && type.kind !== kind) {
                throw new Error(`Type kind mismatch: expected ${kind}, got ${type.kind}`);
            }
            return [type, diagnostics];
        }),
    },
});
//# sourceMappingURL=type.js.map