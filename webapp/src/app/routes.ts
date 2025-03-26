import { Route} from '@angular/router';
import {ProjectsComponent} from "./pages/projects/projects.component";
import { authGuardFn } from "@auth0/auth0-angular";
import {PageNotFoundComponent} from "./pages/pagenotfound/pagenotfound.component";
import {ProjectComponent} from "./pages/projects/project/project.component";

export const ROUTES: Route[] = [
  { path: '', pathMatch: 'full', component: ProjectsComponent, canActivate: [authGuardFn], },
  { path: 'projects/:id', component: ProjectComponent, canActivate: [authGuardFn], },
  { path: '**', component: PageNotFoundComponent }, // must be last in order
];
