import { defineConfig } from "oxfmt";

export default defineConfig({
  printWidth: 80,
  embeddedLanguageFormatting: "auto",
  sortPackageJson: true,
  sortImports: {
    groups: [
      "type-import",
      ["value-builtin", "value-external"],
      "type-internal",
      "value-internal",
      ["type-parent", "type-sibling", "type-index"],
      ["value-parent", "value-sibling", "value-index"],
      "unknown",
    ],
  },
  ignorePatterns: ["pnpm-lock.yaml", "charts/**/*.yaml", "routeTree.gen.ts"],
});
