export const environment = {
  production: true,
  serveDay: "2025-07-12",
  apiUrl: "", // Empty to use relative URLs (proxy will handle it in production)
  auth0: {
    domain: "dev-dnuncdnpl8446bmt.us.auth0.com",
    clientId: "BaYqp3c6XO3GQTqIIStocfRmVxxFRhBc",
    authorizationParams: {
      redirect_uri: "https://serveday.journeycolorado.com",
      audience: "https://serve.journeyco.com",
    },
    errorPath: "/",
    // The AuthHttpInterceptor configuration
    httpInterceptor: {
      allowedList: [
        "https://serveday.journeycolorado.com/api/admin/*",
        "http://serveday.journeycolorado.com/api/admin/*",
        "/api/admin/*",
      ],
    },
  },
};
