import {Component, OnInit} from '@angular/core';
import {APIService} from "@services";
import {Project} from "@models";

@Component({
    selector: 'app-projects',
    templateUrl: './projects.component.html',
    styleUrl: './projects.component.css',
    standalone: false,
})
export class ProjectsComponent implements OnInit {
  projects: Project[] = [];
  constructor(private APIService: APIService) {}

  ngOnInit() {
    this.APIService.getProjects().subscribe(data => {
      this.projects = data;
    })
  }

}
