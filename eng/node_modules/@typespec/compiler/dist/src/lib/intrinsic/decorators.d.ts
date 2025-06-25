import type { Program } from "../../core/program.js";
import type { DecoratorContext, ModelIndexer, Scalar, Type } from "../../core/types.js";
export declare const indexerDecorator: (context: DecoratorContext, target: Type, key: Scalar, value: Type) => void;
export declare function getIndexer(program: Program, target: Type): ModelIndexer | undefined;
//# sourceMappingURL=decorators.d.ts.map