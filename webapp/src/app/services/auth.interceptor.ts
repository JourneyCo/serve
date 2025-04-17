import { inject } from "@angular/core";
import {
  HttpRequest,
  HttpHandlerFn,
  HttpErrorResponse,
  HttpInterceptorFn,
} from "@angular/common/http";
import { throwError } from "rxjs";
import { catchError, switchMap } from "rxjs/operators";
import { AuthService } from "./auth.service";
import { environment } from "../../environments/environment";

export const authInterceptor: HttpInterceptorFn = (
  request: HttpRequest<unknown>,
  next: HttpHandlerFn,
) => {
  const authService = inject(AuthService);

  // Skip auth for auth config endpoint to avoid circular dependency
  if (request.url.includes("/auth/config")) {
    return next(request);
  }

  // Only add authorization header for API calls to our server
  // If apiUrl is empty, we're using same-origin requests
  if (environment.apiUrl && !request.url.startsWith(environment.apiUrl)) {
    return next(request);
  }

  return authService.getToken().pipe(
    switchMap((token: any) => {
      // Only add authorization if we have a token
      if (!token) {
        return next(request);
      }

      // Clone the request and add the authorization header
      const authRequest = request.clone({
        setHeaders: {
          Authorization: `Bearer ${token}`,
        },
      });
      return next(authRequest);
    }),
    catchError((error: HttpErrorResponse) => {
      // Handle auth errors
      if (error.status === 401) {
        // Optionally redirect to login or refresh token
        console.log("Authentication error, redirecting to login...");
        authService.login();
      }
      return throwError(() => error);
    }),
  );
};
