const {GenerateSW} = require('workbox-webpack-plugin');
module.exports = {

    configureWebpack: {
        plugins: [ new GenerateSW()],
    },



    pwa: {
        name: 'Eta Blog',
        themeColor: '#4DBA87',
        msTileColor: '#000000',
        appleMobileWebAppCapable: 'yes',
        appleMobileWebAppStatusBarStyle: 'black',

        // configure the workbox plugin
        workboxPluginMode: 'InjectManifest',
        workboxOptions: {
        }
    }
};
