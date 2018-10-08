const path = require('path')
const CopyWebpackPlugin = require('copy-webpack-plugin')
const WebpackDeletePlugin = require('webpack-delete-plugin')

module.exports = {
  entry: path.resolve(__dirname, 'index.js'),
  output: {
    path: __dirname,
    filename: 'index.dist.js'
  },
  performance: {
    hints: false
  },
  target: 'web',
  plugins: [
    new CopyWebpackPlugin([
      { from: 'index.dist.js', to: path.resolve(__dirname, '..', 'pbp', 'static', 'index.dist.js') }
    ]),
    new WebpackDeletePlugin(['index.dist.js'])
  ]
}
