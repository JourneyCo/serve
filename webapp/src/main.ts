import {bootstrapApplication} from "@angular/platform-browser";
import {AppComponent} from "./app/app.component";
import {provideRouter} from "@angular/router";
import {authHttpInterceptorFn, provideAuth0} from "@auth0/auth0-angular";
import {environment as env} from "./environments/environment.development";
import {provideHttpClient, withInterceptors} from "@angular/common/http";
import {ROUTES} from "./app/routes";

bootstrapApplication(AppComponent, {
  providers: [
    provideRouter(ROUTES),
    provideAuth0({...env.auth0}),
    provideHttpClient(withInterceptors([authHttpInterceptorFn]))
  ]
})
