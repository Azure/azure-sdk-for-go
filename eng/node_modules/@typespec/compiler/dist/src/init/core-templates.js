import { CompilerPackageRoot } from "../core/node-host.js";
import { resolvePath } from "../core/path-utils.js";
export const templatesDir = resolvePath(CompilerPackageRoot, "templates");
let typeSpecCoreTemplates;
export async function getTypeSpecCoreTemplates(host) {
    if (typeSpecCoreTemplates === undefined) {
        const file = await host.readFile(resolvePath(templatesDir, "scaffolding.json"));
        const content = JSON.parse(file.text);
        typeSpecCoreTemplates = {
            baseUri: templatesDir,
            templates: content,
        };
    }
    return typeSpecCoreTemplates;
}
//# sourceMappingURL=core-templates.js.map