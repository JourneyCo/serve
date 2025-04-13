export const environment = {
  production: false,
  apiUrl: 'http://localhost:8080', // Empty to use relative URLs (proxy will handle it)
  googleMapsApiKey: 'AIzaSyDwo9SyW5J7vMQX8DAz6faesTedB2s0Csw', // Will be replaced with the actual API key at runtime
  auth0: {
    domain: 'dev-dnuncdnpl8446bmt.us.auth0.com',
    clientId: 'BaYqp3c6XO3GQTqIIStocfRmVxxFRhBc',
    authorizationParams: {
      redirect_uri: 'http://localhost:3000',
      audience: 'https://serve.journeyco.com',
    },
    errorPath: '/projects',
    // The AuthHttpInterceptor configuration
    httpInterceptor: {
      allowedList: [
        // Attach access tokens to any calls to '/api' (exact match)
        // '/api',

        // Attach access tokens to any calls that start with '/api/'
        '/api/v1/*', '*', 'http://localhost:8080/api/v1/*',

        // Match anything starting with /api/accounts, but also specify the audience and scope the attached
        // access token must have
        // {
        //   uri: '/api/accounts/*',
        //   tokenOptions: {
        //     authorizationParams: {
        //       audience: 'http://my-api/',
        //       scope: 'read:accounts',
        //     }
        //   },
        // },

        // Matching on HTTP method
        // {
        //   uri: '/api/orders',
        //   httpMethod: 'post',
        //   tokenOptions: {
        //     authorizationParams: {
        //       audience: 'http://my-api/',
        //       scope: 'write:orders',
        //     }
        //   },
        // },

        // Using an absolute URI
        // {
        //   uri: 'https://your-domain.auth0.com/api/v2/users',
        //   tokenOptions: {
        //     authorizationParams: {
        //       audience: 'https://your-domain.com/api/v2/',
        //       scope: 'read:users',
        //     }
        //   },
        // },
      ],
    },
  },
};