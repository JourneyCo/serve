import { ApplicationConfig, importProvidersFrom, InjectionToken } from '@angular/core';
import { provideRouter } from '@angular/router';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { provideAnimations } from '@angular/platform-browser/animations';

// Auth0
import { AuthModule, provideAuth0 } from '@auth0/auth0-angular';

// Google Maps
import { GoogleMapsModule } from '@angular/google-maps';

// App routing
import { routes } from './app-routing.module';

// Environment
import { environment as env } from '../environments/environment';

// Define injection token for Google Maps API Key
export const GOOGLE_MAPS_API_KEY = new InjectionToken<string>('google-maps-api-key');

// Custom Interceptor
import { authInterceptor } from './services/auth.interceptor';

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    provideHttpClient(withInterceptors([authInterceptor])),
    provideAnimations(),
    provideAuth0({...env.auth0}),
    importProvidersFrom(
      GoogleMapsModule
    ),
    {
      provide: GOOGLE_MAPS_API_KEY,
      useValue: env.googleMapsApiKey
    }
  ]
};
