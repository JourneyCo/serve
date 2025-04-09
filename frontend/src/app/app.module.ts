import { ApplicationConfig, importProvidersFrom, InjectionToken } from '@angular/core';
import { provideRouter } from '@angular/router';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { provideAnimations } from '@angular/platform-browser/animations';
import { provideAuth0 } from '@auth0/auth0-angular';
import { GoogleMapsModule } from '@angular/google-maps';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { routes } from './app-routing.module';
import { environment as env } from '../environments/environment';
export const GOOGLE_MAPS_API_KEY = new InjectionToken<string>('google-maps-api-key');
import { authInterceptor } from './services/auth.interceptor';

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    provideHttpClient(withInterceptors([authInterceptor])),
    provideAnimations(),
    provideAuth0({...env.auth0}),
    importProvidersFrom(
      GoogleMapsModule,
      NgbModule
    ),
    {
      provide: GOOGLE_MAPS_API_KEY,
      useValue: env.googleMapsApiKey
    }
  ]
};
