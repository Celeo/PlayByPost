const path = require('path')
const WebpackShellPlugin = require('webpack-shell-plugin')

module.exports = {
  entry: path.resolve(__dirname, 'editor.js'),
  output: {
    path: __dirname,
    filename: 'editor.dist.js'
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
        'mv editor.dist.js ../pbp/static'
      ]
    })
  ]
}
