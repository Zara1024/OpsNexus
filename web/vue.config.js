/**
 * vue 配置中心
 */
const { defineConfig } = require('@vue/cli-service')
const sass = require('sass')
const webpack = require('webpack')

const apiProxyTarget = process.env.VUE_APP_API_PROXY_TARGET || 'http://10.0.0.200:8080'

module.exports = defineConfig({
  lintOnSave: false,
  productionSourceMap: false,
  publicPath: '/',
  configureWebpack: {
    plugins: [
      new webpack.DefinePlugin({
        __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: 'false'
      })
    ],
    performance: {
      hints: false
    }
  },
  outputDir: 'dist',
  assetsDir: 'static',
  css: {
    loaderOptions: {
      sass: {
        implementation: sass,
        api: 'modern'
      },
      scss: {
        implementation: sass,
        api: 'modern'
      }
    }
  },
  devServer: {
    port: 8080,
    host: '0.0.0.0',
    https: false,
    open: false,
    historyApiFallback: {
      disableDotRule: true
    },
    proxy: {
      '/api/v1': {
        target: apiProxyTarget,
        changeOrigin: true
      }
    },
    client: {
      overlay: false
    }
  }
})
