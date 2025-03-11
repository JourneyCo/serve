import {CUSTOM_ELEMENTS_SCHEMA, NgModule} from '@angular/core';
import {BrowserModule} from '@angular/platform-browser';
import {HTTP_INTERCEPTORS } from '@angular/common/http';
import {HeaderComponent} from './header/header.component';
import {NgbModule} from "@ng-bootstrap/ng-bootstrap";
import {AppRoutingModule} from './app-routing.module';
import {MapComponent} from "@components";
import {ProjectsComponent} from "./pages/projects/projects.component";
import {AuthHttpInterceptor, authHttpInterceptorFn, AuthModule} from "@auth0/auth0-angular";
import {environment as env} from "../environments/environment.development"
import {APIService } from "@services";
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';

@NgModule({
  declarations: [],
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  bootstrap: [AppModule],
  imports: [BrowserModule,
    NgbModule,
    AppRoutingModule,
    MapComponent,
    AuthModule.forRoot({
      ...env.auth0,
    }),
    ProjectsComponent,
    HeaderComponent,
  ],
  providers:
    [
      APIService,
      provideAnimationsAsync(),
      { provide: HTTP_INTERCEPTORS, useClass: AuthHttpInterceptor, multi: true },
  ]
})
export class AppModule { }
