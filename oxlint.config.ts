import { defineConfig } from "oxlint";

export default defineConfig({
  env: { browser: true },
  plugins: ["eslint", "unicorn", "oxc", "typescript", "react", "react-perf"],
  categories: { correctness: "error" },
  rules: {
    "typescript/no-unused-vars": ["error", { fix: { imports: "safe-fix" } }],
    "react/rules-of-hooks": "error",
    "react/exhaustive-deps": "error",
    "react/only-export-components": [
      "error",
      { allowConstantExport: true, allowExportNames: ["Route"] },
    ],
  },
  overrides: [
    {
      files: ["client/routes/**"],
      rules: { "react/only-export-components": "off" },
    },
  ],
});
