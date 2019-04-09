const path = require('path');
const WebpackBar = require('webpackbar');
const { BundleAnalyzerPlugin } = require('webpack-bundle-analyzer');

const prodConfig = {
    mode: 'production',
    target: 'electron-renderer',
    entry: path.resolve(__dirname, 'src/client/index.jsx'),
    output: {
        path: path.resolve(__dirname, 'resources'),
        filename: 'bundle.js'
    },
    performance: {
        hints: false
    },
    resolve: {
        extensions: ['.js', '.jsx', '.json', '.scss'],
        modules: [
            path.resolve(__dirname, 'src'),
            path.resolve(__dirname, 'node_modules')
        ]
    },
    module: {
        rules: [{
            test: /\.(sa|sc|c)ss$/,
            use: [
                'style-loader',
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
            loader: 'babel-loader',
            exclude: /node_modules/,
            options: {
                presets: ['@babel/preset-react', '@babel/preset-env'],
                plugins: [
                    ["@babel/transform-runtime"]
                ]
            }
        },{
            test: /\.(woff|woff2|eot|ttf|otf)$/,
            exclude: /node_modules/,
            use: ['url-loader']
        }]
    },
    plugins: [
        new WebpackBar()
    ]
};

if (process.env.NODE_ANALYZE) {
    prodConfig.plugins.push(new BundleAnalyzerPlugin());
}

module.exports = prodConfig;