import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { tap } from 'rxjs/operators';
import { AuthService } from '@services';

export const AdminGuard: CanActivateFn = (route, state) => {
  const authService = inject(AuthService);
  const router = inject(Router);

  return authService.isAdmin().pipe(
    tap(isAdmin => {
      if (!isAdmin) {
        router.navigate(['/']);
      }
    })
  );
};
