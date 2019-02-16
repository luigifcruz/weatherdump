const webpack = require('webpack');
const path = require('path');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");

const devConfig = {
    mode: 'development',
    target: 'web',
    devtool: 'inline-source-map',
    performance: {
        hints: false
    },
    entry: [
        'react-hot-loader/patch',
        'webpack-hot-middleware/client?path=/__webpack_hmr&timeout=20000',
        './src/client/index.js'
    ],
    output: {
        path: path.resolve(__dirname, 'build'),
        filename: 'bundle.js',
        publicPath: "/"
    },
    module: {
        rules: [{
            test: /\.(sa|sc|c)ss$/,
            use: [
                'css-hot-loader',
                MiniCssExtractPlugin.loader,
                'css-loader',
                'sass-loader',
                {
                    loader: 'sass-resources-loader',
                    options: {
                        resources: [
                            path.resolve(__dirname, 'src/styles/Resources.scss')
                        ],
                    },
                },
            ],
        },{
            test: /\.(js|jsx)$/,
            loader: "babel-loader",
            exclude: /(node_modules)/,
            options: {
                presets: ['@babel/react', '@babel/env']
            }
        }]
    },
    plugins: [
        new MiniCssExtractPlugin({
            filename: "[name].css",
            allChunks: false
        }),
        new webpack.HotModuleReplacementPlugin(),
        new webpack.DefinePlugin({
            'process.env.NODE_ENV': JSON.stringify('development')
        }),
    ]
};

module.exports = devConfig;