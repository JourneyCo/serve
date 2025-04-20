import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { tap } from 'rxjs/operators';
import { AuthService } from '@services';

export const AuthGuard: CanActivateFn = (route, state) => {
  const authService = inject(AuthService);

  return authService.isAuthenticated().pipe(
    tap(isAuthenticated => {
      if (!isAuthenticated) {
        // Store the attempted URL for redirecting after login
        localStorage.setItem('authRedirect', state.url);
        authService.login();
      }
    })
  );
};
