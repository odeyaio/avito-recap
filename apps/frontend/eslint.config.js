import eslint from "@eslint/js";
import tanstackQuery from "@tanstack/eslint-plugin-query";
import globals from "globals";
import jsxA11y from "eslint-plugin-jsx-a11y";
import reactHooks from "eslint-plugin-react-hooks";
import reactRefresh from "eslint-plugin-react-refresh";
import tseslint from "typescript-eslint";

export default tseslint.config(
  {
    ignores: ["dist", "src/api/generated", "src/api/mocks/generated"],
  },
  eslint.configs.recommended,
  ...tseslint.configs.recommended,
  jsxA11y.flatConfigs.recommended,
  ...tanstackQuery.configs["flat/recommended"],
  {
    files: ["**/*.{ts,tsx}"],
    languageOptions: {
      ecmaVersion: 2022,
      globals: globals.browser,
    },
    plugins: {
      "react-hooks": reactHooks,
      "react-refresh": reactRefresh,
    },
    rules: {
      ...reactHooks.configs.flat.recommended.rules,
      "react-refresh/only-export-components": [
        "warn",
        { allowConstantExport: true },
      ],
    },
  },
  {
    files: ["src/pages/**/*.tsx", "src/routes.tsx"],
    rules: {
      // Route modules follow React Router's data-router convention of
      // co-locating a loader/action next to the route's Component.
      "react-refresh/only-export-components": "off",
    },
  },
);
