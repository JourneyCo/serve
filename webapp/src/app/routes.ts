import { Route } from '@angular/router';
import { AuthGuard } from '@auth0/auth0-angular';
import {MapComponent} from "@components";
import {ProjectsComponent} from "./pages/projects/projects.component";

export const ROUTES: Route[] = [
  // { path: '', component: MapComponent }, // Default route
  // { path: 'projects', component: ProjectsComponent }
  {path: '', pathMatch: 'full', loadComponent: () => import('./components/map/map.component').then(mod => mod.MapComponent)},
  // {path: 'profile', loadComponent: () => import('./features/profile/profile.component').then(mod => mod.ProfileComponent), canActivate: [AuthGuard]},
  // {path: 'public', loadComponent: () => import('./features/public/public.component').then(mod => mod.PublicComponent)},
  // {path: 'protected', loadComponent: () => import('./features/protected/protected.component').then(mod => mod.ProtectedComponent), canActivate: [AuthGuard]},
  // {path: 'admin', loadComponent: () => import('./features/admin/admin.component').then(mod => mod.AdminComponent), canActivate: [AuthGuard]},
  // {path: 'callback', loadComponent: () => import('./features/callback/callback.component').then(mod => mod.CallbackComponent)},
  // {path: '**', loadComponent: () => import('./features/not-found/not-found.component').then(mod => mod.NotFoundComponent)},
];
