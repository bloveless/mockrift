const HtmlWebPackPlugin = require('html-webpack-plugin');
const ManifestPlugin = require('webpack-manifest-plugin');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');
const { DefinePlugin } = require('webpack');
const path = require('path');

module.exports = (env, argv) => ({
    devServer: {
        host: '0.0.0.0',
        contentBase: path.join(__dirname, 'dist'),
        compress: true,
        historyApiFallback: {
            index: '/static/react/index.html',
            rewrites: [
                { from: /^.*\.html$/, to: '/static/react/index.html' },
            ],
        },
    },
    output: {
        path: path.resolve(__dirname, '/mockrift/ui/static/react/'),
        publicPath: '/static/react/',
        filename: argv.mode === 'production' ? '[name].[chunkhash].bundle.js' : '[name].bundle.js',
    },
    optimization: {
        splitChunks: {
            chunks: 'all',
        },
    },
    module: {
        rules: [
            {
                test: /\.(js|jsx)$/,
                exclude: /node_modules/,
                use: {
                    loader: "babel-loader",
                },
            },
            {
                test: /\.html$/,
                use: [
                    {
                        loader: "html-loader",
                    },
                ],
            },
        ],
    },
    plugins: [
        new DefinePlugin({
            'process.env.MODE': JSON.stringify(argv.mode),
        }),
        new CleanWebpackPlugin(),
        new HtmlWebPackPlugin({
            template: "./src/index.html",
            filename: "index.html",
        }),
        new ManifestPlugin(),
    ],
});
