const path    = require('path');
const webpack = require('webpack');

const VueLoaderPlugin      = require('vue-loader/lib/plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');

module.exports = (env, argv) => {

    return {
        entry:     {
            app: './src/app.js',
        },
        output:    {
            filename:   'app.js',
            path:       path.resolve(__dirname, 'public'),
            publicPath: '/'
        },
        module:    {
            rules: [
                {
                    test:    /\.m?js$/,
                    exclude: /(node_modules|bower_components)/,
                    use:     {
                        loader:  'babel-loader',
                        options: {
                            presets: ['@babel/preset-env']
                        }
                    }
                },
                {
                    test:   /\.vue$/,
                    loader: 'vue-loader',
                },
                {
                    test: /\.scss$/,
                    use:  [
                        argv.mode === 'development' ? 'style-loader' : MiniCssExtractPlugin.loader,
                        'css-loader',
                        'sass-loader'
                    ]
                },
                {
                    test:    /\.(svg|jpg|jpeg|png|gif|ico)$/,
                    exclude: /fonts/,
                    use:     {
                        loader:  'file-loader',
                        options: {
                            name:       '[name].[ext]',
                            outputPath: 'i/',
                            publicPath: '/i/'
                        }
                    }
                }
            ]
        },
        plugins:   [
            new VueLoaderPlugin(),
            new MiniCssExtractPlugin({
                filename:      'css/style.css',
                chunkFilename: '[id].css'
            })
        ],
        resolve:   {
            extensions: ['.js', '.vue'],
            alias:      {
                'vue$': 'vue/dist/vue.runtime.min.js'
            }
        },
        devServer: {
            contentBase:      path.join(__dirname, 'public'),
            watchContentBase: true,
            compress:         true,
            port:             3210
        }
    };

};

