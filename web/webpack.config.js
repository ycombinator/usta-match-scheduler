const HtmlWebpackPlugin = require('html-webpack-plugin');
const path = require('path');

module.exports = {
 module: {
    rules: [
        {
            test: /\.js/,
            exclude: /node_modules/,
            use: {
                loader: 'babel-loader',
            }
        },
        {
            test: /\.css$/i,
            use: ["style-loader", "css-loader"],
          },
    ]
 },
 plugins: [
    new HtmlWebpackPlugin({
        template: 'public/index.html',
        filename: 'index.html',
    }),
 ],
 devServer: {
    static: {
        directory: path.join(__dirname, 'dist'),
    },
    compress: true,
    port: 9000,
    allowedHosts: "all"
 },
 resolve: {
    extensions: [".js", ".jsx"]
 }
}