import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { AppComponent } from './app.component';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';
import { HeaderComponent } from './header/header.component';
import {NgbModule} from "@ng-bootstrap/ng-bootstrap";
import { AppRoutingModule } from './app-routing.module';
import {MapComponent, TableComponent} from "@components";
import {ProjectsComponent} from "./pages/projects/projects.component";
import {LoginButtonComponent} from "./components/auth0/login-button.component";
import {AuthModule} from "@auth0/auth0-angular";
import { environment as env} from "../environments/environment.development"

@NgModule({
  declarations: [
        AppComponent,
        HeaderComponent,
        ProjectsComponent,
    ],
  bootstrap: [AppComponent],
  imports: [BrowserModule,
        NgbModule,
        AppRoutingModule,
        MapComponent,
        TableComponent,
        LoginButtonComponent,
    AuthModule.forRoot({
      ...env.auth0,}),
  ],
  providers: [provideHttpClient(withInterceptorsFromDi())] })
export class AppModule { }
