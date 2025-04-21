import { ApplicationConfig, importProvidersFrom } from '@angular/core';
import { provideRouter } from '@angular/router';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { provideAnimations } from '@angular/platform-browser/animations';
import {authHttpInterceptorFn, provideAuth0} from '@auth0/auth0-angular';
import { GoogleMapsModule } from '@angular/google-maps';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { routes } from './app-routing.module';
import { environment as env } from '../environments/environment';

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    provideAnimations(),
    provideAuth0({...env.auth0}),
    importProvidersFrom(
      GoogleMapsModule,
      NgbModule
    ),
    provideHttpClient(withInterceptors([authHttpInterceptorFn])),
  ]
};
