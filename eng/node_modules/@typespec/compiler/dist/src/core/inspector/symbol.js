import pc from "picocolors";
import { SyntaxKind } from "../types.js";
/**
 * @internal
 */
export function inspectSymbol(sym, links = {}) {
    let output = `
${pc.blue(pc.inverse(` sym `))} ${pc.white(sym.name)}
${pc.dim("flags")} ${inspectSymbolFlags(sym.flags)}
  `.trim();
    if (sym.declarations && sym.declarations.length > 0) {
        const decls = sym.declarations.map((d) => SyntaxKind[d.kind]).join("\n");
        output += `\n${pc.dim("declarations")} ${decls}`;
    }
    if (sym.exports) {
        output += `\n${pc.dim("exports")} ${[...sym.exports.keys()].join(", ")}`;
    }
    if (sym.id) {
        output += `\n${pc.dim("id")} ${sym.id}`;
    }
    if (sym.members) {
        output += `\n${pc.dim("members")} ${[...sym.members.keys()].join(", ")}`;
    }
    if (sym.metatypeMembers) {
        output += `\n${pc.dim("metatypeMembers")} ${[...sym.metatypeMembers.keys()].join(", ")}`;
    }
    if (sym.parent) {
        output += `\n${pc.dim("parent")} ${sym.parent.name}`;
    }
    if (sym.symbolSource) {
        output += `\n${pc.dim("symbolSource")} ${sym.symbolSource.name}`;
    }
    if (sym.type) {
        output += `\n${pc.dim("type")} ${"name" in sym.type && sym.type.name ? String(sym.type.name) : sym.type.kind}`;
    }
    if (sym.value) {
        output += `\n${pc.dim("value")} present`;
    }
    if (Object.keys(links).length > 0) {
        output += `\nlinks\n`;
        if (links.declaredType) {
            output += `\n${pc.dim("declaredType")} ${"name" in links.declaredType && links.declaredType.name
                ? String(links.declaredType.name)
                : links.declaredType.kind}`;
        }
        if (links.instantiations) {
            output += `\n${pc.dim("instantiations")} initialized`;
        }
        if (links.type) {
            output += `\n${pc.dim("type")} ${"name" in links.type && links.type.name ? String(links.type.name) : links.type.kind}`;
        }
    }
    return output;
}
const flagsNames = [
    [2 /* SymbolFlags.Model */, "Model"],
    [4 /* SymbolFlags.Scalar */, "Scalar"],
    [8 /* SymbolFlags.Operation */, "Operation"],
    [16 /* SymbolFlags.Enum */, "Enum"],
    [32 /* SymbolFlags.Interface */, "Interface"],
    [64 /* SymbolFlags.Union */, "Union"],
    [128 /* SymbolFlags.Alias */, "Alias"],
    [256 /* SymbolFlags.Namespace */, "Namespace"],
    [512 /* SymbolFlags.Decorator */, "Decorator"],
    [1024 /* SymbolFlags.TemplateParameter */, "TemplateParameter"],
    [2048 /* SymbolFlags.Function */, "Function"],
    [4096 /* SymbolFlags.FunctionParameter */, "FunctionParameter"],
    [8192 /* SymbolFlags.Using */, "Using"],
    [16384 /* SymbolFlags.DuplicateUsing */, "DuplicateUsing"],
    [32768 /* SymbolFlags.SourceFile */, "SourceFile"],
    [65536 /* SymbolFlags.Member */, "Member"],
    [131072 /* SymbolFlags.Const */, "Const"],
];
export function inspectSymbolFlags(flags) {
    const names = [];
    for (const [flag, name] of flagsNames) {
        if (flags & flag)
            names.push(name);
    }
    return names.join(", ");
}
//# sourceMappingURL=symbol.js.map