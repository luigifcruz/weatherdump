const webpack = require('webpack');
const path = require('path');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const { BundleAnalyzerPlugin } = require('webpack-bundle-analyzer');

const prodConfig = {
    mode: 'production',
    target: 'web',
    entry: './src/client/index.js',
    output: {
        path: path.resolve(__dirname, 'resources'),
        filename: 'bundle.js'
    },
    performance: {
        hints: false
    },
    module: {
        rules: [{
            test: /\.(sa|sc|c)ss$/,
            use: [
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
        })
    ]
};

if (process.env.NODE_ANALYZE) {
    prodConfig.plugins.push(new BundleAnalyzerPlugin());
}

module.exports = prodConfig;