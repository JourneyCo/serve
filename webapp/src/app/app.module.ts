import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { AppComponent } from './app.component';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';
import { HeaderComponent } from './header/header.component';
import {NgbModule} from "@ng-bootstrap/ng-bootstrap";
import { AppRoutingModule } from './app-routing.module';
import {MapComponent} from "./components/map/map.component";
import {ProjectsComponent} from "./pages/projects/projects.component";
import {TableComponent} from "./components/table/table.component";

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
        TableComponent],
  providers: [provideHttpClient(withInterceptorsFromDi())] })
export class AppModule { }
