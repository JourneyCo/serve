export const environment = {
  production: false,
  serveDay: "2025-07-12",
  apiUrl: "http://localhost:8080", // Empty to use relative URLs (proxy will handle it)
  auth0: {
    domain: "dev-dnuncdnpl8446bmt.us.auth0.com",
    clientId: "BaYqp3c6XO3GQTqIIStocfRmVxxFRhBc",
    authorizationParams: {
      redirect_uri: "http://localhost:3000",
      audience: "https://serve.journeyco.com",
    },
    errorPath: "/",
    // The AuthHttpInterceptor configuration
    httpInterceptor: {
      allowedList: [
        // Attach access tokens to any calls that start with '/api/'
        "/api/v1/users/profile",
        "/api/v1/admin",
        {uriMatcher: (uri: string) => /\/api\/v1\/projects\/\d+\/register$/.test(uri)},
        {uriMatcher: (uri: string) => /\/api\/v1\/projects\/\d+\/cancel$/.test(uri)},
        {uriMatcher: (uri: string) => /\/api\/v1\/projects\/\d+\/registrations$/.test(uri)},
      ],
    },
  },
};
