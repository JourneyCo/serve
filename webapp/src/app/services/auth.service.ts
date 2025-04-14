import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { AuthModule, AuthService as Auth0Service } from '@auth0/auth0-angular';
import { Observable, from, of, throwError } from 'rxjs';
import { tap, catchError, switchMap, shareReplay, map, timeout } from 'rxjs/operators';
import { environment } from '../../environments/environment';

// interface AuthConfig {
//   domain: string;
//   clientId: string;
//   audience: string;
//   redirectUri: string;
// }

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private authConfigUrl = `${environment.apiUrl}/auth/config`;
  private isInitialized = false;

  constructor(
    private http: HttpClient,
    private auth0Service: Auth0Service,
  ) {}

  // initAuth(): Observable<any> {
  //   if (this.isInitialized) {
  //     return of(true);
  //   }
  //
  //   // Add timeout of 10 seconds to avoid hanging indefinitely
  //   return this.http.get<AuthConfig>(this.authConfigUrl, {
  //     headers: { 'Cache-Control': 'no-cache' }
  //   }).pipe(
  //     timeout(10000), // 10 seconds timeout
  //     tap((config: AuthConfig) => {
  //       console.log('Auth0 config loaded successfully');
  //       // Dynamically configure Auth0
  //       // This is a hack since Angular doesn't allow dynamic module configuration after bootstrap
  //       // In a real app, we'd use Angular's APP_INITIALIZER for this
  //       const auth0Module = (this.auth0Service as any)['auth0Client'];
  //       if (auth0Module) {
  //         auth0Module['options'] = {
  //           domain: config.domain,
  //           clientId: config.clientId,
  //           authorizationParams: {
  //             redirect_uri: config.redirectUri,
  //             audience: config.audience
  //           }
  //         };
  //       }
  //
  //       this.isInitialized = true;
  //     }),
  //     catchError(error => {
  //       console.error('Failed to load Auth0 configuration', error);
  //       // Allow the app to continue even without Auth0 in development
  //       this.isInitialized = true;
  //       return throwError(() => 'Failed to initialize authentication. Please try again later.');
  //     })
  //   );
  // }

  login(): void {
    this.auth0Service.loginWithRedirect();
  }

  logout(): void {
    this.auth0Service.logout({
      logoutParams: {
        returnTo: window.location.origin
      }
    });
  }

  isAuthenticated(): Observable<boolean> {
    return this.auth0Service.isAuthenticated$;
  }

  // getToken(): Observable<string> {
  //   // Make sure auth is initialized first
  //   if (!this.isInitialized) {
  //     return this.initAuth().pipe(
  //       switchMap(() => this.auth0Service.getAccessTokenSilently()),
  //       catchError(error => {
  //         console.error('Error getting auth token:', error);
  //         return of(''); // Return empty token if there's an error
  //       })
  //     );
  //   }
  //   return this.auth0Service.getAccessTokenSilently().pipe(
  //     catchError(error => {
  //       console.error('Error getting auth token:', error);
  //       return of(''); // Return empty token if there's an error
  //     })
  //   );
  // }
}
