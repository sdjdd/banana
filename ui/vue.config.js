module.exports = {
    devServer: {
        proxy: 'http://localhost:80'
    },
    css: {
        loaderOptions: {
            less: {
                javascriptEnabled: true
            }
        }
    }
}