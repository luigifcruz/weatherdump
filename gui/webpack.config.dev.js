const webpack = require('webpack');
const path = require('path');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const WebpackBar = require('webpackbar');

const buildDate = new Date().toISOString()

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
        './src/client/index.jsx'
    ],
    output: {
        path: path.resolve(__dirname, 'resources'),
        filename: 'bundle.js',
        publicPath: "/"
    },
    module: {
        rules: [
            {
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
                                path.resolve(__dirname, 'src/styles/palette.scss'),
                                path.resolve(__dirname, 'src/styles/mixins.scss')
                            ],
                        },
                    },
                ],
            },{
                test: /\.(js|jsx)$/,
                loader: "babel-loader",
                exclude: /(node_modules)/,
                resolve: {
                    extensions: [".js", ".jsx", ".json", ".scss"],
                    modules: [
                        path.resolve(__dirname, 'src'),
                        "node_modules"
                    ]
                },
                options: {
                    presets: ['@babel/react', '@babel/env']
                }
            },{
                test: /\.(woff|woff2|eot|ttf|otf)$/,
                use: ['url-loader']
            }
        ]
    },
    plugins: [
        new MiniCssExtractPlugin({
            filename: "[name].css",
            allChunks: false
        }),
        new WebpackBar(),
        new webpack.HotModuleReplacementPlugin(),
        new webpack.DefinePlugin({
            'process.env.NODE_ENV': JSON.stringify('development'),
            'BUILD_DATE': JSON.stringify(buildDate)
        }),
    ]
};

module.exports = devConfig;