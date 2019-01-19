const webpack = require('webpack');
const path = require('path');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const CompressionWebpackPlugin = require('compression-webpack-plugin');
const DelWebpackPlugin = require('del-webpack-plugin');
const { BundleAnalyzerPlugin } = require('webpack-bundle-analyzer');

const ROOT_DIR = path.resolve(__dirname, '../');
const DIST_DIR = path.resolve(ROOT_DIR, 'build');

const prodConfig = {
    mode: 'production',
    target: 'web',
    entry: './src/client/index.js',
    output: {
        path: DIST_DIR,
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
                            path.resolve(ROOT_DIR, 'gui/src/styles/Resources.scss')
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
        new CompressionWebpackPlugin({
            filename: '[path].gz[query]',
            algorithm: 'gzip',
            test: new RegExp('\\.(js|css)$'),
            cache: true
        }),
        new DelWebpackPlugin({
            exclude: ['bundle.js.gz', 'index.html', 'main.css.gz'],
            keepGeneratedAssets: false,
            info: true,
        })
    ]
};

if (process.env.NODE_ANALYZE) {
    prodConfig.plugins.push(new BundleAnalyzerPlugin());
}

module.exports = prodConfig;