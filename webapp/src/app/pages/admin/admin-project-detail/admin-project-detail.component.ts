import { Component } from '@angular/core';
import {ProjectDetailComponent} from '../../projects/project-detail/project-detail.component';
import {AdminProjectPanelComponent} from '@components';
import {ProjectService} from '@services';
import {Observable} from 'rxjs';
import {ActivatedRoute} from '@angular/router';
import {Project} from '@models';

@Component({
  selector: 'app-admin-project-detail',
  imports: [
    ProjectDetailComponent,
    AdminProjectPanelComponent
  ],
  templateUrl: './admin-project-detail.component.html',
  styleUrl: './admin-project-detail.component.scss'
})
export class AdminProjectDetailComponent {
  isAdmin: Observable<boolean>;
  proj_id: string;
  project: Project

  constructor(
      private route: ActivatedRoute,
      private projectService: ProjectService
  ) {
    this.getProject();
  }

  getProject() {
    this.proj_id = this.route.snapshot.paramMap.get("id") || "";
    const id = Number(this.proj_id);
    this.projectService.getProject(id).subscribe(data => {
      this.project = data
    })
  }


}
