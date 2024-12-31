import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import {RouterModule, Routes} from '@angular/router';
import {ProjectsComponent} from "./pages/projects/projects.component";
import {MapComponent} from "./components/map/map.component";

const routes: Routes = [
  { path: '', component: MapComponent }, // Default route
  { path: 'projects', component: ProjectsComponent }
];

@NgModule({
  declarations: [],
  imports: [
    CommonModule, RouterModule.forRoot(routes)
  ],
  exports: [RouterModule]
})
export class AppRoutingModule { }
