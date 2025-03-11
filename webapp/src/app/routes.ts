import { Route} from '@angular/router';
import {ProjectsComponent} from "./pages/projects/projects.component";
import {AdminComponent} from "./pages/admin/admin.component";
import { authGuardFn } from "@auth0/auth0-angular";

export const ROUTES: Route[] = [
  // { path: '', component: MapComponent }, // Default route
  { path: 'projects', component: ProjectsComponent, canActivate: [authGuardFn], },
  {path: '', pathMatch: 'full', component: AdminComponent},
  // {path: 'profile', loadComponent: () => import('./features/profile/profile.component').then(mod => mod.ProfileComponent), canActivate: [AuthGuard]},
  // {path: 'public', loadComponent: () => import('./features/public/public.component').then(mod => mod.PublicComponent)},
  // {path: 'protected', loadComponent: () => import('./features/protected/protected.component').then(mod => mod.ProtectedComponent), canActivate: [AuthGuard]},
  { path: 'admin', component: AdminComponent, canActivate: [authGuardFn], },
  // {path: 'callback', loadComponent: () => import('./features/callback/callback.component').then(mod => mod.CallbackComponent)},
  // {path: '**', loadComponent: () => import('./features/not-found/not-found.component').then(mod => mod.NotFoundComponent)},
];
