import { ApplicationConfig, importProvidersFrom } from '@angular/core';
import { provideRouter } from '@angular/router';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { provideAnimations } from '@angular/platform-browser/animations';
import {authHttpInterceptorFn, provideAuth0} from '@auth0/auth0-angular';
import { GoogleMapsModule } from '@angular/google-maps';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { routes } from './app-routing.module';
import { environment as env } from '../environments/environment';
import {NgxLinkifyjsModule} from 'ngx-linkifyjs-v2';
import {provideEnvironmentNgxMask} from 'ngx-mask';
import {RECAPTCHA_V3_SITE_KEY, RecaptchaV3Module} from 'ng-recaptcha-2';
import {CookieService} from 'ngx-cookie-service';

export const appConfig: ApplicationConfig = {
  providers: [
    { provide: RECAPTCHA_V3_SITE_KEY, useValue: 'fakevalue' },
    provideRouter(routes),
    provideAnimations(),
    provideAuth0({...env.auth0}),
    provideEnvironmentNgxMask(),
    importProvidersFrom(
      GoogleMapsModule,
      NgbModule,
        NgxLinkifyjsModule.forRoot(),
        RecaptchaV3Module
    ),
    provideHttpClient(withInterceptors([authHttpInterceptorFn])),
      CookieService
  ]
};
