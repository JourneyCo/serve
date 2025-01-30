import {CUSTOM_ELEMENTS_SCHEMA, NgModule} from '@angular/core';
import {BrowserModule} from '@angular/platform-browser';
import {AppComponent} from './app.component';
import {provideHttpClient, withInterceptorsFromDi} from '@angular/common/http';
import {HeaderComponent} from './header/header.component';
import {NgbModule} from "@ng-bootstrap/ng-bootstrap";
import {AppRoutingModule} from './app-routing.module';
import {MapComponent} from "@components";
import {ProjectsComponent} from "./pages/projects/projects.component";
import {LoginButtonComponent} from "./components/auth0/login-button.component";
import {AuthModule} from "@auth0/auth0-angular";
import {environment as env} from "../environments/environment.development"
import {APIService} from "@services";
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';

@NgModule({
  declarations: [AppComponent],
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  bootstrap: [AppComponent],
  imports: [BrowserModule,
    NgbModule,
    AppRoutingModule,
    MapComponent,
    LoginButtonComponent,
    AuthModule.forRoot({
      ...env.auth0,
    }),
    ProjectsComponent,
    HeaderComponent,
  ],
  providers:
    [provideHttpClient(withInterceptorsFromDi()),
      APIService,
      provideAnimationsAsync(),
  ]
})
export class AppModule { }
