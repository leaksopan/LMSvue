const { defineConfig } = require("@vue/cli-service");
const webpack = require("webpack");

module.exports = defineConfig({
  transpileDependencies: true,
  configureWebpack: {
    plugins: [
      new webpack.DefinePlugin({
        // Vue feature flags
        __VUE_OPTIONS_API__: true,
        __VUE_PROD_DEVTOOLS__: false,
        __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: false,
      }),
    ],
  },
  devServer: {
    port: 8080,
    proxy: {
      "/api": {
        target: "http://localhost:3001",
        ws: true,
        changeOrigin: true,
        logLevel: "debug",
        secure: false,
        // Don't rewrite the path
        pathRewrite: { "^/api": "/api" },
      },
    },
  },
});
