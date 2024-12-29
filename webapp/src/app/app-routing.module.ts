import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import {RouterModule, Routes} from '@angular/router';
import {TechComponent} from "./tech/tech.component";
import {ProjectsComponent} from "./pages/projects/projects.component";

const routes: Routes = [
  { path: '', component: ProjectsComponent }, // Default route
  { path: 'contact', component: TechComponent }
];

@NgModule({
  declarations: [],
  imports: [
    CommonModule, RouterModule.forRoot(routes)
  ],
  exports: [RouterModule]
})
export class AppRoutingModule { }
