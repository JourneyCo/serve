import { Injectable } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import { AuthService as Auth0Service } from "@auth0/auth0-angular";
import { Observable, of } from "rxjs";
import {
  tap,
  catchError,
  switchMap,
  shareReplay,
  map
} from "rxjs/operators";
import { User } from '@models';
import { UserService } from "./user.service";
import { jwtDecode } from "jwt-decode";

@Injectable({
  providedIn: "root",
})
export class AuthService {
  private cachedUser: User | null = null;

  constructor(
    private http: HttpClient,
    private auth0Service: Auth0Service,
    private userService: UserService,
  ) {}

  login(): void {
    this.auth0Service.loginWithRedirect({
      appState: {
        target: "/",
      },
      authorizationParams: {
        prompt: "login",
      },
    });
  }

  logout(): void {
    this.auth0Service.logout({
      logoutParams: {
        returnTo: window.location.origin,
      },
    });
    this.cachedUser = null;
  }

  isAuthenticated(): Observable<boolean> {
    return this.auth0Service.isAuthenticated$;
  }

  getCurrentUser(): Observable<User | null> {
    if (this.cachedUser) {
      return of(this.cachedUser);
    }

    return this.auth0Service.user$.pipe(
      switchMap((auth0User) => {
        if (!auth0User) {
          return of(null);
        }

        // Try to get the user profile from our server
        return this.userService.getUserProfile().pipe(
          tap((user) => {
            this.cachedUser = user;
          }),
          catchError((error) => {
            console.error("Error fetching user profile:", error);
            return of(null);
          }),
        );
      }),
      shareReplay(1),
    );
  }

  isAdmin(): Observable<boolean> {
    return this.auth0Service.getAccessTokenSilently().pipe(
      map((token) => {
        const decodedToken: any = jwtDecode(token);
        const perms = decodedToken.permissions || [];
        return perms.includes("edit:projects");
      }),
    );
  }
}
