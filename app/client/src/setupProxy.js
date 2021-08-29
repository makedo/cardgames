const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
  app.use(
    '/api',
    createProxyMiddleware({
      target: 'http://localhost:8080',
      changeOrigin: true,
      hostRewrite: true,
      pathRewrite: {
        '^/api': '',
      },
      secure: false,
    })
  );

};