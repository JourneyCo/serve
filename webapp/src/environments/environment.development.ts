export const environment = {
    production: false,
    apiUrl: 'http://localhost:8080/api/v1',
    auth0: {
      domain: 'dev-thr5tcgxepxvehp4.us.auth0.com',
      clientId: 'QlleJ7sv8ie22FEnpR7ycyNrpLAbLSv9',
      authorizationParams: {
        redirect_uri: 'http://localhost:3000/callback',
      },
      errorPath: '/callback',
    },
};
