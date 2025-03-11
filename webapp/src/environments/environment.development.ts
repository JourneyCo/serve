export const environment = {
    production: false,
    apiUrl: 'http://localhost:8080/api/v1',
    auth0: {
      domain: 'dev-dnuncdnpl8446bmt.us.auth0.com',
      clientId: 'BaYqp3c6XO3GQTqIIStocfRmVxxFRhBc',
      authorizationParams: {
        redirect_uri: 'http://localhost:3000/projects',
      },
      errorPath: '/projects',
    },
};
