const path = require('path')
const WebpackShellPlugin = require('webpack-shell-plugin')

module.exports = {
  entry: path.resolve(__dirname, 'index.js'),
  output: {
    path: __dirname,
    filename: 'index.dist.js'
  },
  resolve: {
    alias: {
      vue: 'vue/dist/vue.js'
    }
  },
  performance: {
    hints: false
  },
  target: 'web',
  plugins: [
    new WebpackShellPlugin({
      onBuildEnd: [
        'mv index.dist.js ../pbp/static'
      ]
    })
  ]
}
