import { Routes } from '@angular/router';

// Components
import { HomeComponent } from './pages/home.component';
import { ProfileComponent } from './pages/profile/profile.component';
import { ProjectsComponent } from './pages/projects/projects.component';
import { ProjectDetailComponent } from './pages/projects/project-detail/project-detail.component';

// Guards
import { AuthGuard } from './guards/auth.guard';
import { AdminGuard } from './guards/admin.guard';
import {RegistrationComponent} from './pages/registration/registration.component';
import {CallbackComponent} from './components/callback/callback.component';
import {AdminComponent} from './pages/admin/admin.component';
import {AdminProjectDetailComponent} from './pages/admin/admin-project-detail/admin-project-detail.component';

export const routes: Routes = [
  { path: '', component: HomeComponent },
  {
    path: 'profile',
    component: ProfileComponent,
    canActivate: [AuthGuard]
  },
  {
    path: 'projects',
    component: ProjectsComponent,
  },
  {
    path: 'projects/:id',
    component: ProjectDetailComponent,
  },
  {
    path: 'projects/:id/register',
    component: RegistrationComponent,
  },
  {
    path: 'admin',
    component: AdminComponent,
    canActivate: [AuthGuard, AdminGuard]
  },
  {
    path: 'admin/:id',
    component: AdminProjectDetailComponent,
    canActivate: [AuthGuard, AdminGuard]
  },
  {
    path: 'callback',
    component: CallbackComponent,
    canActivate: [AuthGuard]
  },
  { path: '**', redirectTo: '' }
];
