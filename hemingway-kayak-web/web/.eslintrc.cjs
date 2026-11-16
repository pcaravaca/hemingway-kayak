module.exports = {
  root: true,
  env: {
    browser: true,
    es2022: true,
  },
  extends: ['eslint:recommended', 'plugin:vue/vue3-essential'],
  parserOptions: {
    ecmaVersion: 'latest',
    sourceType: 'module',
  },
  ignorePatterns: ['dist/**', 'node_modules/**'],
  rules: {
    // El backend devuelve errores en `error`, no en `message`; el manejo es intencional.
    'vue/no-unused-vars': ['error', { ignorePattern: '^_' }],
  },
}
